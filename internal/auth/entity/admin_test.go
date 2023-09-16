package entity

import (
	"testing"

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
}
