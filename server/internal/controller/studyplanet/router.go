package studyplanet

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"studyplanet/internal/config"
	"studyplanet/internal/middleware"
	"studyplanet/internal/service"
)

// BindRoutes 集中维护对外 HTTP 路由；业务经 service 接口层解耦，实现位于 internal/logic。
func BindRoutes(server *ghttp.Server, cfg *config.Config) {
	logic := service.StudyPlanet()
	server.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)
		group.GET("/health", logic.Health)
		group.GET("/parent/auth-mode", logic.AuthMode)
		group.POST("/parent/login", logic.ParentLogin)
		group.GET("/parent/casdoor/login", logic.CasdoorLogin)
		group.GET("/parent/casdoor/callback", logic.CasdoorCallback)

		group.GET("/students", logic.ListStudents)
		group.POST("/sessions", logic.CreateSession)
		group.GET("/sessions", logic.ListSessions)
		group.POST("/sessions/:id/finish", logic.FinishSession)
		group.GET("/words", logic.ListWords)
		group.GET("/words/:id", logic.WordDetail)
		group.POST("/words/:id/progress", logic.WordProgress)
		group.GET("/readings/:id", logic.ReadingDetail)
		group.POST("/readings/:id/answer", logic.ReadingAnswer)
		group.GET("/math", logic.ListMath)
		group.POST("/math/:id/answer", logic.MathAnswer)
		group.GET("/tasks", logic.ListTasks)
		group.POST("/tasks/:id/complete", logic.CompleteTask)
		group.GET("/points", logic.PointsSummary)
		group.GET("/points/log", logic.PointsLog)
		group.GET("/rewards", logic.ListRewards)
		group.POST("/rewards/:id/redeem", logic.Redeem)
		group.GET("/subjects", logic.ListSubjects)
		group.GET("/content/pick", logic.PickQuestions)
		group.GET("/content/item", logic.ContentItem)
		group.POST("/content/answer", logic.ContentAnswer)
		group.GET("/wrong-questions", logic.ListWrongQuestions)
		group.GET("/leaderboard/weekly", logic.WeeklyLeaderboard)

		group.Group("/parent", func(parent *ghttp.RouterGroup) {
			parent.Middleware(middleware.ParentAuth(cfg.Parent.JWTSecret))
			parent.POST("/tasks", logic.AddTask)
			parent.DELETE("/tasks/:id", logic.DeleteTask)
			parent.POST("/rewards", logic.AddReward)
			parent.POST("/redemptions/:id/confirm", logic.ConfirmRedemption)
			parent.POST("/set-pin", logic.SetPin)
			parent.POST("/students", logic.CreateStudent)
			parent.PUT("/students/:id", logic.UpdateStudent)
			parent.DELETE("/students/:id", logic.DeleteStudent)
			parent.POST("/content/import", logic.ImportContent)
			parent.GET("/content/stats", logic.SubjectStats)
		})
	})
}
