package api

import (
	"context"
	"net/http"

	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/internal/auth/request"
	"github.com/and-period/furumane/internal/auth/response"
	"github.com/and-period/furumane/internal/auth/service"
	"github.com/and-period/furumane/internal/util"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/gin-gonic/gin"
)

func (c *controller) adminAuthRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/auth")
	g.POST("", c.SignInAdmin)
	g.DELETE("", c.SignOutAdmin)
	g.GET("", c.GetAdminAuth)
	g.POST("/refresh", c.RefreshAdminToken)
}

// SignInAdmin 管理者サインイン（メールアドレス認証）
func (c *controller) SignInAdmin(ctx *gin.Context) {
	req := &request.SignInAdminRequest{}
	if err := c.bind(ctx, req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	rs, err := c.adminAuth.SignIn(ctx, req.Key, req.Password)
	if err != nil {
		httpError(ctx, err)
		return
	}
	admin, err := c.getAdminAuth(ctx, rs)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.SignInAdminResponse{
		AdminAuth: service.NewAdminAuth(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// SignOutAdmin 管理者サインアウト
func (c *controller) SignOutAdmin(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err.Error())
		return
	}
	if err := c.adminAuth.SignOut(ctx, token); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// GetAdminAuth 管理者認証情報取得
func (c *controller) GetAdminAuth(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err.Error())
		return
	}
	rs := &cognito.AuthResult{AccessToken: token}
	admin, err := c.getAdminAuth(ctx, rs)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.GetAdminAuthResponse{
		AdminAuth: service.NewAdminAuth(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// RefreshAdminToken 管理者アクセストークンの更新
func (c *controller) RefreshAdminToken(ctx *gin.Context) {
	req := &request.RefreshAdminTokenRequest{}
	if err := c.bind(ctx, req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	rs, err := c.adminAuth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		httpError(ctx, err)
		return
	}
	admin, err := c.getAdminAuth(ctx, rs)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.RefreshAdminTokenResponse{
		AdminAuth: service.NewAdminAuth(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *controller) getAdminAuth(ctx context.Context, rs *cognito.AuthResult) (*entity.AdminAuth, error) {
	username, err := c.adminAuth.GetUsername(ctx, rs.AccessToken)
	if err != nil {
		return nil, err
	}
	admin, err := c.db.Admin.GetByCognitoID(ctx, username)
	if err != nil {
		return nil, err
	}
	return entity.NewAdminAuth(admin, rs), nil
}
