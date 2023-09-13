package service

import (
	"context"

	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/proto/auth"
)

func (s *service) SignInAdmin(ctx context.Context, req *auth.SignInAdminRequest) (*auth.SignInAdminResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	rs, err := s.adminAuth.SignIn(ctx, req.Key, req.Password)
	if err != nil {
		return nil, gRPCError(err)
	}
	admin, err := s.getAdminAuth(ctx, rs)
	if err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.SignInAdminResponse{
		Auth: admin.Proto(),
	}
	return res, nil
}

func (s *service) SignInAdminWithOAuth(
	ctx context.Context, req *auth.SignInAdminWithOAuthRequest,
) (*auth.SignInAdminWithOAuthResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	rs := &cognito.AuthResult{AccessToken: req.AccessToken}
	admin, err := s.getAdminAuth(ctx, rs)
	if err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.SignInAdminWithOAuthResponse{
		Auth: admin.Proto(),
	}
	return res, nil
}

func (s *service) SignOutAdmin(ctx context.Context, req *auth.SignOutAdminRequest) (*auth.SignOutAdminResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	if err := s.adminAuth.SignOut(ctx, req.AccessToken); err != nil {
		return nil, gRPCError(err)
	}
	return &auth.SignOutAdminResponse{}, nil
}

func (s *service) RefreshAdminToken(ctx context.Context, req *auth.RefreshAdminTokenRequest) (*auth.RefreshAdminTokenResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	rs, err := s.adminAuth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, gRPCError(err)
	}
	admin, err := s.getAdminAuth(ctx, rs)
	if err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.RefreshAdminTokenResponse{
		Auth: admin.Proto(),
	}
	return res, nil
}

func (s *service) getAdminAuth(ctx context.Context, rs *cognito.AuthResult) (*entity.AdminAuth, error) {
	username, err := s.adminAuth.GetUsername(ctx, rs.AccessToken)
	if err != nil {
		return nil, err
	}
	admin, err := s.db.Admin.GetByCognitoID(ctx, username)
	if err != nil {
		return nil, err
	}
	return entity.NewAdminAuth(admin, rs), nil
}
