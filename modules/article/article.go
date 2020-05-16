package article

import (
	"base-site-api/internal/config"
	"base-site-api/middleware/auth"
	"github.com/gofiber/fiber"
)

func New(config *config.Config, api *fiber.Group) {
	handler := NewHandler(NewService(NewRepository(config.Database)))

	articles := api.Group("/v1/articles")
	articles.Use(auth.New(&auth.Config{
		SigningKey: config.SigningKey,
		Filter: auth.FilterOutGet,
	}))

	articles.Get("/", handler.List)
	articles.Post("/", handler.Create)
	articles.Put("/:id", handler.Update)
	articles.Delete("/:id", handler.Remove)
	articles.Get("/:id", handler.GetDetail)
}
