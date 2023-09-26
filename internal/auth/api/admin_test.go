package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/internal/auth/request"
	"github.com/and-period/furumane/internal/auth/response"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSignUpAdmin(t *testing.T) {
	t.Parallel()
	adminID := uuid.New()
	admin := &entity.Admin{
		ID:           uuid.Base58Encode(adminID),
		CognitoID:    uuid.Base58Encode(adminID),
		ProviderType: entity.ProviderTypeEmail,
		Email:        "test@example.com",
		PhoneNumber:  "09012341234",
	}
	tests := []struct {
		name   string
		setup  func(mocks *mocks)
		req    *request.SignUpAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Create(gomock.Any(), admin, gomock.Any()).Return(nil)
			},
			req: &request.SignUpAdminRequest{
				Email:                "test@example.com",
				PhoneNumber:          "09012341234",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.SignUpAdminResponse{
					AdminID: uuid.Base58Encode(adminID),
				},
			},
		},
		{
			name:  "invalid argument",
			setup: func(mocks *mocks) {},
			req:   &request.SignUpAdminRequest{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to create admin",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Create(gomock.Any(), admin, gomock.Any()).Return(assert.AnError)
			},
			req: &request.SignUpAdminRequest{
				Email:                "test@example.com",
				PhoneNumber:          "09012341234",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin"
			testPost(t, tt.setup, tt.expect, path, tt.req, withUUID(adminID))
		})
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
		setup  func(mocks *mocks)
		req    *request.VerifyAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Get(gomock.Any(), "admin-id", "cognito_id", "verified_at").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmSignUp(gomock.Any(), "cognito-id", "verify-code").Return(nil)
			},
			req: &request.VerifyAdminRequest{
				AdminID:    "admin-id",
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name:  "invalid argument",
			setup: func(mocks *mocks) {},
			req:   &request.VerifyAdminRequest{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get admin",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Get(gomock.Any(), "admin-id", "cognito_id", "verified_at").Return(nil, assert.AnError)
			},
			req: &request.VerifyAdminRequest{
				AdminID:    "admin-id",
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "already verified",
			setup: func(mocks *mocks) {
				admin := &entity.Admin{VerifiedAt: now}
				mocks.db.admin.EXPECT().Get(gomock.Any(), "admin-id", "cognito_id", "verified_at").Return(admin, nil)
			},
			req: &request.VerifyAdminRequest{
				AdminID:    "admin-id",
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: http.StatusPreconditionFailed,
			},
		},
		{
			name: "failed to confirm sign up",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Get(gomock.Any(), "admin-id", "cognito_id", "verified_at").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmSignUp(gomock.Any(), "cognito-id", "verify-code").Return(assert.AnError)
			},
			req: &request.VerifyAdminRequest{
				AdminID:    "admin-id",
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/verified"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
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
		setup  func(mocks *mocks)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUser(gomock.Any(), "access-token").Return(auser, nil)
				mocks.db.admin.EXPECT().Create(gomock.Any(), admin, gomock.Any()).Return(nil)
				mocks.db.admin.EXPECT().UpdateVerifiedAt(gomock.Any(), uuid.Base58Encode(adminID)).Return(nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.SignUpAdminWithOAuthResponse{
					Admin: &response.Admin{
						ID:           uuid.Base58Encode(adminID),
						ProviderType: entity.ProviderTypeOAuth,
						Email:        "test@example.com",
						CreatedAt:    time.Time{},
						UpdatedAt:    time.Time{},
					},
				},
			},
		},
		{
			name: "failed to get user",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUser(gomock.Any(), "access-token").Return(nil, assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to create admin",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUser(gomock.Any(), "access-token").Return(auser, nil)
				mocks.db.admin.EXPECT().Create(gomock.Any(), admin, gomock.Any()).Return(assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to update verified at",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUser(gomock.Any(), "access-token").Return(auser, nil)
				mocks.db.admin.EXPECT().Create(gomock.Any(), admin, gomock.Any()).Return(nil)
				mocks.db.admin.EXPECT().UpdateVerifiedAt(gomock.Any(), uuid.Base58Encode(adminID)).Return(assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/oauth"
			testPost(t, tt.setup, tt.expect, path, nil, withUUID(adminID))
		})
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
		name    string
		setup   func(mocks *mocks)
		adminID string
		expect  *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Get(gomock.Any(), "admin-id").Return(admin, nil)
			},
			adminID: "admin-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.GetAdminResponse{
					Admin: &response.Admin{
						ID:           "admin-id",
						ProviderType: entity.ProviderTypeOAuth,
						Email:        "test@example.com",
						CreatedAt:    current,
						UpdatedAt:    current,
					},
				},
			},
		},
		{
			name: "failed to get admin",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Get(gomock.Any(), "admin-id").Return(nil, assert.AnError)
			},
			adminID: "admin-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/admin/%s"
			path := fmt.Sprintf(format, tt.adminID)
			testGet(t, tt.setup, tt.expect, path)
		})
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
		setup  func(mocks *mocks)
		req    *request.UpdateAdminEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(gomock.Any(), params).Return(nil)
			},
			req: &request.UpdateAdminEmailRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name:  "invalid argument",
			setup: func(mocks *mocks) {},
			req:   &request.UpdateAdminEmailRequest{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get username",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("", assert.AnError)
			},
			req: &request.UpdateAdminEmailRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get admin",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(nil, assert.AnError)
			},
			req: &request.UpdateAdminEmailRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "not allow provider type",
			setup: func(mocks *mocks) {
				admin := &entity.Admin{ProviderType: entity.ProviderTypeOAuth}
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(admin, nil)
			},
			req: &request.UpdateAdminEmailRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: http.StatusPreconditionFailed,
			},
		},
		{
			name: "failed to change email",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(gomock.Any(), params).Return(assert.AnError)
			},
			req: &request.UpdateAdminEmailRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/email"
			testPut(t, tt.setup, tt.expect, path, tt.req)
		})
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
		setup  func(mocks *mocks)
		req    *request.VerifyAdminEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(gomock.Any(), params).Return("test@example.com", nil)
				mocks.db.admin.EXPECT().UpdateEmail(gomock.Any(), "admin-id", "test@example.com").Return(nil)
			},
			req: &request.VerifyAdminEmailRequest{
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name:  "invalid argument",
			setup: func(mocks *mocks) {},
			req:   &request.VerifyAdminEmailRequest{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get username",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("", assert.AnError)
			},
			req: &request.VerifyAdminEmailRequest{
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get admin",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(nil, assert.AnError)
			},
			req: &request.VerifyAdminEmailRequest{
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to confirm change email",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(gomock.Any(), params).Return("", assert.AnError)
			},
			req: &request.VerifyAdminEmailRequest{
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to update email",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(gomock.Any(), params).Return("test@example.com", nil)
				mocks.db.admin.EXPECT().UpdateEmail(gomock.Any(), "admin-id", "test@example.com").Return(assert.AnError)
			},
			req: &request.VerifyAdminEmailRequest{
				VerifyCode: "verify-code",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/email/verified"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
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
		setup  func(mocks *mocks)
		req    *request.UpdateAdminPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().ChangePassword(gomock.Any(), params).Return(nil)
			},
			req: &request.UpdateAdminPasswordRequest{
				OldPassword:          "password",
				NewPassword:          "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name:  "invalid argument",
			setup: func(mocks *mocks) {},
			req:   &request.UpdateAdminPasswordRequest{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "unmatch password and password confirmation",
			setup: func(mocks *mocks) {},
			req: &request.UpdateAdminPasswordRequest{
				OldPassword:          "password",
				NewPassword:          "password",
				PasswordConfirmation: "password-confirmation",
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to change password",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().ChangePassword(gomock.Any(), params).Return(assert.AnError)
			},
			req: &request.UpdateAdminPasswordRequest{
				OldPassword:          "password",
				NewPassword:          "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/password"
			testPut(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestForgotAdminPassword(t *testing.T) {
	t.Parallel()
	admin := &entity.Admin{
		CognitoID: "cognito-id",
	}
	tests := []struct {
		name   string
		setup  func(mocks *mocks)
		req    *request.ForgotAdminPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(gomock.Any(), "test@example.com", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ForgotPassword(gomock.Any(), "cognito-id").Return(nil)
			},
			req: &request.ForgotAdminPasswordRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name:  "invalid argument",
			setup: func(mocks *mocks) {},
			req:   &request.ForgotAdminPasswordRequest{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get admin",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(gomock.Any(), "test@example.com", "cognito_id").Return(nil, assert.AnError)
			},
			req: &request.ForgotAdminPasswordRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to forgot password",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(gomock.Any(), "test@example.com", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ForgotPassword(gomock.Any(), "cognito-id").Return(assert.AnError)
			},
			req: &request.ForgotAdminPasswordRequest{
				Email: "test@example.com",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/password/forgot"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
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
		setup  func(mocks *mocks)
		req    *request.ResetAdminPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(gomock.Any(), "test@example.com", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmForgotPassword(gomock.Any(), params).Return(nil)
			},
			req: &request.ResetAdminPasswordRequest{
				Email:                "test@example.com",
				VerifyCode:           "verify-code",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name:  "invalid argument",
			setup: func(mocks *mocks) {},
			req:   &request.ResetAdminPasswordRequest{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "unmatch password and password confirmation",
			setup: func(mocks *mocks) {},
			req: &request.ResetAdminPasswordRequest{
				Email:                "test@example.com",
				VerifyCode:           "verify-code",
				Password:             "password",
				PasswordConfirmation: "password-confirmation",
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get admin",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(gomock.Any(), "test@example.com", "cognito_id").Return(nil, assert.AnError)
			},
			req: &request.ResetAdminPasswordRequest{
				Email:                "test@example.com",
				VerifyCode:           "verify-code",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to confirm forgot password",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().GetByEmail(gomock.Any(), "test@example.com", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmForgotPassword(gomock.Any(), params).Return(assert.AnError)
			},
			req: &request.ResetAdminPasswordRequest{
				Email:                "test@example.com",
				VerifyCode:           "verify-code",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/password/reset"
			testPut(t, tt.setup, tt.expect, path, tt.req)
		})
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
		name    string
		setup   func(mocks *mocks)
		adminID string
		expect  *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Get(gomock.Any(), "admin-id").Return(admin, nil)
				mocks.db.admin.EXPECT().Delete(gomock.Any(), "admin-id", gomock.Any()).Return(nil)
			},
			adminID: "admin-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to get admin",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Get(gomock.Any(), "admin-id").Return(nil, assert.AnError)
			},
			adminID: "admin-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to delete",
			setup: func(mocks *mocks) {
				mocks.db.admin.EXPECT().Get(gomock.Any(), "admin-id").Return(admin, nil)
				mocks.db.admin.EXPECT().Delete(gomock.Any(), "admin-id", gomock.Any()).Return(assert.AnError)
			},
			adminID: "admin-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/admin/%s"
			path := fmt.Sprintf(format, tt.adminID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
