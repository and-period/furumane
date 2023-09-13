package entity

type ProviderType int32 // 認証種別

const (
	ProviderTypeUnknown ProviderType = 0
	ProviderTypeEmail   ProviderType = 1 // メールアドレス認証
	ProviderTypeOAuth   ProviderType = 2 // OAuth認証
)
