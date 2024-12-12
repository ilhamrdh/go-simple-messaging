package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        uint   `gorm:"primarykey"`
	Fullname  string `json:"full_name" gorm:"type:varchar(100);" validate:"required,min=6"`
	Username  string `json:"username" gorm:"unique;" validate:"required,min=6,max=32"`
	Password  string `json:"password" gorm:"type:varchar(255);" validate:"required,min=6"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (l User) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type UserSession struct {
	ID                  uint      `gorm:"primarykey"`
	UserID              int       `json:"user_id" gorm:"type:int;" validate:"required"`
	Token               string    `json:"token" gorm:"type:varchar(255);" validate:"required"`
	RefreshToken        string    `json:"refresh_token" gorm:"type:varchar(255);" validate:"required"`
	TokenExpired        time.Time `json:"-" validate:"required"`
	RefreshTokenExpired time.Time `json:"-" validate:"required"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (l UserSession) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
