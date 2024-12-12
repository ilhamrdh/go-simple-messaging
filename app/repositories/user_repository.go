package repositories

import (
	"context"
	"time"

	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/pkg/database"
	"go.elastic.co/apm"
)

func InsertNewUser(ctx context.Context, user *models.User) error {
	span, _ := apm.StartSpan(ctx, "InsertNewUser", "repository")
	defer span.End()
	return database.DB.Create(user).Error
}

func GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	user := models.User{}

	err := database.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func CreateUserSession(ctx context.Context, session *models.UserSession) error {
	return database.DB.Create(session).Error
}

func GetUserSessionByToken(ctx context.Context, token string) (models.UserSession, error) {
	userSession := models.UserSession{}
	err := database.DB.Where("token = ?", token).First(&userSession).Error
	return userSession, err
}

func DeleteUserSession(ctx context.Context, token string) error {
	return database.DB.Exec("DELETE FROM user_sessions WHERE token = ?", token).Error
}

func UpdateRefreshToken(ctx context.Context, token, refreshToken string, tokenExpired time.Time) error {
	return database.DB.Exec("UPDATE user_sessions SET token = ?, token_expired = ? WHERE refresh_token = ?", token, tokenExpired, refreshToken).Error
}
