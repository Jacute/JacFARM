package server

import (
	"JacFARM/internal/config"
	"JacFARM/internal/http/handlers"
	"JacFARM/internal/http/middlewares"
	"net/http"
	"net/http/pprof"

	"github.com/bytedance/sonic"
	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

// @title           JacFARM API Docs

// @host      localhost:15050
// @BasePath  /jacfarm-api

// @securityDefinitions.basic  BasicAuth
func setupRouter(h *handlers.Handlers, cfg *config.HTTPConfig, apiKey string) *fiber.App {
	r := fiber.New(fiber.Config{
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		IdleTimeout:     cfg.IdleTimeout,
		BodyLimit:       20 * 1024 * 1024, // 20 MB
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		JSONEncoder:     sonic.Marshal,
		JSONDecoder:     sonic.Unmarshal,
		AppName:         "",
	})
	r.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	apiV1 := r.Group("/api/v1")
	apiV1.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowHeaders:     []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))
	apiV1.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	flagGroup := apiV1.Group("/flags")
	flagGroup.Get("/", h.ListFlags())
	flagGroup.Post("/", h.PutFlag())
	flagGroup.Get("/statuses", h.GetStatuses())
	flagGroup.Get("/count", h.GetFlagsCount())

	exploitGroup := apiV1.Group("/exploits")
	exploitGroup.Get("/", h.ListExploits())
	exploitGroup.Get("/short", h.ListShortExploits())
	exploitGroup.Post("/", h.UploadExploit())
	exploitGroup.Delete("/:id", h.DeleteExploit())
	exploitGroup.Post("/:id/toggle", h.ToggleExploit())

	teamGroup := apiV1.Group("/teams")
	teamGroup.Get("/", h.ListTeams())
	teamGroup.Post("/", h.AddTeam())
	teamGroup.Delete("/:id", h.DeleteTeam())
	teamGroup.Get("/short", h.ListShortTeams())

	configGroup := apiV1.Group("/config")
	configGroup.Get("/", h.GetConfig())
	configGroup.Patch("/:id", h.UpdateConfig())

	logGroup := apiV1.Group("/logs")
	logGroup.Get("/", h.ListLogs())
	logGroup.Get("/levels", h.ListLogLevels())
	logGroup.Get("/modules", h.ListModules())

	serviceGroup := apiV1.Group("/service")
	serviceGroup.Post("/flags", middlewares.ServiceAuthMiddleware(apiKey), h.PutFlag())
	serviceGroup.Get("/teams", middlewares.ServiceAuthMiddleware(apiKey), h.ListTeams())
	serviceGroup.Get("/config", middlewares.ServiceAuthMiddleware(apiKey), h.GetConfig())

	return r
}

func setupPprofRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	return mux
}
