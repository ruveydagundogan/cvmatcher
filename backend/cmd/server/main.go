package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/application/iam/usecase"
	llmusecase "github.com/ruveydagundogan/llm-decision-score/backend/internal/application/llmscoring/usecase"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/auth"
	backendllmhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/backendllm"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/health"
	iamhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/iam"
	llmhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/llmscoring"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/router"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/llm"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/memory"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/metrics"
	pgresaudit "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/postgres/audit"
	pgresiam "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/postgres/iam"
	pgresscoring "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/postgres/llmscoring"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/config"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/database"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/logger"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/middleware"

	auditrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/audit/repository"
	iamrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/iam/repository"
	scoringrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/llmscoring/repository"
)

func main() {
	cfg := config.Load()

	if cfg.JWT.Secret == "" {
		env := os.Getenv("GO_ENV")
		if env == "production" || env == "staging" {
			slog.Error("JWT_SECRET is required in production/staging")
			os.Exit(1)
		}
		slog.Warn("JWT_SECRET not set, using development fallback")
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Format)

	log.Info("starting llm-decision-score backend",
		"server_host", cfg.Server.Host,
		"server_port", cfg.Server.Port,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		userRepo    iamrepo.UserRepository
		roleRepo    iamrepo.RoleRepository
		scoringRepo scoringrepo.ScoringRepository
		auditRepo   auditrepo.AuditRepository
		dbPool      interface{ Close() }
	)

	pool, err := database.NewPostgresPool(ctx, cfg.Database, log)
	if err != nil {
		log.Warn("postgresql unavailable, running in in-memory mode", "error", err)
		memUser := memory.NewInMemoryUserRepo()
		userRepo = memUser
		roleRepo = memory.NewInMemoryRoleRepo()
		scoringRepo = memory.NewInMemoryScoringRepo()
		auditRepo = memory.NewInMemoryAuditRepo()
	} else {
		log.Info("postgresql connected, using persistent storage")
		if err := database.RunMigrations(ctx, pool, log); err != nil {
			log.Error("failed to run migrations", "error", err)
		}
		userRepo = pgresiam.NewUserRepository(pool)
		roleRepo = pgresiam.NewRoleRepository(pool)
		scoringRepo = pgresscoring.NewScoringRepository(pool)
		auditRepo = pgresaudit.NewAuditRepository(pool)
		dbPool = pool
	}

	jwtService := auth.NewJWTService(cfg.JWT)
	bcryptCost := cfg.BcryptCost
	if bcryptCost == 0 {
		bcryptCost = 10
	}
	bcryptService := auth.NewBcryptAuthService(bcryptCost)

	registerUC := usecase.NewRegisterUseCase(userRepo, roleRepo, bcryptService, jwtService, auditRepo, log)
	loginUC := usecase.NewLoginUseCase(userRepo, jwtService, auditRepo, log)
	getProfileUC := usecase.NewGetProfileUseCase(userRepo, log)
	updateProfileUC := usecase.NewUpdateProfileUseCase(userRepo, log)

	scoreUC := llmusecase.NewScoreUseCase(scoringRepo, auditRepo, log)
	historyUC := llmusecase.NewGetHistoryUseCase(scoringRepo, log)
	deleteHistoryUC := llmusecase.NewDeleteHistoryUseCase(scoringRepo, auditRepo, log)
	statsUC := llmusecase.NewGetStatsUseCase(scoringRepo, log)

	healthHandler := health.NewHandler(nil)
	iamH := iamhandler.NewHandler(registerUC, loginUC, getProfileUC, updateProfileUC)
	llmH := llmhandler.NewHandler(scoreUC, historyUC, deleteHistoryUC, statsUC)

	rateLimiter := middleware.NewRateLimiter(100, 200, log)

	var backendLLMHandler *backendllmhandler.Handler
	if cfg.MLCLLM.Enabled {
		llmClient := llm.NewClient(cfg.MLCLLM.BaseURL, 60*time.Second)
		backendLLMHandler = backendllmhandler.NewHandler(llmClient)
		log.Info("server-side LLM enabled", "base_url", cfg.MLCLLM.BaseURL)
	} else {
		log.Info("server-side LLM disabled, using browser-based WebLLM")
	}

	httpMetrics := metrics.NewHTTPMetrics("llm-decision-score")

	deps := router.Dependencies{
		HealthHandler:    healthHandler,
		IAMHandler:       iamH,
		LLMHandler:       llmH,
		BackendLLMHandler: backendLLMHandler,
		JWTValidator:     jwtService,
		Config:           cfg,
		RateLimiter:      rateLimiter,
		Metrics:          httpMetrics,
	}

	r := router.NewRouter(deps)

	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("server listening", "addr", cfg.Server.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	if dbPool != nil {
		dbPool.Close()
	}
	log.Info("server stopped")
}
