package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	adminhandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/admin"
	backendllmhandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/backendllm"
	cvhandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/cv"
	healthhandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/health"
	iamhandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/iam"
	jdhandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/jd"
	knowledgehandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/knowledge"
	llmhandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/llmscoring"
	matchinghandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/matching"
	mcphandler "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/http/handler/mcp"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/metrics"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/config"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/middleware"
)

type Dependencies struct {
	HealthHandler      *healthhandler.Handler
	IAMHandler         *iamhandler.Handler
	LLMHandler         *llmhandler.Handler
	BackendLLMHandler  *backendllmhandler.Handler
	CVHandler          *cvhandler.Handler
	JDHandler          *jdhandler.Handler
	MatchingHandler    *matchinghandler.Handler
	MCPHandler         *mcphandler.Handler
	KnowledgeHandler   *knowledgehandler.Handler
	AdminHandler       *adminhandler.Handler
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

		if deps.MCPHandler != nil {
			r.Post("/mcp/query", deps.MCPHandler.Query)
			r.Get("/mcp/adapters", deps.MCPHandler.ListAdapters)
		}

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(deps.JWTValidator))

			r.Get("/me", deps.IAMHandler.GetProfile)
			r.Put("/me", deps.IAMHandler.UpdateProfile)

			r.Get("/history", deps.LLMHandler.GetHistory)
			r.Delete("/history", deps.LLMHandler.DeleteHistory)
			r.Get("/stats", deps.LLMHandler.GetStats)

			if deps.CVHandler != nil {
				r.Post("/cvs", deps.CVHandler.Create)
				r.Get("/cvs", deps.CVHandler.List)
				r.Get("/cvs/{id}", deps.CVHandler.GetByID)
				r.Delete("/cvs/{id}", deps.CVHandler.Delete)
				r.Post("/cvs/{id}/parse", deps.CVHandler.Parse)
			}

			if deps.JDHandler != nil {
				r.Post("/jds", deps.JDHandler.Create)
				r.Get("/jds", deps.JDHandler.List)
				r.Get("/jds/{id}", deps.JDHandler.GetByID)
				r.Put("/jds/{id}", deps.JDHandler.Update)
				r.Delete("/jds/{id}", deps.JDHandler.Delete)
				r.Post("/jds/{id}/analyze", deps.JDHandler.Analyze)
			}

			if deps.MatchingHandler != nil {
				r.Post("/matches", deps.MatchingHandler.RunMatch)
				r.Get("/matches", deps.MatchingHandler.List)
				r.Get("/matches/{id}", deps.MatchingHandler.GetByID)
				r.Get("/dashboard/stats", deps.MatchingHandler.GetDashboardStats)
			}

			if deps.KnowledgeHandler != nil {
				r.Post("/knowledge", deps.KnowledgeHandler.Create)
				r.Get("/knowledge", deps.KnowledgeHandler.List)
				r.Get("/knowledge/search", deps.KnowledgeHandler.Search)
				r.Get("/knowledge/categories", deps.KnowledgeHandler.ListCategories)
				r.Get("/knowledge/{id}", deps.KnowledgeHandler.GetByID)
				r.Delete("/knowledge/{id}", deps.KnowledgeHandler.Delete)
			}

			if deps.AdminHandler != nil {
				r.Get("/admin/adapters", deps.AdminHandler.ListAdapters)
				r.Post("/admin/adapters", deps.AdminHandler.CreateAdapter)
				r.Delete("/admin/adapters/{id}", deps.AdminHandler.DeleteAdapter)

				r.Get("/admin/prompts", deps.AdminHandler.ListPrompts)
				r.Post("/admin/prompts", deps.AdminHandler.CreatePrompt)
				r.Put("/admin/prompts/{id}", deps.AdminHandler.UpdatePrompt)
				r.Post("/admin/prompts/{id}/activate", deps.AdminHandler.ActivatePrompt)
				r.Delete("/admin/prompts/{id}", deps.AdminHandler.DeletePrompt)

				r.Get("/admin/settings", deps.AdminHandler.GetSettings)
				r.Put("/admin/settings", deps.AdminHandler.SaveSettings)

				r.Get("/admin/logs", deps.AdminHandler.ListLogs)
			}
		})
	})

	return r
}
