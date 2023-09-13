package service

import (
	"context"
	"errors"

	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/uuid"
	"github.com/and-period/furumane/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) SignUpAdmin(
	ctx context.Context, req *auth.SignUpAdminRequest,
) (*auth.SignUpAdminResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	if req.Password != req.PasswordConfirmation {
		return nil, invalidArgument(errors.New("service: unmatch password and password confirmation"))
	}
	cognitoID := uuid.Base58Encode(s.uuid())
	params := &entity.AdminParams{
		AdminID:      uuid.Base58Encode(s.uuid()),
		CognitID:     cognitoID,
		ProviderType: entity.ProviderTypeEmail,
		Email:        req.Email,
	}
	admin := entity.NewAdmin(params)
	fn := func(ctx context.Context) error {
		params := &cognito.SignUpParams{
			Username: cognitoID,
			Email:    req.Email,
			Password: req.Password,
		}
		return s.adminAuth.SignUp(ctx, params)
	}
	if err := s.db.Admin.Create(ctx, admin, fn); err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.SignUpAdminResponse{
		AdminId: admin.ID,
	}
	return res, nil
}

func (s *service) VerifyAdmin(ctx context.Context, req *auth.VerifyAdminRequest) (*auth.VerifyAdminResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	admin, err := s.db.Admin.Get(ctx, req.AdminId, "cognito_id", "verified_at")
	if err != nil {
		return nil, gRPCError(err)
	}
	if !admin.VerifiedAt.IsZero() {
		return nil, status.Error(codes.FailedPrecondition, "this admin is already verified")
	}
	if err := s.adminAuth.ConfirmSignUp(ctx, admin.CognitoID, req.VerifyCode); err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.VerifyAdminResponse{}
	return res, nil
}

func (s *service) SignUpAdminWithOAuth(
	ctx context.Context, req *auth.SignUpAdminWithOAuthRequest,
) (*auth.SignUpAdminWithOAuthResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	au, err := s.adminAuth.GetUser(ctx, req.AccessToken)
	if err != nil {
		return nil, gRPCError(err)
	}
	params := &entity.AdminParams{
		AdminID:      uuid.Base58Encode(s.uuid()),
		CognitID:     au.Username,
		ProviderType: entity.ProviderTypeOAuth,
		Email:        au.Email,
	}
	admin := entity.NewAdmin(params)
	fn := func(ctx context.Context) error {
		return nil // Cognitoへはすでに登録済みのため何もしない
	}
	if err := s.db.Admin.Create(ctx, admin, fn); err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.SignUpAdminWithOAuthResponse{
		Admin: admin.Proto(),
	}
	return res, nil
}

func (s *service) GetAdmin(ctx context.Context, req *auth.GetAdminRequest) (*auth.GetAdminResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, req.AccessToken)
	if err != nil {
		return nil, gRPCError(err)
	}
	admin, err := s.db.Admin.GetByCognitoID(ctx, username)
	if err != nil {
		return nil, gRPCError(err)
	}
	if admin.VerifiedAt.IsZero() {
		return nil, status.Error(codes.Unauthenticated, "not verified")
	}
	res := &auth.GetAdminResponse{
		Admin: admin.Proto(),
	}
	return res, nil
}

func (s *service) UpdateAdminEmail(ctx context.Context, req *auth.UpdateAdminEmailRequest) (*auth.UpdateAdminEmailResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, req.AccessToken)
	if err != nil {
		return nil, gRPCError(err)
	}
	admin, err := s.db.Admin.GetByCognitoID(ctx, username)
	if err != nil {
		return nil, gRPCError(err)
	}
	if admin.ProviderType != entity.ProviderTypeEmail {
		return nil, status.Error(codes.FailedPrecondition, "not allow provider type to change email")
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: req.AccessToken,
		Username:    username,
		OldEmail:    admin.Email,
		NewEmail:    req.Email,
	}
	if err := s.adminAuth.ChangeEmail(ctx, params); err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.UpdateAdminEmailResponse{}
	return res, nil
}

func (s *service) VerifyAdminEmail(ctx context.Context, req *auth.VerifyAdminEmailRequest) (*auth.VerifyAdminEmailResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, req.AccessToken)
	if err != nil {
		return nil, gRPCError(err)
	}
	admin, err := s.db.Admin.GetByCognitoID(ctx, username)
	if err != nil {
		return nil, gRPCError(err)
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: req.AccessToken,
		Username:    username,
		VerifyCode:  req.VerifyCode,
	}
	email, err := s.adminAuth.ConfirmChangeEmail(ctx, params)
	if err != nil {
		return nil, gRPCError(err)
	}
	if err := s.db.Admin.UpdateEmail(ctx, admin.ID, email); err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.VerifyAdminEmailResponse{}
	return res, nil
}

func (s *service) UpdateAdminPassword(
	ctx context.Context, req *auth.UpdateAdminPasswordRequest,
) (*auth.UpdateAdminPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	if req.NewPassword != req.PasswordConfirmation {
		return nil, invalidArgument(errors.New("service: unmatch password and password confirmation"))
	}
	params := &cognito.ChangePasswordParams{
		AccessToken: req.AccessToken,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	if err := s.adminAuth.ChangePassword(ctx, params); err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.UpdateAdminPasswordResponse{}
	return res, nil
}

func (s *service) ForgotAdminPassword(
	ctx context.Context, req *auth.ForgotAdminPasswordRequest,
) (*auth.ForgotAdminPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	admin, err := s.db.Admin.GetByEmail(ctx, req.Email, "cognito_id")
	if err != nil {
		return nil, gRPCError(err)
	}
	if err := s.adminAuth.ForgotPassword(ctx, admin.CognitoID); err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.ForgotAdminPasswordResponse{}
	return res, nil
}

func (s *service) ResetAdminPassword(ctx context.Context, req *auth.ResetAdminPasswordRequest) (*auth.ResetAdminPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	if req.NewPassword != req.PasswordConfirmation {
		return nil, invalidArgument(errors.New("service: unmatch password and password confirmation"))
	}
	admin, err := s.db.Admin.GetByEmail(ctx, req.Email, "cognito_id")
	if err != nil {
		return nil, gRPCError(err)
	}
	params := &cognito.ConfirmForgotPasswordParams{
		Username:    admin.CognitoID,
		VerifyCode:  req.VerifyCode,
		NewPassword: req.NewPassword,
	}
	if err := s.adminAuth.ConfirmForgotPassword(ctx, params); err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.ResetAdminPasswordResponse{}
	return res, nil
}

func (s *service) DeleteAdmin(ctx context.Context, req *auth.DeleteAdminRequest) (*auth.DeleteAdminResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, invalidArgument(err)
	}
	admin, err := s.db.Admin.Get(ctx, req.AdminId)
	if err != nil {
		return nil, gRPCError(err)
	}
	fn := func(ctx context.Context) error {
		return s.adminAuth.DeleteUser(ctx, admin.CognitoID)
	}
	if err := s.db.Admin.Delete(ctx, admin.ID, fn); err != nil {
		return nil, gRPCError(err)
	}
	res := &auth.DeleteAdminResponse{}
	return res, nil
}
