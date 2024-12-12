package router

import (
	"time"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/repositories"
	"github.com/kooroshh/fiber-boostrap/pkg/jwt"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
)

func MiddlewareAuth(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		log.Println("authorization is empty")
		return response.ResponseError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	_, err := repositories.GetUserSessionByToken(ctx.Context(), auth)
	if err != nil {
		log.Println("failed to get user session on db: ", err)
		return response.ResponseError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	claim, err := jwt.ValidateToken(auth)
	if err != nil {
		log.Println(err)
		return response.ResponseError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token is expired: ", claim.ExpiresAt)
		return response.ResponseError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}

func MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		log.Println("authorization is empty")
		return response.ResponseError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	claim, err := jwt.ValidateToken(auth)
	if err != nil {
		log.Println(err)
		return response.ResponseError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token is expired: ", claim.ExpiresAt)
		return response.ResponseError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)
	return ctx.Next()
}
