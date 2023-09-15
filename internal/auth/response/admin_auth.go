package response

// AdminAuth 管理者認証情報
type AdminAuth struct {
	AdminID      string `json:"adminId"`      // 管理者ID
	AccessToken  string `json:"accessToken"`  // アクセストークン
	RefreshToken string `json:"refreshToken"` // リフレッシュトークン
	ExpiresIn    int32  `json:"expiresIn"`    // 有効期限(sec)
}

type SignInAdminResponse struct {
	AdminAuth *AdminAuth `json:"auth"` // 管理者認証情報
}

type SignInAdminWithOAuthResponse struct {
	AdminAuth *AdminAuth `json:"auth"` // 管理者認証情報
}

type GetAdminAuthResponse struct {
	AdminAuth *AdminAuth `json:"auth"` // 管理者認証情報
}

type RefreshAdminTokenResponse struct {
	AdminAuth *AdminAuth `json:"auth"` // 管理者認証情報
}
