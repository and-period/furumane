package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumane/internal/auth/database"
	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/internal/auth/request"
	"github.com/and-period/furumane/internal/auth/response"
	"github.com/and-period/furumane/internal/auth/service"
	"github.com/and-period/furumane/internal/util"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/uuid"
	"github.com/gin-gonic/gin"
)

func (c *controller) adminRoutes(rg *gin.RouterGroup) {
	g := rg.Group("")
	g.POST("", c.SignUpAdmin)
	g.POST("/verified", c.VerifyAdmin)
	g.POST("/oauth", c.SignUpAdminWithOAuth)
	g.PUT("/email", c.UpdateAdminEmail)
	g.POST("/email/verified", c.VerifyAdminEmail)
	g.PUT("/password", c.UpdateAdminPassword)
	g.POST("/password/forgot", c.ForgotAdminPassword)
	g.PUT("/password/reset", c.ResetAdminPassword)
	g.GET("/:adminId", c.GetAdmin)
	g.DELETE("/:adminId", c.DeleteAdmin)
}

// SignUpAdmin 管理者登録（メールアドレス認証）
func (c *controller) SignUpAdmin(ctx *gin.Context) {
	req := &request.SignUpAdminRequest{}
	if err := c.bind(ctx, req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	cognitoID := uuid.Base58Encode(c.uuid())
	params := &entity.AdminParams{
		AdminID:      uuid.Base58Encode(c.uuid()),
		CognitID:     cognitoID,
		ProviderType: entity.ProviderTypeEmail,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
	}
	admin := entity.NewAdmin(params)
	fn := func(ctx context.Context) error {
		params := &cognito.SignUpParams{
			Username:    cognitoID,
			Email:       admin.Email,
			PhoneNumber: admin.InternationalPhoneNumber(),
			Password:    req.Password,
		}
		return c.adminAuth.SignUp(ctx, params)
	}
	err := c.db.Admin.Create(ctx, admin, fn)
	if err != nil && !errors.Is(err, database.ErrAlreadyExists) {
		httpError(ctx, err)
		return
	}
	res := &response.SignUpAdminResponse{
		AdminID: admin.ID,
	}
	ctx.JSON(http.StatusOK, res)
}

// VerifyAdmin 管理者登録後の確認 (メールアドレス認証)
func (c *controller) VerifyAdmin(ctx *gin.Context) {
	req := &request.VerifyAdminRequest{}
	if err := c.bind(ctx, req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	admin, err := c.db.Admin.Get(ctx, req.AdminID, "cognito_id", "verified_at")
	if err != nil {
		httpError(ctx, err)
		return
	}
	if !admin.VerifiedAt.IsZero() {
		preconditionFailed(ctx, "this admin is already verified")
		return
	}
	if err := c.adminAuth.ConfirmSignUp(ctx, admin.CognitoID, req.VerifyCode); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// SignUpAdminWithOAuth 管理者登録 (OAuth認証)
func (c *controller) SignUpAdminWithOAuth(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err.Error())
		return
	}
	au, err := c.adminAuth.GetUser(ctx, token)
	if err != nil {
		httpError(ctx, err)
		return
	}
	params := &entity.AdminParams{
		AdminID:      uuid.Base58Encode(c.uuid()),
		CognitID:     au.Username,
		ProviderType: entity.ProviderTypeOAuth,
		Email:        au.Email,
	}
	admin := entity.NewAdmin(params)
	fn := func(ctx context.Context) error {
		return nil // Cognitoへはすでに登録済みのため何もしない
	}
	err = c.db.Admin.Create(ctx, admin, fn)
	if err != nil && !errors.Is(err, database.ErrAlreadyExists) {
		httpError(ctx, err)
		return
	}
	if err := c.db.Admin.UpdateVerifiedAt(ctx, admin.ID); err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.SignUpAdminWithOAuthResponse{
		Admin: service.NewAdmin(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// GetAdmin 管理者情報取得
func (c *controller) GetAdmin(ctx *gin.Context) {
	adminID := util.GetParam(ctx, "adminId")
	admin, err := c.db.Admin.Get(ctx, adminID)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.GetAdminResponse{
		Admin: service.NewAdmin(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// UpdateAdminEmail 管理者メールアドレス更新
func (c *controller) UpdateAdminEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err.Error())
		return
	}
	req := &request.UpdateAdminEmailRequest{}
	if err := c.bind(ctx, req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	username, err := c.adminAuth.GetUsername(ctx, token)
	if err != nil {
		httpError(ctx, err)
		return
	}
	admin, err := c.db.Admin.GetByCognitoID(ctx, username)
	if err != nil {
		httpError(ctx, err)
		return
	}
	if admin.ProviderType != entity.ProviderTypeEmail {
		preconditionFailed(ctx, "not allow provider type to change email")
		return
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: token,
		Username:    username,
		OldEmail:    admin.Email,
		NewEmail:    req.Email,
	}
	if err := c.adminAuth.ChangeEmail(ctx, params); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// VerifyAdminEmail 管理者メールアドレス更新後の確認
func (c *controller) VerifyAdminEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err.Error())
		return
	}
	req := &request.VerifyAdminEmailRequest{}
	if err := c.bind(ctx, req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	username, err := c.adminAuth.GetUsername(ctx, token)
	if err != nil {
		httpError(ctx, err)
		return
	}
	admin, err := c.db.Admin.GetByCognitoID(ctx, username)
	if err != nil {
		httpError(ctx, err)
		return
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: token,
		Username:    username,
		VerifyCode:  req.VerifyCode,
	}
	email, err := c.adminAuth.ConfirmChangeEmail(ctx, params)
	if err != nil {
		httpError(ctx, err)
		return
	}
	if err := c.db.Admin.UpdateEmail(ctx, admin.ID, email); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// UpdateAdminPassword 管理者パスワード更新
func (c *controller) UpdateAdminPassword(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err.Error())
		return
	}
	req := &request.UpdateAdminPasswordRequest{}
	if err := c.bind(ctx, req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	params := &cognito.ChangePasswordParams{
		AccessToken: token,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	if err := c.adminAuth.ChangePassword(ctx, params); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// ForgotAdminPassword 管理者パスワードリセット (メール/SMS送信)
func (c *controller) ForgotAdminPassword(ctx *gin.Context) {
	req := &request.ForgotAdminPasswordRequest{}
	if err := c.bind(ctx, req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	admin, err := c.db.Admin.GetByEmail(ctx, req.Email, "cognito_id")
	if err != nil {
		httpError(ctx, err)
		return
	}
	if err := c.adminAuth.ForgotPassword(ctx, admin.CognitoID); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// ResetAdminPassword 管理者パスワードリセット (パスワード更新)
func (c *controller) ResetAdminPassword(ctx *gin.Context) {
	req := &request.ResetAdminPasswordRequest{}
	if err := c.bind(ctx, req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	admin, err := c.db.Admin.GetByEmail(ctx, req.Email, "cognito_id")
	if err != nil {
		httpError(ctx, err)
		return
	}
	params := &cognito.ConfirmForgotPasswordParams{
		Username:    admin.CognitoID,
		VerifyCode:  req.VerifyCode,
		NewPassword: req.Password,
	}
	if err := c.adminAuth.ConfirmForgotPassword(ctx, params); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// DeleteAdmin 管理者退会
func (c *controller) DeleteAdmin(ctx *gin.Context) {
	adminID := util.GetParam(ctx, "adminId")
	admin, err := c.db.Admin.Get(ctx, adminID)
	if errors.Is(err, database.ErrNotFound) {
		ctx.Status(http.StatusNoContent) // 退会済み
		return
	}
	if err != nil {
		httpError(ctx, err)
		return
	}
	fn := func(ctx context.Context) error {
		return c.adminAuth.DeleteUser(ctx, admin.CognitoID)
	}
	if err := c.db.Admin.Delete(ctx, admin.ID, fn); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
