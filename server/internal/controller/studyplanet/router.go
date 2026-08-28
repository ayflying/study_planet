package studyplanet

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"studyplanet/internal/config"
	"studyplanet/internal/middleware"
)

// BindRoutes 集中维护对外 HTTP 路由。
// 路由由 api 定义（api/studyplanet/v1 的 g.Meta）驱动，group.Bind(ctrl) 自动注册；
// 分组规则：
//   - /api/students*、/api/tasks*、/api/rewards*、/api/points*、/api/learn* 等学生接口公开访问；
//   - /api/parent/* 家长专属操作统一挂 ParentAuth（JWT）鉴权。
//     注意 auth-mode/login/casdoor 回调虽在 /api/parent 路径下但必须匿名访问，
//     因此拆为「匿名 parent 认证组」与「鉴权 parent 管理组」两个子分组。
func BindRoutes(server *ghttp.Server, cfg *config.Config) {
	ctrl := NewV1()
	server.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)

		// ---- 公开：系统与学生端 ----
		group.Bind(
			ctrl.Health,
			ctrl.AuthMode,
			ctrl.ParentLogin,
			ctrl.CasdoorLogin,
			ctrl.CasdoorCallback,
			ctrl.ListStudents,
			ctrl.ListTasks,
			ctrl.CompleteTask,
			ctrl.ListRewards,
			ctrl.Redeem,
			ctrl.PointsSummary,
			ctrl.PointsLog,
			ctrl.ListWords,
			ctrl.WordDetail,
			ctrl.WordProgress,
			ctrl.ReadingDetail,
			ctrl.ReadingAnswer,
			ctrl.ListMath,
			ctrl.MathAnswer,
			ctrl.CreateSession,
			ctrl.ListSessions,
			ctrl.FinishSession,
			ctrl.ListSubjects,
			ctrl.PickQuestions,
			ctrl.ContentItem,
			ctrl.ContentAnswer,
			ctrl.WeeklyLeaderboard,
			ctrl.ListWrongQuestions,
		)

		// ---- 家长鉴权：管理操作 ----
		group.Group("/", func(parent *ghttp.RouterGroup) {
			parent.Middleware(middleware.ParentAuth(cfg.Parent.JWTSecret))
			parent.Bind(
				ctrl.SetPin,
				ctrl.CreateStudent,
				ctrl.UpdateStudent,
				ctrl.DeleteStudent,
				ctrl.AddTask,
				ctrl.DeleteTask,
				ctrl.AddReward,
				ctrl.ConfirmRedemption,
				ctrl.ImportContent,
				ctrl.SubjectStats,
			)
		})
	})
}
