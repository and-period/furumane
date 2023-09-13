package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/uuid"
	"github.com/and-period/furumane/proto/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func TestSignUpAdmin(t *testing.T) {
	t.Parallel()
	adminID := uuid.New()
	admin := &entity.Admin{
		ID:           uuid.Base58Encode(adminID),
		CognitoID:    uuid.Base58Encode(adminID),
		ProviderType: entity.ProviderTypeEmail,
		Email:        "test@example.com",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.SignUpAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().Create(ctx, admin, gomock.Any()).Return(nil)
			},
			req: &auth.SignUpAdminRequest{
				Email:                "test@example.com",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.SignUpAdminResponse{
					AdminId: uuid.Base58Encode(adminID),
				},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.SignUpAdminRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name:  "unmatch password and password confirmation",
			setup: func(ctx context.Context, mocks *mocks) {},
			req: &auth.SignUpAdminRequest{
				Email:                "test@example.com",
				Password:             "password",
				PasswordConfirmation: "password-confirmation",
			},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().Create(ctx, admin, gomock.Any()).Return(assert.AnError)
			},
			req: &auth.SignUpAdminRequest{
				Email:                "test@example.com",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.SignUpAdmin(ctx, tt.req)
		}, withUUID(adminID)))
	}
}

func TestVerifyAdmin(t *testing.T) {
	t.Parallel()
	now := time.Now()
	admin := &entity.Admin{
		ID:           "admin-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeEmail,
		Email:        "test@example.com",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.VerifyAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().Get(ctx, "admin-id", "cognito_id", "verified_at").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmSignUp(ctx, "cognito-id", "verify-code").Return(nil)
			},
			req: &auth.VerifyAdminRequest{
				AdminId:    "admin-id",
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.VerifyAdminResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.VerifyAdminRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().Get(ctx, "admin-id", "cognito_id", "verified_at").Return(nil, assert.AnError)
			},
			req: &auth.VerifyAdminRequest{
				AdminId:    "admin-id",
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "already verified",
			setup: func(ctx context.Context, mocks *mocks) {
				admin := &entity.Admin{VerifiedAt: now}
				mocks.db.admin.EXPECT().Get(ctx, "admin-id", "cognito_id", "verified_at").Return(admin, nil)
			},
			req: &auth.VerifyAdminRequest{
				AdminId:    "admin-id",
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: codes.FailedPrecondition,
			},
		},
		{
			name: "failed to confirm sign up",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().Get(ctx, "admin-id", "cognito_id", "verified_at").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmSignUp(ctx, "cognito-id", "verify-code").Return(assert.AnError)
			},
			req: &auth.VerifyAdminRequest{
				AdminId:    "admin-id",
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.VerifyAdmin(ctx, tt.req)
		}))
	}
}

func TestSignUpAdminWithOAuth(t *testing.T) {
	t.Parallel()
	adminID := uuid.New()
	auser := &cognito.AuthUser{
		Username:    "cognito-id",
		Email:       "test@example.com",
		PhoneNumber: "",
	}
	admin := &entity.Admin{
		ID:           uuid.Base58Encode(adminID),
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeOAuth,
		Email:        "test@example.com",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.SignUpAdminWithOAuthRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUser(ctx, "access-token").Return(auser, nil)
				mocks.db.admin.EXPECT().Create(ctx, admin, gomock.Any()).Return(nil)
			},
			req: &auth.SignUpAdminWithOAuthRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.SignUpAdminWithOAuthResponse{
					Admin: &auth.Admin{
						Id:           uuid.Base58Encode(adminID),
						ProviderType: auth.ProviderType_OAUTH,
						Email:        "test@example.com",
						CreatedAt:    time.Time{}.Unix(),
						UpdatedAt:    time.Time{}.Unix(),
					},
				},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.SignUpAdminWithOAuthRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUser(ctx, "access-token").Return(nil, assert.AnError)
			},
			req: &auth.SignUpAdminWithOAuthRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUser(ctx, "access-token").Return(auser, nil)
				mocks.db.admin.EXPECT().Create(ctx, admin, gomock.Any()).Return(assert.AnError)
			},
			req: &auth.SignUpAdminWithOAuthRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.SignUpAdminWithOAuth(ctx, tt.req)
		}, withUUID(adminID)))
	}
}

func TestGetAdmin(t *testing.T) {
	t.Parallel()
	admin := &entity.Admin{
		ID:           "admin-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeOAuth,
		Email:        "test@example.com",
		CreatedAt:    current,
		UpdatedAt:    current,
		VerifiedAt:   current,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.GetAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(admin, nil)
			},
			req: &auth.GetAdminRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.GetAdminResponse{
					Admin: &auth.Admin{
						Id:           "admin-id",
						ProviderType: auth.ProviderType_OAUTH,
						Email:        "test@example.com",
						CreatedAt:    current.Unix(),
						UpdatedAt:    current.Unix(),
					},
				},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.GetAdminRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			req: &auth.GetAdminRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(nil, assert.AnError)
			},
			req: &auth.GetAdminRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "unverified",
			setup: func(ctx context.Context, mocks *mocks) {
				admin := &entity.Admin{VerifiedAt: time.Time{}}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(admin, nil)
			},
			req: &auth.GetAdminRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.Unauthenticated,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.GetAdmin(ctx, tt.req)
		}))
	}
}

func TestUpdateAdminEmail(t *testing.T) {
	t.Parallel()
	admin := &entity.Admin{
		ID:           "admin-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeEmail,
		Email:        "hoge@example.com",
		CreatedAt:    current,
		UpdatedAt:    current,
		VerifiedAt:   current,
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: "access-token",
		Username:    "cognito-id",
		OldEmail:    "hoge@example.com",
		NewEmail:    "test@example.com",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.UpdateAdminEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(ctx, params).Return(nil)
			},
			req: &auth.UpdateAdminEmailRequest{
				AccessToken: "access-token",
				Email:       "test@example.com",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.UpdateAdminEmailResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.UpdateAdminEmailRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			req: &auth.UpdateAdminEmailRequest{
				AccessToken: "access-token",
				Email:       "test@example.com",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(nil, assert.AnError)
			},
			req: &auth.UpdateAdminEmailRequest{
				AccessToken: "access-token",
				Email:       "test@example.com",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to change email",
			setup: func(ctx context.Context, mocks *mocks) {
				admin := &entity.Admin{ProviderType: entity.ProviderTypeOAuth}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(admin, nil)
			},
			req: &auth.UpdateAdminEmailRequest{
				AccessToken: "access-token",
				Email:       "test@example.com",
			},
			expect: &testResponse{
				code: codes.FailedPrecondition,
			},
		},
		{
			name: "failed to change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(ctx, params).Return(assert.AnError)
			},
			req: &auth.UpdateAdminEmailRequest{
				AccessToken: "access-token",
				Email:       "test@example.com",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.UpdateAdminEmail(ctx, tt.req)
		}))
	}
}

func TestVerifyAdminEmail(t *testing.T) {
	t.Parallel()
	admin := &entity.Admin{
		ID:           "admin-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeEmail,
		Email:        "test@example.com",
		CreatedAt:    current,
		UpdatedAt:    current,
		VerifiedAt:   current,
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: "access-token",
		Username:    "cognito-id",
		VerifyCode:  "verify-code",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.VerifyAdminEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test@example.com", nil)
				mocks.db.admin.EXPECT().UpdateEmail(ctx, "admin-id", "test@example.com").Return(nil)
			},
			req: &auth.VerifyAdminEmailRequest{
				AccessToken: "access-token",
				VerifyCode:  "verify-code",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.VerifyAdminEmailResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.VerifyAdminEmailRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			req: &auth.VerifyAdminEmailRequest{
				AccessToken: "access-token",
				VerifyCode:  "verify-code",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(nil, assert.AnError)
			},
			req: &auth.VerifyAdminEmailRequest{
				AccessToken: "access-token",
				VerifyCode:  "verify-code",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to confirm change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("", assert.AnError)
			},
			req: &auth.VerifyAdminEmailRequest{
				AccessToken: "access-token",
				VerifyCode:  "verify-code",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to update email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(ctx, "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test@example.com", nil)
				mocks.db.admin.EXPECT().UpdateEmail(ctx, "admin-id", "test@example.com").Return(assert.AnError)
			},
			req: &auth.VerifyAdminEmailRequest{
				AccessToken: "access-token",
				VerifyCode:  "verify-code",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.VerifyAdminEmail(ctx, tt.req)
		}))
	}
}

func TestUpdateAdminPassword(t *testing.T) {
	t.Parallel()
	params := &cognito.ChangePasswordParams{
		AccessToken: "access-token",
		OldPassword: "password",
		NewPassword: "password",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.UpdateAdminPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().ChangePassword(ctx, params).Return(nil)
			},
			req: &auth.UpdateAdminPasswordRequest{
				AccessToken:          "access-token",
				OldPassword:          "password",
				NewPassword:          "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.UpdateAdminPasswordResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.UpdateAdminPasswordRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name:  "unmatch password and password confirmation",
			setup: func(ctx context.Context, mocks *mocks) {},
			req: &auth.UpdateAdminPasswordRequest{
				AccessToken:          "access-token",
				OldPassword:          "password",
				NewPassword:          "password",
				PasswordConfirmation: "password-confirmation",
			},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to change password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().ChangePassword(ctx, params).Return(assert.AnError)
			},
			req: &auth.UpdateAdminPasswordRequest{
				AccessToken:          "access-token",
				OldPassword:          "password",
				NewPassword:          "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.UpdateAdminPassword(ctx, tt.req)
		}))
	}
}

func TestForgotAdminPassword(t *testing.T) {
	t.Parallel()
	admin := &entity.Admin{
		CognitoID: "cognito-id",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.ForgotAdminPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(ctx, "test@example.com", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ForgotPassword(ctx, "cognito-id").Return(nil)
			},
			req: &auth.ForgotAdminPasswordRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.ForgotAdminPasswordResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.ForgotAdminPasswordRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(ctx, "test@example.com", "cognito_id").Return(nil, assert.AnError)
			},
			req: &auth.ForgotAdminPasswordRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to forgot password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(ctx, "test@example.com", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ForgotPassword(ctx, "cognito-id").Return(assert.AnError)
			},
			req: &auth.ForgotAdminPasswordRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.ForgotAdminPassword(ctx, tt.req)
		}))
	}
}

func TestResetAdminPassword(t *testing.T) {
	t.Parallel()
	admin := &entity.Admin{
		CognitoID: "cognito-id",
	}
	params := &cognito.ConfirmForgotPasswordParams{
		Username:    "cognito-id",
		VerifyCode:  "verify-code",
		NewPassword: "password",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.ResetAdminPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(ctx, "test@example.com", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmForgotPassword(ctx, params).Return(nil)
			},
			req: &auth.ResetAdminPasswordRequest{
				Email:                "test@example.com",
				VerifyCode:           "verify-code",
				NewPassword:          "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.ResetAdminPasswordResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.ResetAdminPasswordRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name:  "unmatch password and password confirmation",
			setup: func(ctx context.Context, mocks *mocks) {},
			req: &auth.ResetAdminPasswordRequest{
				Email:                "test@example.com",
				VerifyCode:           "verify-code",
				NewPassword:          "password",
				PasswordConfirmation: "password-confirmation",
			},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(ctx, "test@example.com", "cognito_id").Return(nil, assert.AnError)
			},
			req: &auth.ResetAdminPasswordRequest{
				Email:                "test@example.com",
				VerifyCode:           "verify-code",
				NewPassword:          "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to confirm forgot password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(ctx, "test@example.com", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmForgotPassword(ctx, params).Return(assert.AnError)
			},
			req: &auth.ResetAdminPasswordRequest{
				Email:                "test@example.com",
				VerifyCode:           "verify-code",
				NewPassword:          "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.ResetAdminPassword(ctx, tt.req)
		}))
	}
}

func TestDeleteAdmin(t *testing.T) {
	t.Parallel()
	admin := &entity.Admin{
		ID:           "admin-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeEmail,
		Email:        "hoge@example.com",
		CreatedAt:    current,
		UpdatedAt:    current,
		VerifiedAt:   current,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.DeleteAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
				mocks.db.admin.EXPECT().Delete(ctx, "admin-id", gomock.Any()).Return(nil)
			},
			req: &auth.DeleteAdminRequest{
				AdminId: "admin-id",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.DeleteAdminResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.DeleteAdminRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().Get(ctx, "admin-id").Return(nil, assert.AnError)
			},
			req: &auth.DeleteAdminRequest{
				AdminId: "admin-id",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
				mocks.db.admin.EXPECT().Delete(ctx, "admin-id", gomock.Any()).Return(assert.AnError)
			},
			req: &auth.DeleteAdminRequest{
				AdminId: "admin-id",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.DeleteAdmin(ctx, tt.req)
		}))
	}
}
