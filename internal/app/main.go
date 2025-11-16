package app

import (
	"pull-request-review/config"
	"pull-request-review/internal/delivery/http/handlers"
	"pull-request-review/internal/infrastructure/adapters/logger"
	"pull-request-review/internal/infrastructure/adapters/router"
	"pull-request-review/internal/infrastructure/database"
	"pull-request-review/internal/infrastructure/http/route"
	"pull-request-review/internal/infrastructure/http/server"
	"pull-request-review/internal/infrastructure/repository"
	"pull-request-review/internal/service"
)

type Application struct {
	router router.Router
}

func Run(cfg *config.Config, db *database.Database, appLogger logger.Logger) {
	app := initializeApp(db, cfg, appLogger)

	srv := server.NewServer(app.router, cfg.Server, appLogger)

	srv.Start()
	srv.WaitForShutdown()
}

func initializeApp(db *database.Database, cfg *config.Config, appLogger logger.Logger) *Application {
	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	prRepo := repository.NewPullRequestRepositoryPgx(db)
	reviewAssignmentRepo := repository.NewReviewAssignmentRepository(db)

	teamService := service.NewTeamService(teamRepo, userRepo, appLogger)
	userService := service.NewUserService(userRepo, teamRepo, appLogger)
	prService := service.NewPullRequestService(
		prRepo,
		userRepo,
		reviewAssignmentRepo,
		appLogger,
		cfg.Service.MaxReviewersCount,
	)
	statisticsService := service.NewStatisticsService(reviewAssignmentRepo, prRepo, appLogger)

	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService, prService)
	prHandler := handlers.NewPullRequestHandler(prService)
	healthHandler := handlers.NewHealthHandler(db)
	statsHandler := handlers.NewStatisticsHandler(statisticsService)

	r := router.NewGinRouter()
	route.SetupRoutes(
		r, &route.Handlers{
			TeamHandler:        teamHandler,
			UserHandler:        userHandler,
			PullRequestHandler: prHandler,
			StatisticsHandler:  statsHandler,
			HealthHandler:      healthHandler,
		},
		appLogger,
		cfg.Server.RequestTimeout,
	)

	return &Application{
		router: r,
	}
}