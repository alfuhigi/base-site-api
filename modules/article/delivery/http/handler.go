package http

import (
	"base-site-api/modules/article/service"
	"base-site-api/responses"

	"github.com/gofiber/fiber"
)

type ArticleHandler struct {
	service service.Service
}

func New(s service.Service) *ArticleHandler {
	return &ArticleHandler{
		service: s,
	}
}

func (h *ArticleHandler) List(c *fiber.Ctx) {
	articles, err := h.service.FindAll(c.Params("sort"))

	if err != nil {
		c.Status(500).Send(responses.ErrorResponse{
			Message: "Problem with getting the articles",
			Error:   err.Error(),
		})
		return
	}

	c.Send(&articles)
}
