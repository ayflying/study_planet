package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"

	"studyplanet/internal/config"
	studyplanetcontroller "studyplanet/internal/controller/studyplanet"
	"studyplanet/internal/db"
	"studyplanet/internal/seed"
	studyplanetservice "studyplanet/internal/service/studyplanet"
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
		if strings.EqualFold(cfg.Database.Driver, "sqlite") || strings.EqualFold(cfg.Database.Driver, "sqlite3") {
			if dir := filepath.Dir(cfg.Database.DSN); dir != "" && dir != "." && dir != "/" {
				if err := os.MkdirAll(dir, 0o755); err != nil {
					return err
				}
			}
		}
		if err := db.Migrate(cfg.Database.Driver, cfg.Database.DSN); err != nil {
			return err
		}

		sqlDB, err := db.Open(cfg.Database.Driver, cfg.Database.DSN)
		if err != nil {
			return err
		}
		defer sqlDB.Close()

		if cfg.Seed.Enabled {
			if err = seed.Run(sqlDB, cfg.Parent.Pin); err != nil {
				log.Printf("种子数据警告: %v", err)
			}
		}

		store := studyplanetservice.NewStore(sqlDB, cfg)
		s := g.Server()
		s.SetPort(cfg.Server.Port)
		s.SetDumpRouterMap(true)
		if staticRoot := resolveStaticRoot(); staticRoot != "" {
			s.SetServerRoot(staticRoot)
			s.SetIndexFiles([]string{"index.html"})
			s.SetFileServerEnabled(true)
		}
		studyplanetcontroller.BindRoutes(s, store, cfg)

		log.Printf("学霸星球 StudyPlanet 服务端已启动，监听 :%d", cfg.Server.Port)
		s.Run()
		return nil
	},
}
