package entity

import "github.com/and-period/furumane/pkg/cognito"

// AdminAuth - 管理者認証情報
type AdminAuth struct {
	AdminID      string // 管理者ID
	AccessToken  string // アクセストークン
	RefreshToken string // 更新トークン
	ExpiresIn    int32  // 有効期限
}

func NewAdminAuth(admin *Admin, rs *cognito.AuthResult) *AdminAuth {
	return &AdminAuth{
		AdminID:      admin.ID,
		AccessToken:  rs.AccessToken,
		RefreshToken: rs.RefreshToken,
		ExpiresIn:    rs.ExpiresIn,
	}
}
