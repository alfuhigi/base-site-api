package page

import (
	"base-site-api/config"
	"base-site-api/middleware/auth"
	"github.com/gofiber/fiber"
)

// TODO: prepare whole new module
func New(config *config.Config, api *fiber.Group) {
	handler := NewHandler(NewService(NewRepository(config.Database)))

	pages := api.Group("/v1/pages")
	pages.Use(auth.New(&auth.Config{
		SigningKey: config.SigningKey,
		Filter:     auth.FilterGetOnly,
	}))

	pages.Get("/", handler.ListCategories)
	pages.Get("/:page-category", handler.ListPages)
	pages.Get("/:slug", handler.GetDetail)
	pages.Post("/:page-category", handler.Create)
	pages.Put("/:id", handler.Update)
	pages.Delete("/:id", handler.Remove)

}
