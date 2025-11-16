package route

import (
	"net/http"
	"time"

	"pull-request-review/internal/delivery/http/handlers"
	"pull-request-review/internal/infrastructure/adapters/logger"
	"pull-request-review/internal/infrastructure/adapters/router"
	"pull-request-review/internal/infrastructure/http/middleware"
)

type Handlers struct {
	TeamHandler        *handlers.TeamHandler
	UserHandler        *handlers.UserHandler
	PullRequestHandler *handlers.PullRequestHandler
	StatisticsHandler  *handlers.StatisticsHandler
	HealthHandler      *handlers.HealthHandler
}

func SetupRoutes(r router.Router, handlers *Handlers, logger logger.Logger, requestTimeout time.Duration) {
	r.Use(
		middleware.Recovery(logger),
		middleware.Logger(logger),
		middleware.Timeout(requestTimeout),
	)

	r.GET("/health", http.HandlerFunc(handlers.HealthHandler.Check))

	teamGroup := r.Group("/team")
	teamGroup.POST("/add", http.HandlerFunc(handlers.TeamHandler.AddTeam))
	teamGroup.GET("/get", http.HandlerFunc(handlers.TeamHandler.GetTeam))

	userGroup := r.Group("/users")
	userGroup.POST("/setIsActive", http.HandlerFunc(handlers.UserHandler.SetIsActive))
	userGroup.GET("/getReview", http.HandlerFunc(handlers.UserHandler.GetReviews))

	prGroup := r.Group("/pullRequest")
	prGroup.POST("/create", http.HandlerFunc(handlers.PullRequestHandler.CreatePullRequest))
	prGroup.POST("/merge", http.HandlerFunc(handlers.PullRequestHandler.MergePullRequest))
	prGroup.POST("/reassign", http.HandlerFunc(handlers.PullRequestHandler.ReassignReviewer))

	r.GET("/statistics", http.HandlerFunc(handlers.StatisticsHandler.GetStatistics))
}