// Package studyplanet 每日任务模块：列表 / 完成 / 家长发布与删除。
package studyplanet

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
)

// ListTasks 学生任务列表（按 student_id，公开接口）。
func (s *sStudyPlanet) ListTasks(ctx context.Context, req *v1.ListTasksReq) (res *v1.ListTasksRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	m := daoTasks.Ctx(ctx).Where("child_id", cid)
	if req.Status != "" {
		m = m.Where("status", req.Status)
	}
	all, err := m.Order("due_date").All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询任务失败")
	}
	today := todayStr()
	out := make(v1.ListTasksRes, 0, len(all))
	for _, r := range all {
		t := taskOf(r)
		if t.Status != "done" && t.DueDate != "" && t.DueDate < today {
			t.Status = "overdue"
		}
		out = append(out, t)
	}
	return &out, nil
}

// CompleteTask 完成任务（学生操作，公开接口）。
func (s *sStudyPlanet) CompleteTask(ctx context.Context, req *v1.CompleteTaskReq) (res *v1.CompleteTaskRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	t, err := daoTasks.Ctx(ctx).Where("id", req.ID).Where("child_id", cid).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询任务失败")
	}
	if t.IsEmpty() {
		return nil, errNotFound("未找到该任务")
	}
	if t["status"].String() == "done" {
		return &v1.CompleteTaskRes{OK: true, Already: true}, nil
	}
	if _, err := daoTasks.Ctx(ctx).Where("id", req.ID).Data(doTasks{
		Status:      "done",
		CompletedAt: gtime.Now(),
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "更新任务失败")
	}
	points := t["points"].Int()
	s.award(cid, points, "完成任务:+"+g.NewVar(points).String())
	// 完成任务可能掉落零食
	s.addSnackDrop(ctx, cid)
	return &v1.CompleteTaskRes{OK: true}, nil
}

// AddTask 发布任务（家长鉴权）：学生必须归属当前家长。
func (s *sStudyPlanet) AddTask(ctx context.Context, req *v1.AddTaskReq) (res *v1.AddTaskRes, err error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errParam("请填写任务名称")
	}
	if err := s.ensureChildOwned(ctx, req.StudentID); err != nil {
		return nil, err
	}
	data := doTasks{
		Title:   req.Title,
		Type:    req.Type,
		Points:  req.Points,
		Status:  "pending",
		ChildId: req.StudentID,
	}
	// due_date 为空时留空（MySQL DATE 列不接受空字符串）
	if strings.TrimSpace(req.DueDate) != "" {
		if t, err := gtime.StrToTime(req.DueDate, "Y-m-d"); err == nil {
			data.DueDate = t
		} else {
			data.DueDate = gtime.NewFromStr(req.DueDate)
		}
	}
	if _, err := daoTasks.Ctx(ctx).Data(data).Insert(); err != nil {
		return nil, gerror.Wrap(err, "发布任务失败")
	}
	return &v1.AddTaskRes{OK: true}, nil
}

// StudentAddTask 学生自建任务（公开接口，无需家长鉴权）。
func (s *sStudyPlanet) StudentAddTask(ctx context.Context, req *v1.StudentAddTaskReq) (res *v1.StudentAddTaskRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	if strings.TrimSpace(req.Title) == "" {
		return nil, errParam("请填写任务名称")
	}
	data := doTasks{
		Title:   req.Title,
		Type:    req.Type,
		Points:  req.Points,
		Status:  "pending",
		ChildId: cid,
	}
	if req.Points <= 0 || req.Points > 50 {
		data.Points = 5 // 学生自建任务默认5积分，上限50
	}
	if strings.TrimSpace(req.DueDate) != "" {
		if t, err := gtime.StrToTime(req.DueDate, "Y-m-d"); err == nil {
			data.DueDate = t
		} else {
			data.DueDate = gtime.NewFromStr(req.DueDate)
		}
	}
	if _, err := daoTasks.Ctx(ctx).Data(data).Insert(); err != nil {
		return nil, gerror.Wrap(err, "创建任务失败")
	}
	return &v1.StudentAddTaskRes{OK: true}, nil
}

// StudentDeleteTask 学生删除自建任务（公开接口，仅删除自己名下任务）。
func (s *sStudyPlanet) StudentDeleteTask(ctx context.Context, req *v1.StudentDeleteTaskReq) (res *v1.StudentDeleteTaskRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	t, err := daoTasks.Ctx(ctx).Where("id", req.ID).Where("child_id", cid).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询任务失败")
	}
	if t.IsEmpty() {
		return nil, errNotFound("未找到该任务")
	}
	if _, err := daoTasks.Ctx(ctx).Where("id", req.ID).Delete(); err != nil {
		return nil, gerror.Wrap(err, "删除任务失败")
	}
	return &v1.StudentDeleteTaskRes{OK: true}, nil
}

// DeleteTask 删除任务（家长鉴权）：任务须属于当前家长的学生。
func (s *sStudyPlanet) DeleteTask(ctx context.Context, req *v1.DeleteTaskReq) (res *v1.DeleteTaskRes, err error) {
	t, err := daoTasks.Ctx(ctx).Fields("id,child_id").Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询任务失败")
	}
	if t.IsEmpty() {
		return nil, errNotFound("未找到该任务")
	}
	if err := s.ensureChildOwned(ctx, t["child_id"].Int()); err != nil {
		return nil, err
	}
	if _, err := daoTasks.Ctx(ctx).Where("id", req.ID).Delete(); err != nil {
		return nil, gerror.Wrap(err, "删除任务失败")
	}
	return &v1.DeleteTaskRes{OK: true}, nil
}
