package modules

import (
	"github.com/gofiber/fiber"
	"math"
	"strconv"
)

// Handler wrap common logic  for handlers
type Handler struct {
}

// CalculatePagination return Pagination struct with calculated TotalPages
func (h *Handler) CalculatePagination(page int, size int, count int) *Pagination {
	return &Pagination{
		Page:       page,
		PageSize:   size,
		Total:      count,
		TotalPages: math.Ceil(float64(count) / float64(size)),
	}
}

// JSON handle json response and error handling for handler
func (h *Handler) JSON(c *fiber.Ctx, status int, data interface{}) {
	if err := c.Status(status).JSON(data); err != nil {
		c.Next(err)
	}
}

// Error Handle Error response
func (h *Handler) Error(c *fiber.Ctx, status int) {
	c.Next(fiber.NewError(status))
}

// ErrorWithMessage Handle Error response with custom message
func (h *Handler) ErrorWithMessage(c *fiber.Ctx, status int, message string) {
	c.Next(fiber.NewError(status, message))
}

// ParseUserID parse user id from context and convert it to uint
func (h *Handler) ParseUserID(c *fiber.Ctx) uint {
	return c.Locals("userID").(uint)
}

// ParseID parse id from url and convert to uint
func (h *Handler) ParseID(c *fiber.Ctx) (uint, error) {
	id := c.Params("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return uint(0), err
	}

	return uint(uid), err
}
