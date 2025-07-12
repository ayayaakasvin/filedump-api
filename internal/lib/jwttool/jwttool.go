package jwttool

import (
	"fmt"
	"log"
	"os"
	"time"
	"up-down-server/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var _JWTSecret []byte

func init() {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		log.Fatalf("NO JWT SECRET KEY PROVIDED")
	}

	_JWTSecret = []byte(key)
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return _JWTSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %s", err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GenerateAccessToken(userId int, sessionId string, ttl time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.CustomClaims{
		UserID:    userId,
		SessionID: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(_JWTSecret)
	if err != nil {
		log.Print(err)
		return ""
	}

	return tokenString
}

func GenerateRefreshToken(userId int, ttl time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.CustomClaims{
		UserID:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(_JWTSecret)
	if err != nil {
		log.Print(err)
		return ""
	}

	return tokenString
}

func FetchUserID(userIdAny any) (int, error) {
	switch v := userIdAny.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return -1, fmt.Errorf("invalid user id type")
	}
}
