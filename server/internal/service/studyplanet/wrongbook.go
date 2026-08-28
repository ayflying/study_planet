package studyplanet

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"studyplanet/internal/leaderboard"
)

// gLog 非关键路径日志（失败不影响主流程）。
func gLog(format string, args ...interface{}) {
	g.Log().Errorf(gctx.New(), format, args...)
}

// ---------- 错题本：答错登记 / 巩固练习出题 ----------

// recordWrong 答错时登记错题（已 resolved 的重新激活，wrong_count 累加）。
// 只做附加记录，失败不影响作答主流程。
func (s *Store) recordWrong(childID int, subject string, refID int) {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.DB.Exec(
		`INSERT INTO wrong_questions(child_id,subject,ref_id,wrong_count,resolved,last_wrong_at) VALUES(?,?,?,1,0,?)
		 ON DUPLICATE KEY UPDATE wrong_count=wrong_count+1, resolved=0, last_wrong_at=VALUES(last_wrong_at)`,
		childID, subject, refID, now,
	)
	if err != nil {
		// SQLite 方言回退
		_, err = s.DB.Exec(
			`INSERT INTO wrong_questions(child_id,subject,ref_id,wrong_count,resolved,last_wrong_at) VALUES(?,?,?,1,0,?)
			 ON CONFLICT(child_id,subject,ref_id) DO UPDATE SET wrong_count=wrong_count+1, resolved=0, last_wrong_at=excluded.last_wrong_at`,
			childID, subject, refID, now,
		)
		if err != nil {
			gLog("错题登记失败 child=%d %s#%d: %v", childID, subject, refID, err)
			return
		}
	}
}

// resolveWrong 巩固练习中答对后消除错题。
func (s *Store) resolveWrong(childID int, subject string, refID int) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if _, err := s.DB.Exec(
		"UPDATE wrong_questions SET resolved=1, last_reviewed_at=? WHERE child_id=? AND subject=? AND ref_id=?",
		now, childID, subject, refID,
	); err != nil {
		gLog("错题消除失败 child=%d %s#%d: %v", childID, subject, refID, err)
	}
}

// wrongIDs 学生当前未消除的错题 id 列表（按错误次数倒序，优先巩固错得多的）。
func (s *Store) wrongIDs(childID int, subject string, limit int) []int {
	if limit <= 0 {
		return nil
	}
	var ids []int
	if err := s.DB.Select(&ids,
		"SELECT ref_id FROM wrong_questions WHERE child_id=? AND subject=? AND resolved=0 ORDER BY wrong_count DESC, last_wrong_at DESC LIMIT ?",
		childID, subject, limit,
	); err != nil {
		return nil
	}
	return ids
}

// shouldReview 判断本题是否来自错题本（用于响应标记与答对消除）。
func isReviewRef(reviewRefs map[int]bool, id int) bool { return reviewRefs[id] }

// MixInWrongQuestions 从候选题目中按间隔混入错题，返回混合后的题目 id 顺序。
// 规则：每 poolEvery 道新题插入 1 道错题（至少 1 道、最多 wrongMax 道），错题不足时有多少插多少。
func MixInWrongQuestions(freshIDs []int, wrongIDs []int, poolEvery int, wrongMax int) []int {
	if len(wrongIDs) == 0 {
		return freshIDs
	}
	if poolEvery <= 0 {
		poolEvery = 2
	}
	if wrongMax <= 0 || wrongMax > len(wrongIDs) {
		wrongMax = len(wrongIDs)
	}
	mixed := make([]int, 0, len(freshIDs)+wrongMax)
	wi := 0
	for i, id := range freshIDs {
		mixed = append(mixed, id)
		// 每做完 poolEvery 道新题，插入一道错题巩固
		if (i+1)%poolEvery == 0 && wi < wrongMax {
			mixed = append(mixed, wrongIDs[wi])
			wi++
		}
	}
	// 剩余错题补到末尾
	for ; wi < wrongMax; wi++ {
		mixed = append(mixed, wrongIDs[wi])
	}
	return mixed
}

// reviewRefs 构建本题集合中的错题标记（服务端自行判断，避免信任前端）。
func (s *Store) reviewRefs(r *ghttp.Request, subject string, refIDs []int) map[int]bool {
	cid := s.resolveChild(r)
	if cid < 0 || len(refIDs) == 0 {
		return nil
	}
	q := "SELECT ref_id FROM wrong_questions WHERE child_id=? AND subject=? AND resolved=0 AND ref_id IN ("
	args := []interface{}{cid, subject}
	for i, id := range refIDs {
		if i > 0 {
			q += ","
		}
		q += "?"
		args = append(args, id)
	}
	q += ")"
	var ids []int
	if err := s.DB.Select(&ids, q, args...); err != nil {
		return nil
	}
	m := make(map[int]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}

// ---------- 每周经验排行榜接口 ----------

// WeeklyLeaderboard 周榜：GET /api/leaderboard/weekly?limit=20
// 返回当前 ISO 周经验值最高的学生名单。
func (s *Store) WeeklyLeaderboard(r *ghttp.Request) {
	limit := r.GetQuery("limit").Int()
	week := leaderboard.WeekKey(time.Now())
	entries := s.Board.Top(r.Context(), week, limit, func(id int) (string, string) {
		var name, avatar string
		_ = s.DB.Get(&name, "SELECT name FROM children WHERE id=?", id)
		_ = s.DB.Get(&avatar, "SELECT avatar FROM children WHERE id=?", id)
		return name, avatar
	})
	// 学生自己的排名（可能不在前 N）
	cid := s.resolveChild(r)
	var myXP, myRank int
	_ = s.DB.Get(&myXP, "SELECT COALESCE(xp,0) FROM children WHERE id=?", cid)
	for _, e := range entries {
		if e.ChildID == cid {
			myRank = e.Rank
			break
		}
	}
	s.ok(r, map[string]interface{}{
		"week":     week,
		"redis":    s.Board.Enabled(),
		"entries":  entries,
		"my_xp":    myXP,
		"my_rank":  myRank,
		"my_id":    cid,
	})
}

// ---------- 错题本接口 ----------

// ListWrongQuestions 学生错题本：GET /api/wrong-questions?subject=
func (s *Store) ListWrongQuestions(r *ghttp.Request) {
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	subject := r.GetQuery("subject").String()
	q := "SELECT id,child_id,subject,ref_id,wrong_count,resolved,last_wrong_at,COALESCE(last_reviewed_at,'') AS last_reviewed_at FROM wrong_questions WHERE child_id=?"
	args := []interface{}{cid}
	if subject != "" {
		q += " AND subject=?"
		args = append(args, subject)
	}
	q += " ORDER BY resolved, wrong_count DESC, last_wrong_at DESC LIMIT 200"
	var ws []map[string]interface{}
	rows, err := s.DB.Queryx(q, args...)
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := map[string]interface{}{}
		if err := rows.MapScan(m); err != nil {
			s.fail(r, 500, err.Error())
			return
		}
		// MySQL 驱动把 VARCHAR 扫成 []byte，这里统一转字符串（避免 JSON 输出 base64）
		for k, v := range m {
			if b, isBytes := v.([]byte); isBytes {
				m[k] = string(b)
			}
		}
		ws = append(ws, m)
	}
	if ws == nil {
		ws = []map[string]interface{}{}
	}
	s.ok(r, ws)
}
