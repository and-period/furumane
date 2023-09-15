package request

type SignUpAdminRequest struct {
	Email                string `json:"email" validate:"required,max=256,email"`                   // メールアドレス
	Password             string `json:"password" validate:"min=8,max=32,password"`                 // パスワード
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"` // パスワード（確認用）
}

type VerifyAdminRequest struct {
	AdminID    string `json:"adminId" validate:"required"`    // 管理者ID
	VerifyCode string `json:"verifyCode" validate:"required"` // 検証コード
}

type UpdateAdminEmailRequest struct {
	Email string `json:"email" validate:"required,max=256,email"` // メールアドレス
}

type VerifyAdminEmailRequest struct {
	VerifyCode string `json:"verifyCode" validate:"required"` // 検証コード
}

type UpdateAdminPasswordRequest struct {
	OldPassword          string `json:"oldPassword"`                      // 現在のパスワード
	NewPassword          string `validate:"min=8,max=32,password"`        // 新しいパスワード
	PasswordConfirmation string `validate:"required,eqfield=NewPassword"` // パスワード（確認用）
}

type ForgotAdminPasswordRequest struct {
	Email string `json:"email" validate:"required"` // メールアドレス
}

type ResetAdminPasswordRequest struct {
	Email                string `json:"email" validate:"required"`      // メールアドレス
	VerifyCode           string `json:"verifyCode" validate:"required"` // 検証コード
	Password             string `validate:"min=8,max=32,password"`      // 新しいパスワード
	PasswordConfirmation string `validate:"required,eqfield=Password"`  // パスワード（確認用）
}
