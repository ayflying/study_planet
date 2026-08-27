package studyplanet

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"studyplanet/internal/config"
	"studyplanet/internal/middleware"
	studyplanetservice "studyplanet/internal/service/studyplanet"
)

// BindRoutes 集中维护对外 HTTP 路由；业务实现位于 internal/service/studyplanet。
func BindRoutes(server *ghttp.Server, store *studyplanetservice.Store, cfg *config.Config) {
	server.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)
		group.GET("/health", store.Health)
		group.GET("/parent/auth-mode", store.AuthMode)
		group.POST("/parent/login", store.ParentLogin)
		group.GET("/parent/casdoor/login", store.CasdoorLogin)
		group.GET("/parent/casdoor/callback", store.CasdoorCallback)

		group.GET("/students", store.ListStudents)
		group.POST("/sessions", store.CreateSession)
		group.GET("/sessions", store.ListSessions)
		group.POST("/sessions/:id/finish", store.FinishSession)
		group.GET("/words", store.ListWords)
		group.GET("/words/:id", store.WordDetail)
		group.POST("/words/:id/progress", store.WordProgress)
		group.GET("/readings/:id", store.ReadingDetail)
		group.POST("/readings/:id/answer", store.ReadingAnswer)
		group.GET("/math", store.ListMath)
		group.POST("/math/:id/answer", store.MathAnswer)
		group.GET("/tasks", store.ListTasks)
		group.POST("/tasks/:id/complete", store.CompleteTask)
		group.GET("/points", store.PointsSummary)
		group.GET("/points/log", store.PointsLog)
		group.GET("/rewards", store.ListRewards)
		group.POST("/rewards/:id/redeem", store.Redeem)

		group.Group("/parent", func(parent *ghttp.RouterGroup) {
			parent.Middleware(middleware.ParentAuth(cfg.Parent.JWTSecret))
			parent.POST("/tasks", store.AddTask)
			parent.DELETE("/tasks/:id", store.DeleteTask)
			parent.POST("/rewards", store.AddReward)
			parent.POST("/redemptions/:id/confirm", store.ConfirmRedemption)
			parent.POST("/set-pin", store.SetPin)
			parent.POST("/students", store.CreateStudent)
			parent.PUT("/students/:id", store.UpdateStudent)
			parent.DELETE("/students/:id", store.DeleteStudent)
		})
	})
}
