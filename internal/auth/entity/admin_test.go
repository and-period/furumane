package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumane/proto/auth"
	"github.com/stretchr/testify/assert"
)

func TestAdmin(t *testing.T) {
	t.Parallel()
	params := &AdminParams{
		AdminID:      "admin-id",
		CognitID:     "cognito-id",
		ProviderType: ProviderTypeEmail,
		Email:        "test@example.com",
	}
	actual := NewAdmin(params)

	t.Run("constructor", func(t *testing.T) {
		expect := &Admin{
			ID:           "admin-id",
			CognitoID:    "cognito-id",
			ProviderType: ProviderTypeEmail,
			Email:        "test@example.com",
		}
		assert.Equal(t, expect, actual)
	})
	t.Run("proto", func(t *testing.T) {
		expect := &auth.Admin{
			Id:           "admin-id",
			ProviderType: auth.ProviderType_EMAIL,
			Email:        "test@example.com",
			CreatedAt:    time.Time{}.Unix(),
			UpdatedAt:    time.Time{}.Unix(),
		}
		assert.Equal(t, expect, actual.Proto())
	})
}
