package server

import (
	"JacFARM/internal/config"
	"JacFARM/internal/http/handlers"
	"JacFARM/internal/http/middlewares"

	"github.com/bytedance/sonic"
	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

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

	serviceGroup := apiV1.Group("/service")
	serviceGroup.Post("/flags", middlewares.ServiceAuthMiddleware(apiKey), h.PutFlag())

	return r
}
