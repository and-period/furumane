package service

import (
	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/internal/auth/response"
)

type AdminAuth struct {
	response.AdminAuth
}

func NewAdminAuth(auth *entity.AdminAuth) *AdminAuth {
	return &AdminAuth{
		AdminAuth: response.AdminAuth{
			AdminID:      auth.AdminID,
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
			ExpiresIn:    auth.ExpiresIn,
		},
	}
}

func (a *AdminAuth) Response() *response.AdminAuth {
	return &a.AdminAuth
}
