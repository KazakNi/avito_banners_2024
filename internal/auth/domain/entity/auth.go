package entity

import "github.com/golang-jwt/jwt/v4"

type User struct {
	Id          int    `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	AdminRights bool   `json:"is_admin" db:"is_admin"`
}

type CustomClaims struct {
	AdminRights bool
	jwt.RegisteredClaims
}
