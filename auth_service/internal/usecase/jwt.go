package usecase

import (
	"errors"
	"fmt"
	"math"
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

func (jw *JwtManager) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jw.secret), nil
	})

	if err != nil {
		return false, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return false, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return false, errors.New("token expired")
			}
		} else {
			return false, errors.New("invalid expiration claim")
		}

		// Check user ID (int64)
		if uid, ok := claims["uid"].(float64); !ok {
			return false, errors.New("missing user ID claim")
		} else {
			// Verify it's a whole number that fits in int64
			if uid != math.Trunc(uid) {
				return false, errors.New("user ID must be an integer")
			}
			if uid < float64(math.MinInt64) || uid > float64(math.MaxInt64) {
				return false, errors.New("user ID out of range")
			}
		}

		// Check email (string)
		if _, ok := claims["email"].(string); !ok {
			return false, errors.New("missing email claim")
		}
	}

	return true, nil
}
