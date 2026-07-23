package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cvusecase "github.com/ruveydagundogan/llm-decision-score/backend/internal/application/cv/usecase"
	iusecase "github.com/ruveydagundogan/llm-decision-score/backend/internal/application/iam/usecase"
	jdusecase "github.com/ruveydagundogan/llm-decision-score/backend/internal/application/jobdescription/usecase"
	llmusecase "github.com/ruveydagundogan/llm-decision-score/backend/internal/application/llmscoring/usecase"
	matchusecase "github.com/ruveydagundogan/llm-decision-score/backend/internal/application/matching/usecase"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/auth"
	cvhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/cv"
	backendllmhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/backendllm"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/health"
	iamhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/iam"
	jdhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/jd"
	llmhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/llmscoring"
	matchinghandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/matching"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/router"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/llm"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/memory"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/metrics"
	pgresaudit "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/postgres/audit"
	pgrescv "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/postgres/cv"
	pgresiam "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/postgres/iam"
	pgresjd "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/postgres/jobdescription"
	pgresscoring "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/postgres/llmscoring"
	pgresmatch "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/postgres/matching"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/ratelimit"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/tunnel"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/cache"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/config"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/database"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/logger"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/middleware"

	auditrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/audit/repository"
	cvrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/cv/repository"
	iamrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/iam/repository"
	jdrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/jobdescription/repository"
	matchrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/matching/repository"
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
		userRepo     iamrepo.UserRepository
		roleRepo     iamrepo.RoleRepository
		scoringRepo  scoringrepo.ScoringRepository
		auditRepo    auditrepo.AuditRepository
		cvRepo       cvrepo.CVRepository
		jdRepo       jdrepo.JobDescriptionRepository
		matchingRepo matchrepo.MatchingRepository
		dbPool       interface{ Close() }
	)

	pool, err := database.NewPostgresPool(ctx, cfg.Database, log)
	if err != nil {
		log.Warn("postgresql unavailable, running in in-memory mode", "error", err)
		memUser := memory.NewInMemoryUserRepo()
		userRepo = memUser
		roleRepo = memory.NewInMemoryRoleRepo()
		scoringRepo = memory.NewInMemoryScoringRepo()
		auditRepo = memory.NewInMemoryAuditRepo()
		cvRepo = memory.NewInMemoryCVRepo()
		jdRepo = memory.NewInMemoryJDRepo()
		matchingRepo = memory.NewInMemoryMatchingRepo()
	} else {
		log.Info("postgresql connected, using persistent storage")
		if err := database.RunMigrations(ctx, pool, log); err != nil {
			log.Error("failed to run migrations", "error", err)
		}
		userRepo = pgresiam.NewUserRepository(pool)
		roleRepo = pgresiam.NewRoleRepository(pool)
		scoringRepo = pgresscoring.NewScoringRepository(pool)
		auditRepo = pgresaudit.NewAuditRepository(pool)
		cvRepo = pgrescv.NewCVRepository(pool)
		jdRepo = pgresjd.NewJobDescriptionRepository(pool)
		matchingRepo = pgresmatch.NewMatchingRepository(pool)
		dbPool = pool
	}

	jwtService := auth.NewJWTService(cfg.JWT)
	bcryptCost := cfg.BcryptCost
	if bcryptCost == 0 {
		bcryptCost = 10
	}
	bcryptService := auth.NewBcryptAuthService(bcryptCost)

	registerUC := iusecase.NewRegisterUseCase(userRepo, roleRepo, bcryptService, jwtService, auditRepo, log)
	loginUC := iusecase.NewLoginUseCase(userRepo, jwtService, auditRepo, log)
	getProfileUC := iusecase.NewGetProfileUseCase(userRepo, log)
	updateProfileUC := iusecase.NewUpdateProfileUseCase(userRepo, log)

	scoreUC := llmusecase.NewScoreUseCase(scoringRepo, auditRepo, log)
	historyUC := llmusecase.NewGetHistoryUseCase(scoringRepo, log)
	deleteHistoryUC := llmusecase.NewDeleteHistoryUseCase(scoringRepo, auditRepo, log)
	statsUC := llmusecase.NewGetStatsUseCase(scoringRepo, log)

	redisClient, err := cache.NewRedisClient(ctx, cfg.Redis, log)
	if err != nil {
		log.Warn("redis unavailable, running with in-memory rate limiter", "error", err)
		redisClient = nil
	}

	m := metrics.New("llm-decision-score")

	healthHandler := health.NewHandler(nil)
	iamH := iamhandler.NewHandler(registerUC, loginUC, getProfileUC, updateProfileUC)
	llmH := llmhandler.NewHandler(scoreUC, historyUC, deleteHistoryUC, statsUC, m)

	var rateLimiter func(http.Handler) http.Handler
	if redisClient != nil {
		log.Info("redis connected, using distributed rate limiter")
		rateLimiter = ratelimit.NewRedisRateLimiter(redisClient, 100, 200, log).RateLimit
	} else {
		log.Info("redis unavailable, using in-memory rate limiter")
		rateLimiter = middleware.NewRateLimiter(100, 200, log).RateLimit
	}

	var llmClient *llm.Client
	baseURL := cfg.MLCLLM.BaseURL
	if baseURL == "" {
		baseURL = "/tunnel"
	}
	if baseURL[0] == '/' {
		port := os.Getenv("PORT")
		if port == "" {
			port = fmt.Sprintf("%d", cfg.Server.Port)
		}
		baseURL = fmt.Sprintf("http://localhost:%s%s", port, baseURL)
	}
	llmClient = llm.NewClient(baseURL, 60*time.Second)
	log.Info("server-side LLM enabled", "base_url", baseURL)
	tunnelServer := tunnel.NewServer()
	backendLLMHandler := backendllmhandler.NewHandler(llmClient, m)

	cvUseCase := cvusecase.NewCVUseCase(cvRepo, llmClient, log)
	jdUseCase := jdusecase.NewJDUseCase(jdRepo, llmClient, log)
	matchUseCase := matchusecase.NewMatchUseCase(cvRepo, jdRepo, matchingRepo, llmClient, log)

	cvH := cvhandler.NewHandler(cvUseCase)
	jdH := jdhandler.NewHandler(jdUseCase)
	matchH := matchinghandler.NewHandler(matchUseCase)

	deps := router.Dependencies{
		HealthHandler:      healthHandler,
		IAMHandler:         iamH,
		LLMHandler:         llmH,
		BackendLLMHandler:  backendLLMHandler,
		CVHandler:          cvH,
		JDHandler:          jdH,
		MatchingHandler:    matchH,
		JWTValidator:       jwtService,
		Config:             cfg,
		RateLimiter:        rateLimiter,
		Metrics:            m,
	}

	r := router.NewRouter(deps)

	mux := http.NewServeMux()
	if tunnelServer != nil {
		mux.HandleFunc("/tunnel/ws", tunnelServer.HandleWebSocket)
		mux.HandleFunc("/tunnel/", tunnelServer.ProxyHandler)
	}
	mux.Handle("/", r)

	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      mux,
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
