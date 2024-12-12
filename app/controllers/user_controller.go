package controllers

import (
	"time"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repositories"
	"github.com/kooroshh/fiber-boostrap/pkg/jwt"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Register", "controller")
	defer span.End()

	user := new(models.User)

	err := ctx.BodyParser(user)
	if err != nil {
		log.Println("Failed to parse request: ", err)
		return response.ResponseError(ctx, fiber.StatusBadRequest, err.Error())
	}

	err = user.Validate()
	if err != nil {
		log.Println("Failed to validate request: ", err)
		return response.ResponseError(ctx, fiber.StatusBadRequest, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Filed to encrypt the password: ", err)
		return response.ResponseError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	user.Password = string(hash)

	err = repositories.InsertNewUser(spanCtx, user)
	if err != nil {
		log.Println("Failed to insert new user: ", err)
		return response.ResponseError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	res := user
	res.Password = ""
	return response.ResponseSuccess(ctx, fiber.StatusOK, res)
}

func Login(ctx *fiber.Ctx) error {
	user := new(models.LoginRequest)

	err := ctx.BodyParser(user)
	if err != nil {
		log.Println("Failed to parse request: ", err)
		return response.ResponseError(ctx, fiber.StatusBadRequest, err.Error())
	}

	err = user.Validate()
	if err != nil {
		log.Println("Failed to validate request: ", err)
		return response.ResponseError(ctx, fiber.StatusBadRequest, err.Error())
	}

	res, err := repositories.GetUserByUsername(ctx.Context(), user.Username)
	if err != nil {
		log.Println("Failed to get username: ", err)
		return response.ResponseError(ctx, fiber.StatusBadRequest, "username or password invalid")
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(user.Password))
	if err != nil {
		log.Println("Failed to check password: ", err)
		return response.ResponseError(ctx, fiber.StatusBadRequest, "username or password invalid")
	}

	token, err := jwt.GenerateToken(res.Username, res.Fullname, "token")
	if err != nil {
		log.Println("Failed to check password: ", err)
		return response.ResponseError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	refresh_token, err := jwt.GenerateToken(res.Username, res.Fullname, "refresh_token")
	if err != nil {
		log.Println("Failed to check password: ", err)
		return response.ResponseError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	now := time.Now()
	userSession := &models.UserSession{
		UserID:              int(res.ID),
		Token:               token,
		RefreshToken:        refresh_token,
		TokenExpired:        now.Add(jwt.MapTokenType["token"]),
		RefreshTokenExpired: now.Add(jwt.MapTokenType["refresh_token"]),
	}

	err = repositories.CreateUserSession(ctx.Context(), userSession)
	if err != nil {
		log.Println("Failed to create token: ", err)
		return response.ResponseError(ctx, fiber.StatusBadRequest, "invalid token")
	}

	return response.ResponseSuccess(ctx, fiber.StatusOK, fiber.Map{
		"username":      res.Username,
		"full_name":     res.Fullname,
		"token":         token,
		"refresh_token": refresh_token,
	})
}

func Logout(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	err := repositories.DeleteUserSession(ctx.Context(), token)
	if err != nil {
		log.Println("failed delete user session:", err)
		return response.ResponseError(ctx, fiber.StatusInternalServerError, "internal server error")
	}

	return response.ResponseSuccess(ctx, fiber.StatusOK, "bye bye")
}

func RefreshToken(ctx *fiber.Ctx) error {

	now := time.Now()
	refreshToken := ctx.Get("Authorization")
	username := ctx.Locals("username").(string)
	fullname := ctx.Locals("full_name").(string)

	token, err := jwt.GenerateToken(username, fullname, "token")
	if err != nil {
		log.Println("Failed to generate token: ", err)
		return response.ResponseError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	err = repositories.UpdateRefreshToken(ctx.Context(), token, refreshToken, now.Add(jwt.MapTokenType["token"]))
	if err != nil {
		log.Println("Failed to update token: ", err)
		return response.ResponseError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.ResponseSuccess(ctx, fiber.StatusOK, fiber.Map{
		"token": token,
	})
}
