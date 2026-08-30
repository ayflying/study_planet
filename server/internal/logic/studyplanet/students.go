// Package studyplanet 学生账号管理（家长鉴权）：增删改查 + 数据清理。
package studyplanet

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "studyplanet/api/studyplanet/v1"
)

// ListStudents 学生列表（挂可选鉴权）：
//   - 匿名（学生端未登录家长）：返回空列表——匿名用户不应看到任何家庭的孩子；
//     前端家长登录（Casdoor）后本接口自动返回本家孩子。
//   - 带家长 token：只返回归属自己的孩子（数据隔离）。
func (s *sStudyPlanet) ListStudents(ctx context.Context, req *v1.ListStudentsReq) (res *v1.ListStudentsRes, err error) {
	parentID := ctxParentID(ctx)
	out := make(v1.ListStudentsRes, 0, 4)
	if parentID <= 0 {
		return &out, nil
	}
	all, err := daoChildren.Ctx(ctx).Where("parent_id", parentID).Order("id").All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询学生失败")
	}
	for _, r := range all {
		out = append(out, studentOf(r))
	}
	return &out, nil
}

// CreateStudent 新建学生账号：自动归属当前登录家长。
func (s *sStudyPlanet) CreateStudent(ctx context.Context, req *v1.CreateStudentReq) (res *v1.CreateStudentRes, err error) {
	parentID := ctxParentID(ctx)
	if parentID <= 0 {
		return nil, errAuth("登录状态缺少家长身份，请重新登录")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errParam("请填写学生姓名")
	}
	avatar, grade := req.Avatar, req.Grade
	if avatar == "" {
		avatar = "🚀"
	}
	if grade <= 0 {
		grade = 5
	}
	if req.Username != "" {
		cnt, err := daoChildren.Ctx(ctx).Where("username", req.Username).Count()
		if err != nil {
			return nil, gerror.Wrap(err, "校验用户名失败")
		}
		if cnt > 0 {
			return nil, errParam("用户名已被使用")
		}
	}
	id, err := daoChildren.Ctx(ctx).Data(doChildren{
		Name:     req.Name,
		Username: req.Username,
		Avatar:   avatar,
		Grade:    grade,
		ParentId: parentID,
	}).InsertAndGetId()
	if err != nil {
		return nil, gerror.Wrap(err, "创建学生失败")
	}
	rec, err := daoChildren.Ctx(ctx).Where("id", id).One()
	if err != nil || rec.IsEmpty() {
		return nil, errNotFound("查询新学生失败")
	}
	st := v1.CreateStudentRes(studentOf(rec))
	return &st, nil
}

// UpdateStudent 修改学生信息（姓名/头像/年级/用户名）：仅限自己的孩子。
func (s *sStudyPlanet) UpdateStudent(ctx context.Context, req *v1.UpdateStudentReq) (res *v1.UpdateStudentRes, err error) {
	if err := s.ensureChildOwned(ctx, req.ID); err != nil {
		return nil, err
	}
	if req.Name != nil && *req.Name == "" {
		return nil, errParam("姓名不能为空")
	}
	if req.Username != nil && *req.Username != "" {
		cnt, err := daoChildren.Ctx(ctx).Where("username", *req.Username).Where("id<>?", req.ID).Count()
		if err != nil {
			return nil, gerror.Wrap(err, "校验用户名失败")
		}
		if cnt > 0 {
			return nil, errParam("用户名已被使用")
		}
	}
	data := g.Map{}
	if req.Name != nil {
		data["name"] = *req.Name
	}
	if req.Username != nil {
		data["username"] = *req.Username
	}
	if req.Avatar != nil {
		data["avatar"] = *req.Avatar
	}
	if req.Grade != nil {
		data["grade"] = *req.Grade
	}
	if len(data) > 0 {
		if _, err := daoChildren.Ctx(ctx).Where("id", req.ID).Data(data).Update(); err != nil {
			return nil, gerror.Wrap(err, "更新学生失败")
		}
	}
	rec, err := daoChildren.Ctx(ctx).Where("id", req.ID).One()
	if err != nil || rec.IsEmpty() {
		return nil, errNotFound("学生不存在")
	}
	st := v1.UpdateStudentRes(studentOf(rec))
	return &st, nil
}

// DeleteStudent 删除学生；至少保留一个，且清空其学习数据。仅限自己的孩子。
func (s *sStudyPlanet) DeleteStudent(ctx context.Context, req *v1.DeleteStudentReq) (res *v1.DeleteStudentRes, err error) {
	if err := s.ensureChildOwned(ctx, req.ID); err != nil {
		return nil, err
	}
	cnt, err := daoChildren.Ctx(ctx).Where("parent_id", ctxParentID(ctx)).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计学生失败")
	}
	if cnt <= 1 {
		return nil, errParam("至少保留一个学生账号")
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先清场次答题（依赖 session_id），再清其余学习数据，最后删学生
		if _, err := tx.Exec("DELETE FROM session_answers WHERE session_id IN (SELECT id FROM practice_sessions WHERE child_id=?)", req.ID); err != nil {
			return err
		}
		for _, d := range []*gdb.Model{
			daoWrongQ.Ctx(ctx).Where("child_id", req.ID),
			daoSessions.Ctx(ctx).Where("child_id", req.ID),
			daoWordProg.Ctx(ctx).Where("child_id", req.ID),
			daoPointsLog.Ctx(ctx).Where("child_id", req.ID),
			daoTasks.Ctx(ctx).Where("child_id", req.ID),
			daoRedempt.Ctx(ctx).Where("child_id", req.ID),
			daoChildren.Ctx(ctx).Where("id", req.ID),
		} {
			if _, err := d.Delete(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, gerror.Wrap(err, "删除学生失败")
	}
	return &v1.DeleteStudentRes{OK: true}, nil
}
