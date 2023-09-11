package entity

import (
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/proto/auth"
)

// AdminAuth - 管理者認証情報
type AdminAuth struct {
	AdminID      string // 管理者ID
	AccessToken  string // アクセストークン
	RefreshToken string // 更新トークン
	ExpiresIn    int32  // 有効期限
}

func NewAdminAuth(adminID string, rs *cognito.AuthResult) *AdminAuth {
	return &AdminAuth{
		AdminID:      adminID,
		AccessToken:  rs.AccessToken,
		RefreshToken: rs.RefreshToken,
		ExpiresIn:    rs.ExpiresIn,
	}
}

func (a *AdminAuth) Proto() *auth.AdminAuth {
	return &auth.AdminAuth{
		AdminId:      a.AdminID,
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		ExpiresIn:    a.ExpiresIn,
	}
}
