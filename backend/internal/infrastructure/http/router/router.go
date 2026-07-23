package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	backendllmhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/backendllm"
	healthhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/health"
	iamhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/iam"
	llmhandler "github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/http/handler/llmscoring"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/metrics"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/config"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/middleware"
)

type Dependencies struct {
	HealthHandler      *healthhandler.Handler
	IAMHandler         *iamhandler.Handler
	LLMHandler         *llmhandler.Handler
	BackendLLMHandler  *backendllmhandler.Handler
	JWTValidator       middleware.TokenValidator
	Config             *config.Config
	RateLimiter        func(http.Handler) http.Handler
	Metrics            *metrics.Metrics
}

func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(middleware.Logging(nil))
	r.Use(middleware.Recoverer(nil))
	r.Use(middleware.MaxBodyBytes(deps.Config.Server.MaxBodyBytes))
	r.Use(middleware.CORS(deps.Config.Server.CORSAllowedOrigins))

	r.Use(deps.Metrics.HTTP.Middleware)

	r.Get("/health/live", deps.HealthHandler.Live)
	r.Get("/health/ready", deps.HealthHandler.Ready)

	r.Handle("/metrics", deps.Metrics.Handler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", deps.IAMHandler.Register)
		r.Post("/auth/login", deps.IAMHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(deps.RateLimiter)
			r.Post("/score", deps.LLMHandler.Score)
		})

		r.Get("/models", deps.LLMHandler.GetModels)

		if deps.BackendLLMHandler != nil {
			r.Post("/llm/chat", deps.BackendLLMHandler.Chat)
		}

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(deps.JWTValidator))

			r.Get("/me", deps.IAMHandler.GetProfile)
			r.Put("/me", deps.IAMHandler.UpdateProfile)

			r.Get("/history", deps.LLMHandler.GetHistory)
			r.Delete("/history", deps.LLMHandler.DeleteHistory)
			r.Get("/stats", deps.LLMHandler.GetStats)
		})
	})

	return r
}
