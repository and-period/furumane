package response

import (
	"time"

	"github.com/and-period/furumane/internal/auth/entity"
)

type Admin struct {
	ID           string              `json:"id"`           // 管理者ID
	ProviderType entity.ProviderType `json:"providerType"` // 認証種別
	Email        string              `json:"email"`        // メールアドレス
	PhoneNumber  string              `json:"phoneNumber"`  // 電話番号
	CreatedAt    time.Time           `json:"createdAt"`    // 登録日時
	UpdatedAt    time.Time           `json:"updatedAt"`    // 更新日時
}

type SignUpAdminResponse struct {
	AdminID string `json:"adminId"` // 管理者ID
}

type SignUpAdminWithOAuthResponse struct {
	Admin *Admin `json:"admin"` // 管理者情報
}

type GetAdminResponse struct {
	Admin *Admin `json:"admin"` // 管理者情報
}
