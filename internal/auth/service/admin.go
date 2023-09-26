package service

import (
	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/internal/auth/response"
)

type Admin struct {
	response.Admin
}

func NewAdmin(admin *entity.Admin) *Admin {
	return &Admin{
		Admin: response.Admin{
			ID:           admin.ID,
			ProviderType: admin.ProviderType,
			Email:        admin.Email,
			PhoneNumber:  admin.PhoneNumber,
			CreatedAt:    admin.CreatedAt,
			UpdatedAt:    admin.UpdatedAt,
		},
	}
}

func (a *Admin) Response() *response.Admin {
	return &a.Admin
}
