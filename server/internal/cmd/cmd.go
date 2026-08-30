package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"

	"studyplanet/internal/config"
	"studyplanet/internal/battle"
	studyplanetcontroller "studyplanet/internal/controller/studyplanet"
	"studyplanet/internal/db"
	"studyplanet/internal/gdbinit"
	"studyplanet/internal/leaderboard"
	"studyplanet/internal/logic/studyplanet"
	"studyplanet/internal/seed"
	"studyplanet/internal/seedcontent"
)

func resolveStaticRoot() string {
	candidates := []string{"client", filepath.Join("..", "client", "dist")}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	return ""
}

var Main = gcmd.Command{
	Name:  "main",
	Usage: "main",
	Brief: "启动 StudyPlanet HTTP API 服务",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		cfg := config.Load()
		// 数据库固定使用 MySQL：先迁移表结构，再注册 ORM 数据源
		if err := db.Migrate(cfg.Database.DSN); err != nil {
			return err
		}

		if err := gdbinit.Setup(cfg.Database.DSN); err != nil {
			return err
		}
		orm := g.DB()

		if cfg.Seed.Enabled {
			if err := seed.Run(ctx, orm); err != nil {
				log.Printf("种子数据警告: %v", err)
			}
		}

		// 动态内容库：同步学科目录，空题库时导入内置全科题（之后以数据库为准）
		if err := seedcontent.Run(ctx, orm); err != nil {
			log.Printf("内容库初始化警告: %v", err)
		}

		// 每周经验排行榜：Redis 实时 + 每小时持久化到数据库
		board := leaderboard.New(orm, os.Getenv("REDIS_ADDR"))
		go func() {
			ticker := time.NewTicker(time.Hour)
			defer ticker.Stop()
			for range ticker.C {
				if err := board.PersistSnapshot(context.Background()); err != nil {
					log.Printf("周榜持久化失败: %v", err)
				}
			}
		}()

		// 依赖注入：业务实现（logic 层 sStudyPlanet）需要配置与排行榜模块
		studyplanet.SetDeps(cfg, board)
		s := g.Server()
		s.SetPort(cfg.Server.Port)
		s.SetDumpRouterMap(true)
		if staticRoot := resolveStaticRoot(); staticRoot != "" {
			s.SetServerRoot(staticRoot)
			s.SetIndexFiles([]string{"index.html"})
			s.SetFileServerEnabled(true)
		}
		studyplanetcontroller.BindRoutes(s, cfg)

		// 真人对战 WebSocket：/ws/battle（gorilla 升级，独立于 gf 分组中间件）
		battleEngine := battle.New()
		battleEngine.AddXP = studyplanet.ExternalAddXP(board)
		battleEngine.OnSnack = studyplanet.ExternalSnackDrop()
		s.BindHandler("/ws/battle", battleEngine.HandleWS)

		log.Printf("学霸星球 StudyPlanet 服务端已启动，监听 :%d", cfg.Server.Port)
		s.Run()
		return nil
	},
}
