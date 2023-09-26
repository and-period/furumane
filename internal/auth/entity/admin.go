package entity

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID           string         `gorm:"primaryKey;<-:create"` // 管理者ID
	CognitoID    string         `gorm:""`                     // 管理者ID（Cognito用）
	ProviderType ProviderType   `gorm:""`                     // 認証種別
	Email        string         `gorm:"default:null"`         // メールアドレス
	PhoneNumber  string         `gorm:"default:null"`         // 電話番号
	CreatedAt    time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time      `gorm:""`                     // 更新日時
	VerifiedAt   time.Time      `gorm:"default:null"`         // 確認日時
	DeletedAt    gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type AdminParams struct {
	AdminID      string
	CognitID     string
	ProviderType ProviderType
	Email        string
	PhoneNumber  string
}

func NewAdmin(params *AdminParams) *Admin {
	return &Admin{
		ID:           params.AdminID,
		CognitoID:    params.CognitID,
		ProviderType: params.ProviderType,
		Email:        params.Email,
		PhoneNumber:  params.PhoneNumber,
	}
}

func (a *Admin) InternationalPhoneNumber() string {
	if a == nil || a.PhoneNumber == "" {
		return ""
	}
	return strings.Replace(a.PhoneNumber, "0", "+81", 1)
}
