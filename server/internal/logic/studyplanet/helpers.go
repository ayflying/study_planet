// Package studyplanet 通用读取助手：记录 → api 结构转换、学生解析与归属校验。
package studyplanet

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "studyplanet/api/studyplanet/v1"
)

// studentOf children 记录 → 学生档案。
func studentOf(r gdb.Record) v1.Student {
	return v1.Student{
		ID:        r["id"].Int(),
		Name:      r["name"].String(),
		Username:  r["username"].String(),
		Avatar:    r["avatar"].String(),
		Grade:     r["grade"].Int(),
		CreatedAt: r["created_at"].String(),
	}
}

// sessionOf practice_sessions 记录 → 练习场次（时间列可空统一转字符串）。
func sessionOf(r gdb.Record) v1.PracticeSession {
	return v1.PracticeSession{
		ID:         r["id"].Int(),
		ChildID:    r["child_id"].Int(),
		Subject:    r["subject"].String(),
		Level:      r["level"].Int(),
		Total:      r["total"].Int(),
		Correct:    r["correct"].Int(),
		MaxCombo:   r["max_combo"].Int(),
		Bonus:      r["bonus"].Int(),
		Stars:      r["stars"].Int(),
		Finished:   r["finished"].Int(),
		CreatedAt:  r["created_at"].String(),
		FinishedAt: r["finished_at"].String(),
	}
}

// taskOf tasks 记录 → 任务（due_date/completed_at 可空，空值输出 ""）。
func taskOf(r gdb.Record) v1.Task {
	due, completed := "", ""
	if !r["due_date"].IsNil() {
		due = r["due_date"].GTime().Format("Y-m-d")
	}
	if !r["completed_at"].IsNil() {
		completed = r["completed_at"].GTime().Format("Y-m-d H:i:s")
	}
	return v1.Task{
		ID:          r["id"].Int(),
		Title:       r["title"].String(),
		Type:        r["type"].String(),
		DueDate:     due,
		Points:      r["points"].Int(),
		Status:      r["status"].String(),
		CreatedAt:   r["created_at"].String(),
		CompletedAt: completed,
	}
}

// resolveChild 解析本次请求对应的学生 id：来自 req 的 student_id，缺省为 1（兼容旧客户端）。
// 学生不存在时返回 -1。
func (s *sStudyPlanet) resolveChild(ctx context.Context, studentID int) (int, error) {
	id := studentID
	if id <= 0 {
		id = 1
	}
	rec, err := daoChildren.Ctx(ctx).Fields("id,parent_id").Where("id", id).One()
	if err != nil {
		return -1, gerror.Wrap(err, "查询学生失败")
	}
	if rec.IsEmpty() {
		return -1, nil
	}
	return id, nil
}

// childParentID 查询孩子的归属家长 id（孩子不存在返回 0）。
func (s *sStudyPlanet) childParentID(ctx context.Context, childID int) int {
	v, err := daoChildren.Ctx(ctx).Fields("parent_id").Where("id", childID).Value()
	if err != nil || v == nil || v.IsNil() {
		return 0
	}
	return v.Int()
}

// ensureChildOwned 家长端操作前校验孩子归属：不属于自己的孩子一律拒绝。
// requireOwner=true（家长端写操作）时校验学生归属：parentID 与孩子 parent_id 不一致即拒绝，
// NULL 归属（历史数据尚未被接管）同样拒绝，防止越权操作他人孩子。
func (s *sStudyPlanet) ensureChildOwned(ctx context.Context, childID int) error {
	parentID := ctxParentID(ctx)
	if parentID <= 0 {
		return errAuth("登录状态缺少家长身份，请重新登录")
	}
	owner := s.childParentID(ctx, childID)
	if owner == 0 {
		return errNotFound("学生不存在")
	}
	if owner != parentID {
		return errForbidden("无权操作该学生")
	}
	return nil
}

// ensureParentAuth 确认当前请求带家长身份。
func (s *sStudyPlanet) ensureParentAuth(ctx context.Context) error {
	if ctxParentID(ctx) <= 0 {
		return errAuth("登录状态缺少家长身份，请重新登录")
	}
	return nil
}

// ensureName 通用非空校验。
func ensureName(name, msg string) error {
	if strings.TrimSpace(name) == "" {
		return errParam(msg)
	}
	return nil
}
