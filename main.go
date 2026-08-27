package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/net/ghttp"

	"studyplanet/internal/config"
	"studyplanet/internal/db"
	"studyplanet/internal/handler"
	"studyplanet/internal/middleware"
	"studyplanet/internal/seed"
)

func main() {
	cfg := config.Load()

	// 确保 SQLite 文件所在目录存在
	if dir := filepath.Dir(cfg.Database.DSN); dir != "" && dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Fatalf("创建数据目录失败: %v", err)
		}
	}

	// 启动即跑迁移（golang-migrate 幂等）
	if err := db.Migrate(cfg.Database.DSN); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	sqlDB, err := db.Open(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	defer sqlDB.Close()

	// 空库写入五年级示例数据
	if cfg.Seed.Enabled {
		if err := seed.Run(sqlDB, cfg.Parent.Pin); err != nil {
			log.Printf("种子数据警告: %v", err)
		}
	}

	store := handler.NewStore(sqlDB, cfg)

	s := ghttp.GetServer()
	s.SetPort(cfg.Server.Port)
	s.SetDumpRouterMap(true)

	s.Group("/api", func(g *ghttp.RouterGroup) {
		g.Middleware(middleware.CORS)

		// 健康检查 & 登录相关（无需鉴权）
		g.GET("/health", store.Health)
		g.GET("/parent/auth-mode", store.AuthMode)
		g.POST("/parent/login", store.ParentLogin)
		g.GET("/parent/casdoor/login", store.CasdoorLogin)
		g.GET("/parent/casdoor/callback", store.CasdoorCallback)

		// 孩子端开放接口

		// 孩子端开放接口
		g.GET("/students", store.ListStudents)
		g.POST("/sessions", store.CreateSession)
		g.GET("/sessions", store.ListSessions)
		g.POST("/sessions/:id/finish", store.FinishSession)
		g.GET("/words", store.ListWords)
		g.GET("/words/:id", store.WordDetail)
		g.POST("/words/:id/progress", store.WordProgress)
		g.GET("/readings/:id", store.ReadingDetail)
		g.POST("/readings/:id/answer", store.ReadingAnswer)
		g.GET("/math", store.ListMath)
		g.POST("/math/:id/answer", store.MathAnswer)
		g.GET("/tasks", store.ListTasks)
		g.POST("/tasks/:id/complete", store.CompleteTask)
		g.GET("/points", store.PointsSummary)
		g.GET("/points/log", store.PointsLog)
		g.GET("/rewards", store.ListRewards)
		g.POST("/rewards/:id/redeem", store.Redeem)

		// 家长端（JWT 鉴权）
		g.Group("/parent", func(pg *ghttp.RouterGroup) {
			pg.Middleware(middleware.ParentAuth(cfg.Parent.JWTSecret))
			pg.POST("/tasks", store.AddTask)
			pg.DELETE("/tasks/:id", store.DeleteTask)
			pg.POST("/rewards", store.AddReward)
			pg.POST("/redemptions/:id/confirm", store.ConfirmRedemption)
			pg.POST("/set-pin", store.SetPin)
			pg.POST("/students", store.CreateStudent)
			pg.PUT("/students/:id", store.UpdateStudent)
			pg.DELETE("/students/:id", store.DeleteStudent)
		})
	})

	// 站点级路由：首页/学习页/logo（不在 /api 组内）
	s.BindHandler("GET:/", store.Index)
	s.BindHandler("GET:/app", store.Index)
	s.BindHandler("GET:/assets/logo.svg", store.Logo)

	log.Printf("学霸星球 StudyPlanet 服务端已启动，监听 :%d", cfg.Server.Port)
	s.Run()
}
