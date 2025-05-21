package usecase

import (
	"time"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
	"github.com/golang-jwt/jwt"
)

type JwtManager struct {
	secret string
}

func NewJwtManager(secret string) *JwtManager {
	return &JwtManager{
		secret: secret,
	}
}

func (jw *JwtManager) NewToken(user model.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(jw.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
