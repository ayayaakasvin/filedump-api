package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
    UserID    int    `json:"user_id"`
    SessionID string `json:"session_id"`
    jwt.RegisteredClaims
}
