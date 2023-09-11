package service

import (
	"context"
	"errors"

	"github.com/and-period/furumane/proto/auth"
)

func (s *service) SignUpAdmin(
	_ context.Context, req *auth.SignUpAdminRequest,
) (*auth.SignUpAdminResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	if req.Password != req.PasswordConfirmation {
		return nil, invalidArgument(errors.New("service: unmatch password and password confirmation"))
	}
	// TODO: 詳細の実装
	res := &auth.SignUpAdminResponse{}
	return res, nil
}

func (s *service) VerifyAdmin(_ context.Context, req *auth.VerifyAdminRequest) (*auth.VerifyAdminResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	// TODO: 詳細の実装
	res := &auth.VerifyAdminResponse{}
	return res, nil
}

func (s *service) SignUpAdminWithOAuth(
	_ context.Context, req *auth.SignUpAdminWithOAuthRequest,
) (*auth.SignUpAdminWithOAuthResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	// TODO: 詳細の実装
	res := &auth.SignUpAdminWithOAuthResponse{}
	return res, nil
}

func (s *service) UpdateAdminEmail(_ context.Context, req *auth.UpdateAdminEmailRequest) (*auth.UpdateAdminEmailResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	// TODO: 詳細の実装
	res := &auth.UpdateAdminEmailResponse{}
	return res, nil
}

func (s *service) VerifyAdminEmail(_ context.Context, req *auth.VerifyAdminEmailRequest) (*auth.VerifyAdminEmailResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	// TODO: 詳細の実装
	res := &auth.VerifyAdminEmailResponse{}
	return res, nil
}

func (s *service) UpdateAdminPassword(
	_ context.Context, req *auth.UpdateAdminPasswordRequest,
) (*auth.UpdateAdminPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	if req.NewPassword != req.PasswordConfirmation {
		return nil, invalidArgument(errors.New("service: unmatch password and password confirmation"))
	}
	// TODO: 詳細の実装
	res := &auth.UpdateAdminPasswordResponse{}
	return res, nil
}

func (s *service) ForgotAdminPassword(
	_ context.Context, req *auth.ForgotAdminPasswordRequest,
) (*auth.ForgotAdminPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	// TODO: 詳細の実装
	res := &auth.ForgotAdminPasswordResponse{}
	return res, nil
}

func (s *service) ResetAdminPassword(_ context.Context, req *auth.ResetAdminPasswordRequest) (*auth.ResetAdminPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	// TODO: 詳細の実装
	res := &auth.ResetAdminPasswordResponse{}
	return res, nil
}

func (s *service) DeleteAdmin(_ context.Context, req *auth.DeleteAdminRequest) (*auth.DeleteAdminResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	// TODO: 詳細の実装
	res := &auth.DeleteAdminResponse{}
	return res, nil
}
