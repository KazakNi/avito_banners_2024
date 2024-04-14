package dto

import (
	"banners/internal/auth/domain/entity"

	"golang.org/x/crypto/bcrypt"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreatedUserId struct {
	User_id string `json:"user_id"`
}

type Token struct {
	BearerToken string `json:"token"`
}

type User entity.User

func (user *User) CheckPassword(providedPassword string, db_password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(db_password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
