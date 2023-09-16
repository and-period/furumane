package request

type SignInAdminRequest struct {
	Key      string `json:"key" validate:"required"`      // キー
	Password string `json:"password" validate:"required"` // パスワード
}

type RefreshAdminTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"` // リフレッシュトークン
}
