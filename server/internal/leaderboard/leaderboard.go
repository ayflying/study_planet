// Package leaderboard 独立的每周经验值排行榜模块。
//
// 设计：
//   - 实时排名走 Redis ZSET（key 按周切换，如 rank:2026W35，一周后自动过期）；
//   - 每小时把 ZSET 全量快照持久化到 MySQL leaderboard_weekly 表，防 Redis 数据丢失；
//   - Redis 不可用时自动降级：XP 直接累写数据库，查询从持久化表读取，功能不中断；
//   - 对上层只暴露 AddXP / Top / WeekKey / PersistSnapshot，屏蔽存储细节。
package leaderboard

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// Entry 排行榜中的一条排名记录。
type Entry struct {
	Rank    int    `json:"rank"`
	ChildID int    `json:"child_id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	XP      int    `json:"xp"`
}

type Board struct {
	db      *sqlx.DB
	rdb     *redis.Client
	enabled bool // Redis 是否可用
	mu      sync.RWMutex
}

// New 创建排行榜；addr 为空时进入降级模式（纯数据库）。
func New(db *sqlx.DB, addr string) *Board {
	b := &Board{db: db}
	if strings.TrimSpace(addr) != "" {
		b.rdb = redis.NewClient(&redis.Options{Addr: addr, DialTimeout: 2 * time.Second})
		if err := b.rdb.Ping(context.Background()).Err(); err != nil {
			log.Printf("排行榜: Redis 连接失败，降级为数据库模式: %v", err)
		} else {
			b.enabled = true
		}
	} else {
		log.Printf("排行榜: 未配置 REDIS_ADDR，使用数据库模式")
	}
	return b
}

// Enabled 报告 Redis 是否可用。
func (b *Board) Enabled() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.enabled
}

// WeekKey 计算某时刻所属 ISO 周的 key（周一为一周起点），如 2026W35。
func WeekKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%dW%02d", year, week)
}

func redisKey(week string) string { return "rank:" + week }

// AddXP 给学生累加本周经验值：先写持久化表（保证不丢），Redis 可用时再写 ZSET。
func (b *Board) AddXP(ctx context.Context, childID int, delta int) {
	if delta == 0 {
		return
	}
	week := WeekKey(time.Now())
	if err := b.upsertXP(ctx, week, childID, delta, false); err != nil {
		log.Printf("排行榜: 持久化 XP 失败 child=%d: %v", childID, err)
	}
	if !b.Enabled() {
		return
	}
	key := redisKey(week)
	if err := b.rdb.ZIncrBy(ctx, key, float64(delta), strconv.Itoa(childID)).Err(); err != nil {
		log.Printf("排行榜: Redis ZIncrBy 失败（不影响主流程）: %v", err)
		return
	}
	_ = b.rdb.Expire(ctx, key, 8*24*time.Hour).Err() // 当周结束 7 天后自动清理
}

// upsertXP 持久化表写入：full=true 时直接覆盖（快照），否则累加。兼容 MySQL/SQLite。
func (b *Board) upsertXP(ctx context.Context, week string, childID int, delta int, full bool) error {
	_, err := b.db.ExecContext(ctx,
		"INSERT INTO leaderboard_weekly(week_key, child_id, xp) VALUES(?,?,?) "+
			"ON DUPLICATE KEY UPDATE xp = VALUES(xp)",
		week, childID, delta,
	)
	if err == nil {
		return nil
	}
	if full {
		_, err = b.db.ExecContext(ctx,
			"INSERT INTO leaderboard_weekly(week_key, child_id, xp) VALUES(?,?,?) "+
				"ON CONFLICT(week_key, child_id) DO UPDATE SET xp = excluded.xp",
			week, childID, delta,
		)
		return err
	}
	_, err = b.db.ExecContext(ctx,
		"INSERT INTO leaderboard_weekly(week_key, child_id, xp) VALUES(?,?,?) "+
			"ON CONFLICT(week_key, child_id) DO UPDATE SET xp = xp + excluded.xp",
		week, childID, delta,
	)
	return err
}

// Top 返回周榜前 limit 名。Redis 不可用时从持久化表读取。
func (b *Board) Top(ctx context.Context, week string, limit int, nameOf func(int) (string, string)) []Entry {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	type kv struct {
		id int
		xp int
	}
	var rows []kv
	if b.Enabled() {
		vals, err := b.rdb.ZRevRangeWithScores(ctx, redisKey(week), 0, int64(limit-1)).Result()
		if err == nil {
			for _, v := range vals {
				id, _ := strconv.Atoi(v.Member.(string))
				rows = append(rows, kv{id, int(v.Score)})
			}
		} else {
			log.Printf("排行榜: Redis 读取失败，回退数据库: %v", err)
		}
	}
	if rows == nil {
		type dbRow struct {
			ChildID int `db:"child_id"`
			XP      int `db:"xp"`
		}
		var dbRows []dbRow
		if err := b.db.SelectContext(ctx, &dbRows,
			"SELECT child_id, xp FROM leaderboard_weekly WHERE week_key=? ORDER BY xp DESC, child_id LIMIT ?",
			week, limit,
		); err != nil {
			log.Printf("排行榜: 数据库读取失败: %v", err)
			return []Entry{}
		}
		for _, r := range dbRows {
			rows = append(rows, kv{r.ChildID, r.XP})
		}
	}
	entries := make([]Entry, 0, len(rows))
	for _, r := range rows {
		name, avatar := nameOf(r.id)
		entries = append(entries, Entry{Rank: len(entries) + 1, ChildID: r.id, Name: name, Avatar: avatar, XP: r.xp})
	}
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].XP > entries[j].XP })
	for i := range entries {
		entries[i].Rank = i + 1
	}
	return entries
}

// PersistSnapshot 把当前周榜从 Redis 全量快照写入持久化表（每小时由 cmd 定时调用）。
func (b *Board) PersistSnapshot(ctx context.Context) error {
	if !b.Enabled() {
		return nil
	}
	week := WeekKey(time.Now())
	vals, err := b.rdb.ZRangeWithScores(ctx, redisKey(week), 0, -1).Result()
	if err != nil {
		return fmt.Errorf("读取 Redis 周榜失败: %w", err)
	}
	for _, v := range vals {
		id, _ := strconv.Atoi(v.Member.(string))
		if err := b.upsertXP(ctx, week, id, int(v.Score), true); err != nil {
			return fmt.Errorf("持久化 child=%d 失败: %w", id, err)
		}
	}
	log.Printf("排行榜: 已持久化 %d 条周榜记录（%s）", len(vals), week)
	return nil
}
