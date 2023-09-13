package entity

import (
	"testing"

	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/proto/auth"
	"github.com/stretchr/testify/assert"
)

func TestAdminAuth(t *testing.T) {
	t.Parallel()
	a := &Admin{
		ID:           "admin-id",
		CognitoID:    "cognito-id",
		ProviderType: ProviderTypeEmail,
	}
	rs := &cognito.AuthResult{
		IDToken:      "id-token",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}
	actual := NewAdminAuth(a, rs)

	t.Run("constructor", func(t *testing.T) {
		expect := &AdminAuth{
			AdminID:      "admin-id",
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
		}
		assert.Equal(t, expect, actual)
	})
	t.Run("proto", func(t *testing.T) {
		expect := &auth.AdminAuth{
			AdminId:      "admin-id",
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
		}
		assert.Equal(t, expect, actual.Proto())
	})
}
