package handler

import (
	"bytes"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"studyplanet/internal/config"
	"studyplanet/internal/model"
)

// AppVersion 由构建时注入（Dockerfile -ldflags -X 写入），未注入时读 VERSION 文件。
var AppVersion string

// Store 持有数据库连接与运行配置，作为所有 handler 的接收者。
type Store struct {
	DB  *sqlx.DB
	Cfg *config.Config
}

func NewStore(db *sqlx.DB, cfg *config.Config) *Store {
	return &Store{DB: db, Cfg: cfg}
}

func (s *Store) ok(r *ghttp.Request, data interface{}) {
	r.Response.WriteJson(data)
}

func (s *Store) fail(r *ghttp.Request, code int, msg string) {
	r.Response.WriteStatusExit(code, map[string]interface{}{"error": msg})
}

// resolveChild 解析本次请求对应的学生 id：查询参数 student_id，缺省为 1（兼容旧客户端）。
// 学生不存在时返回 -1。
func (s *Store) resolveChild(r *ghttp.Request) int {
	id := r.GetQuery("student_id").Int()
	if id <= 0 {
		id = 1
	}
	var cnt int
	if err := s.DB.Get(&cnt, "SELECT COUNT(*) FROM children WHERE id=?", id); err != nil || cnt == 0 {
		return -1
	}
	return id
}

// award 记录指定学生的积分变动（不阻塞主流程，失败仅记日志）。
func (s *Store) award(childID int, delta int, reason string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.DB.Exec(
		"INSERT INTO points_log(child_id,delta,reason,created_at) VALUES(?,?,?,?)",
		childID, delta, reason, now,
	)
	if err != nil {
		g.Log().Errorf(gctx.New(), "award 记录积分失败 child=%d delta=%d: %v", childID, delta, err)
	}
}

func (s *Store) idParam(r *ghttp.Request) int {
	id, _ := strconv.Atoi(r.Get("id").String())
	return id
}

// ---------- 健康检查 ----------
func (s *Store) Health(r *ghttp.Request) {
	s.ok(r, map[string]interface{}{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"app":     "studyplanet",
		"version": CurrentVersion(),
	})
}

// CurrentVersion 返回构建注入的版本号（未注入时读 VERSION 文件，兜底 dev）。
func CurrentVersion() string {
	if AppVersion != "" {
		return AppVersion
	}
	if b, err := os.ReadFile("VERSION"); err == nil {
		return string(bytes.TrimSpace(b))
	}
	return "dev"
}

// ---------- 单词卡片 ----------
func (s *Store) ListWords(r *ghttp.Request) {
	level := r.GetQuery("level").String()
	q := "SELECT id,level,word,meaning,phonetic,example,created_at FROM words"
	args := []interface{}{}
	if level != "" {
		q += " WHERE level=?"
		args = append(args, level)
	}
	q += " ORDER BY level, id"
	var words []model.Word
	if err := s.DB.Select(&words, q, args...); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, words)
}

func (s *Store) WordDetail(r *ghttp.Request) {
	id := s.idParam(r)
	var w model.Word
	if err := s.DB.Get(&w, "SELECT id,level,word,meaning,phonetic,example,created_at FROM words WHERE id=?", id); err != nil {
		s.fail(r, 404, "未找到该单词")
		return
	}
	known := 0
	var p model.WordProgress
	if err := s.DB.Get(&p, "SELECT word_id,child_id,known FROM word_progress WHERE word_id=? AND child_id=?", id, s.resolveChild(r)); err == nil {
		known = p.Known
	}
	s.ok(r, map[string]interface{}{"word": w, "known": known})
}

func (s *Store) WordProgress(r *ghttp.Request) {
	id := s.idParam(r)
	var body struct {
		Known     bool `json:"known"`
		SessionID int  `json:"session_id"` // 可选：传入则走连击+场次计分
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	known := 0
	if body.Known {
		known = 1
	}
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.DB.Exec(
		`INSERT INTO word_progress(word_id,child_id,known,last_reviewed) VALUES(?,?,?,?)
		 ON CONFLICT(word_id,child_id) DO UPDATE SET known=excluded.known, last_reviewed=excluded.last_reviewed`,
		id, cid, known, now,
	)
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	if body.SessionID > 0 && body.Known {
		s.recordAnswer(r, body.SessionID, id, true, 5, "")
		return
	}
	if body.Known {
		s.award(cid, 5, "单词认读:+5")
	}
	s.ok(r, map[string]interface{}{"ok": true, "known": known})
}

// ---------- 语文阅读 ----------
func (s *Store) ReadingDetail(r *ghttp.Request) {
	id := s.idParam(r)
	var rd model.Reading
	if err := s.DB.Get(&rd, "SELECT id,title,content,level FROM readings WHERE id=?", id); err != nil {
		s.fail(r, 404, "未找到该阅读")
		return
	}
	var qs []model.ReadingQuestion
	if err := s.DB.Select(&qs, "SELECT id,reading_id,question,option_a,option_b,option_c,option_d,answer FROM reading_questions WHERE reading_id=? ORDER BY id", id); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"reading": rd, "questions": qs})
}

func (s *Store) ReadingAnswer(r *ghttp.Request) {
	var body struct {
		QuestionID int    `json:"question_id"`
		Answer     string `json:"answer"`
		SessionID  int    `json:"session_id"` // 可选：传入则走连击+场次计分
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	var q model.ReadingQuestion
	if err := s.DB.Get(&q, "SELECT id,answer FROM reading_questions WHERE id=?", body.QuestionID); err != nil {
		s.fail(r, 404, "未找到该题目")
		return
	}
	correct := strings.EqualFold(strings.TrimSpace(q.Answer), strings.TrimSpace(body.Answer))
	if body.SessionID > 0 {
		s.recordAnswer(r, body.SessionID, body.QuestionID, correct, 2, "")
		return
	}
	if correct {
		if cid := s.resolveChild(r); cid > 0 {
			s.award(cid, 2, "阅读答题:+2")
		}
	}
	s.ok(r, map[string]interface{}{"correct": correct, "correct_answer": q.Answer})
}

// ---------- 数学题目 ----------
func (s *Store) ListMath(r *ghttp.Request) {
	level := r.GetQuery("level").String()
	q := "SELECT id,level,type,question,options,answer,explanation FROM math_problems"
	args := []interface{}{}
	if level != "" {
		q += " WHERE level=?"
		args = append(args, level)
	}
	q += " ORDER BY level, id"
	var ps []model.MathProblem
	if err := s.DB.Select(&ps, q, args...); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, ps)
}

func (s *Store) MathAnswer(r *ghttp.Request) {
	id := s.idParam(r)
	var body struct {
		Answer    string `json:"answer"`
		SessionID int    `json:"session_id"` // 可选：传入则走连击+场次计分
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	var p model.MathProblem
	if err := s.DB.Get(&p, "SELECT id,answer,explanation FROM math_problems WHERE id=?", id); err != nil {
		s.fail(r, 404, "未找到该题目")
		return
	}
	correct := strings.EqualFold(strings.TrimSpace(p.Answer), strings.TrimSpace(body.Answer))
	if body.SessionID > 0 {
		s.recordAnswer(r, body.SessionID, id, correct, 3, "")
		return
	}
	if correct {
		if cid := s.resolveChild(r); cid > 0 {
			s.award(cid, 3, "数学答题:+3")
		}
	}
	s.ok(r, map[string]interface{}{"correct": correct, "explanation": p.Explanation, "answer": p.Answer})
}

// ---------- 每日任务 ----------
func (s *Store) ListTasks(r *ghttp.Request) {
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	status := r.GetQuery("status").String()
	q := "SELECT id,title,type,COALESCE(due_date,'') AS due_date,points,status,created_at,COALESCE(completed_at,'') AS completed_at FROM tasks WHERE child_id=?"
	args := []interface{}{cid}
	if status != "" {
		q += " AND status=?"
		args = append(args, status)
	}
	q += " ORDER BY due_date"
	var tasks []model.Task
	if err := s.DB.Select(&tasks, q, args...); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	today := time.Now().Format("2006-01-02")
	for i := range tasks {
		if tasks[i].Status != "done" && tasks[i].DueDate != "" && tasks[i].DueDate < today {
			tasks[i].Status = "overdue"
		}
	}
	s.ok(r, tasks)
}

func (s *Store) CompleteTask(r *ghttp.Request) {
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	id := s.idParam(r)
	var t model.Task
	if err := s.DB.Get(&t, "SELECT id,points,status FROM tasks WHERE id=? AND child_id=?", id, cid); err != nil {
		s.fail(r, 404, "未找到该任务")
		return
	}
	if t.Status == "done" {
		s.ok(r, map[string]interface{}{"ok": true, "already": true})
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	if _, err := s.DB.Exec("UPDATE tasks SET status='done', completed_at=? WHERE id=?", now, id); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.award(cid, t.Points, "完成任务:"+strconv.Itoa(t.Points))
	s.ok(r, map[string]interface{}{"ok": true})
}

// ---------- 积分 ----------
func (s *Store) PointsSummary(r *ghttp.Request) {
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	var total int
	if err := s.DB.Get(&total, "SELECT COALESCE(SUM(delta),0) FROM points_log WHERE child_id=?", cid); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	today := time.Now().Format("2006-01-02")
	var todayEarned int
	if err := s.DB.Get(&todayEarned, "SELECT COALESCE(SUM(delta),0) FROM points_log WHERE child_id=? AND date(created_at)=?", cid, today); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"total": total, "today_earned": todayEarned, "student_id": cid})
}

func (s *Store) PointsLog(r *ghttp.Request) {
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	var logs []model.PointsLog
	if err := s.DB.Select(&logs, "SELECT id,child_id,delta,reason,created_at FROM points_log WHERE child_id=? ORDER BY id DESC LIMIT 100", cid); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, logs)
}

// ---------- 奖励 / 兑换 ----------
func (s *Store) ListRewards(r *ghttp.Request) {
	var rs []model.Reward
	if err := s.DB.Select(&rs, "SELECT id,name,cost_points,status FROM rewards ORDER BY cost_points"); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, rs)
}

func (s *Store) Redeem(r *ghttp.Request) {
	id := s.idParam(r)
	var rw model.Reward
	if err := s.DB.Get(&rw, "SELECT id,name,cost_points,status FROM rewards WHERE id=?", id); err != nil {
		s.fail(r, 404, "未找到该奖励")
		return
	}
	if rw.Status != "active" {
		s.fail(r, 400, "该奖励暂不可用")
		return
	}
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	var total int
	if err := s.DB.Get(&total, "SELECT COALESCE(SUM(delta),0) FROM points_log WHERE child_id=?", cid); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	if total < rw.CostPoints {
		s.fail(r, 400, "积分不足")
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	if _, err := s.DB.Exec("INSERT INTO redemptions(reward_id,status,requested_at,child_id) VALUES(?, 'pending', ?, ?)", id, now, cid); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"ok": true, "pending": true, "message": "已提交兑换，等待家长确认"})
}

// ---------- 家长端 ----------
func (s *Store) ParentLogin(r *ghttp.Request) {
	// Casdoor 已配置时禁用 PIN 登录，强制走 SSO
	if s.Cfg.Casdoor.Enabled() {
		s.fail(r, http.StatusBadRequest, "已启用 Casdoor 登录，请使用 SSO")
		return
	}
	var body struct {
		Pin string `json:"pin"`
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	var hash string
	if err := s.DB.Get(&hash, "SELECT value FROM settings WHERE key='parent_pin'"); err != nil {
		s.fail(r, 500, "PIN 未设置")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Pin)) != nil {
		s.fail(r, http.StatusUnauthorized, "密码错误")
		return
	}
	tok, err := issueToken(s.Cfg.Parent.JWTSecret, "家长")
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"token": tok})
}

// ---------- 以下为家长鉴权后接口 ----------
func (s *Store) AddTask(r *ghttp.Request) {
	var body struct {
		Title     string `json:"title"`
		Type      string `json:"type"`
		DueDate   string `json:"due_date"`
		Points    int    `json:"points"`
		StudentID int    `json:"student_id"`
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	if body.Title == "" {
		s.fail(r, 400, "请填写任务名称")
		return
	}
	if _, err := s.DB.Exec("INSERT INTO tasks(title,type,due_date,points,status,child_id) VALUES(?,?,?,?,'pending',?)", body.Title, body.Type, body.DueDate, body.Points, body.StudentID); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"ok": true})
}

func (s *Store) DeleteTask(r *ghttp.Request) {
	id := s.idParam(r)
	if _, err := s.DB.Exec("DELETE FROM tasks WHERE id=?", id); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"ok": true})
}

func (s *Store) AddReward(r *ghttp.Request) {
	var body struct {
		Name       string `json:"name"`
		CostPoints int    `json:"cost_points"`
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	if body.Name == "" {
		s.fail(r, 400, "请填写奖励名称")
		return
	}
	if _, err := s.DB.Exec("INSERT INTO rewards(name,cost_points,status) VALUES(?,?,'active')", body.Name, body.CostPoints); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"ok": true})
}

func (s *Store) ConfirmRedemption(r *ghttp.Request) {
	id := s.idParam(r)
	var rd model.Redemption
	if err := s.DB.Get(&rd, "SELECT id,reward_id,status FROM redemptions WHERE id=?", id); err != nil {
		s.fail(r, 404, "未找到该兑换")
		return
	}
	if rd.Status != "pending" {
		s.fail(r, 400, "该兑换不在待确认状态")
		return
	}
	var rw model.Reward
	if err := s.DB.Get(&rw, "SELECT id,name,cost_points FROM rewards WHERE id=?", rd.RewardID); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	if _, err := s.DB.Exec("UPDATE redemptions SET status='confirmed', confirmed_at=? WHERE id=?", now, id); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	if _, err := s.DB.Exec("UPDATE rewards SET status='redeemed' WHERE id=?", rd.RewardID); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.award(rd.ChildID, -rw.CostPoints, "兑换:"+rw.Name)
	s.ok(r, map[string]interface{}{"ok": true})
}

func (s *Store) SetPin(r *ghttp.Request) {
	var body struct {
		Pin string `json:"pin"`
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	if len(body.Pin) < 4 {
		s.fail(r, 400, "PIN 至少 4 位")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Pin), 10)
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	if _, err := s.DB.Exec("INSERT INTO settings(key,value) VALUES('parent_pin',?) ON CONFLICT(key) DO UPDATE SET value=?", string(hash)); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"ok": true})
}

// ---------- 学生账号管理（家长鉴权后） ----------
// ListStudents 全部学生（家长切换用）。
func (s *Store) ListStudents(r *ghttp.Request) {
	var st []model.Student
	if err := s.DB.Select(&st, "SELECT id,name,username,avatar,grade,created_at FROM children ORDER BY id"); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, st)
}

// CreateStudent 新建学生账号。
func (s *Store) CreateStudent(r *ghttp.Request) {
	var body struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Grade    int    `json:"grade"`
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	if body.Name == "" {
		s.fail(r, 400, "请填写学生姓名")
		return
	}
	if body.Avatar == "" {
		body.Avatar = "🚀"
	}
	if body.Grade <= 0 {
		body.Grade = 5
	}
	if body.Username != "" {
		var cnt int
		if err := s.DB.Get(&cnt, "SELECT COUNT(*) FROM children WHERE username=?", body.Username); err != nil {
			s.fail(r, 500, err.Error())
			return
		}
		if cnt > 0 {
			s.fail(r, 400, "用户名已被使用")
			return
		}
	}
	if _, err := s.DB.Exec("INSERT INTO children(name,username,avatar,grade) VALUES(?,?,?,?)",
		body.Name, body.Username, body.Avatar, body.Grade); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	var st model.Student
	if err := s.DB.Get(&st, "SELECT id,name,username,avatar,grade,created_at FROM children WHERE id=last_insert_rowid()"); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, st)
}

// UpdateStudent 修改学生信息（姓名/头像/年级/用户名）。
func (s *Store) UpdateStudent(r *ghttp.Request) {
	id := s.idParam(r)
	var body struct {
		Name     *string `json:"name"`
		Username *string `json:"username"`
		Avatar   *string `json:"avatar"`
		Grade    *int    `json:"grade"`
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	if body.Name != nil && *body.Name == "" {
		s.fail(r, 400, "姓名不能为空")
		return
	}
	if body.Username != nil && *body.Username != "" {
		var cnt int
		if err := s.DB.Get(&cnt, "SELECT COUNT(*) FROM children WHERE username=? AND id<>?", *body.Username, id); err != nil {
			s.fail(r, 500, err.Error())
			return
		}
		if cnt > 0 {
			s.fail(r, 400, "用户名已被使用")
			return
		}
	}
	if body.Name != nil {
		if _, err := s.DB.Exec("UPDATE children SET name=? WHERE id=?", *body.Name, id); err != nil {
			s.fail(r, 500, err.Error())
			return
		}
	}
	if body.Username != nil {
		if _, err := s.DB.Exec("UPDATE children SET username=? WHERE id=?", *body.Username, id); err != nil {
			s.fail(r, 500, err.Error())
			return
		}
	}
	if body.Avatar != nil {
		if _, err := s.DB.Exec("UPDATE children SET avatar=? WHERE id=?", *body.Avatar, id); err != nil {
			s.fail(r, 500, err.Error())
			return
		}
	}
	if body.Grade != nil {
		if _, err := s.DB.Exec("UPDATE children SET grade=? WHERE id=?", *body.Grade, id); err != nil {
			s.fail(r, 500, err.Error())
			return
		}
	}
	var st model.Student
	if err := s.DB.Get(&st, "SELECT id,name,username,avatar,grade,created_at FROM children WHERE id=?", id); err != nil {
		s.fail(r, 404, "学生不存在")
		return
	}
	s.ok(r, st)
}

// DeleteStudent 删除学生；至少保留一个，且清空其学习数据。
func (s *Store) DeleteStudent(r *ghttp.Request) {
	id := s.idParam(r)
	var cnt int
	if err := s.DB.Get(&cnt, "SELECT COUNT(*) FROM children"); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	if cnt <= 1 {
		s.fail(r, 400, "至少保留一个学生账号")
		return
	}
	for _, q := range []string{
		"DELETE FROM word_progress WHERE child_id=?",
		"DELETE FROM points_log WHERE child_id=?",
		"DELETE FROM tasks WHERE child_id=?",
		"DELETE FROM redemptions WHERE child_id=?",
	} {
		if _, err := s.DB.Exec(q, id); err != nil {
			s.fail(r, 500, err.Error())
			return
		}
	}
	if _, err := s.DB.Exec("DELETE FROM children WHERE id=?", id); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"ok": true})
}

// ---------- JWT ----------
func issueToken(secret, name string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  "parent",
		"name": name,
		"exp":  time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
