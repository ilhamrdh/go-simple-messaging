package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/repositories"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
)

func GetHistory(ctx *fiber.Ctx) error {
	res, err := repositories.GetAllMessage(ctx.Context())
	if err != nil {
		log.Println("error repository: ", err)
		return response.ResponseError(ctx, fiber.StatusInternalServerError, "internal server error")
	}
	return response.ResponseSuccess(ctx, fiber.StatusOK, res)
}
