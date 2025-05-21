package testing

import (
	"testing"
	"time"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
	"github.com/Temutjin2k/Tyndau/auth_service/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestJwtManager_ValidateToken(t *testing.T) {
	secret := "your-32-byte-secret-key-here-1234567890"
	manager := usecase.NewJwtManager(secret)

	// Test valid token
	t.Run("valid token", func(t *testing.T) {
		user := model.User{ID: 123, Email: "test@example.com"}
		token, err := manager.NewToken(user, time.Hour)
		require.NoError(t, err)

		valid, err := manager.ValidateToken(token)
		require.NoError(t, err)
		require.True(t, valid)
	})

	// Test expired token
	t.Run("expired token", func(t *testing.T) {
		user := model.User{ID: 123, Email: "test@example.com"}
		token, err := manager.NewToken(user, -time.Hour) // Negative duration for expired token
		require.NoError(t, err)

		valid, err := manager.ValidateToken(token)
		require.Error(t, err)
		require.False(t, valid)
		require.Contains(t, err.Error(), "expired")
	})

	// Test invalid signature
	t.Run("invalid signature", func(t *testing.T) {
		otherManager := usecase.NewJwtManager("different-secret-key-1234567890")
		user := model.User{ID: 123, Email: "test@example.com"}
		token, err := otherManager.NewToken(user, time.Hour)
		require.NoError(t, err)

		valid, err := manager.ValidateToken(token) // Validate with different secret
		require.Error(t, err)
		require.False(t, valid)
	})
}
