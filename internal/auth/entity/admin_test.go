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
		PhoneNumber:  "09012341234",
	}
	actual := NewAdmin(params)

	t.Run("constructor", func(t *testing.T) {
		expect := &Admin{
			ID:           "admin-id",
			CognitoID:    "cognito-id",
			ProviderType: ProviderTypeEmail,
			Email:        "test@example.com",
			PhoneNumber:  "09012341234",
		}
		assert.Equal(t, expect, actual)
	})
	t.Run("international phone number", func(t *testing.T) {
		assert.Equal(t, "+819012341234", actual.InternationalPhoneNumber())
		actual := &Admin{}
		assert.Empty(t, actual.InternationalPhoneNumber())
	})
}
