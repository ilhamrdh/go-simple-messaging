package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseSuccess(ctx *fiber.Ctx, httpCode int, data interface{}) error {
	return ctx.Status(httpCode).JSON(Response{
		Message: "success",
		Data:    data,
	})
}

func ResponseError(ctx *fiber.Ctx, httpCode int, message string) error {
	return ctx.Status(httpCode).JSON(Response{
		Message: message,
	})
}
