package api

import (
	"net/http"
	"testing"

	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/internal/auth/request"
	"github.com/and-period/furumane/internal/auth/response"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSignInAdmin(t *testing.T) {
	t.Parallel()
	result := &cognito.AuthResult{
		IDToken:      "id-token",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}
	admin := &entity.Admin{
		ID:           "admin-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeEmail,
		CreatedAt:    current,
		UpdatedAt:    current,
	}
	tests := []struct {
		name   string
		setup  func(mocks *mocks)
		req    *request.SignInAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(gomock.Any(), "test@example.com", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(admin, nil)
			},
			req: &request.SignInAdminRequest{
				Key:      "test@example.com",
				Password: "password",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.SignInAdminResponse{
					AdminAuth: &response.AdminAuth{
						AdminID:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    3600,
					},
				},
			},
		},
		{
			name:  "bad request",
			setup: func(mocks *mocks) {},
			req:   &request.SignInAdminRequest{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to sign in",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(gomock.Any(), "test@example.com", "password").Return(nil, assert.AnError)
			},
			req: &request.SignInAdminRequest{
				Key:      "test@example.com",
				Password: "password",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get username",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(gomock.Any(), "test@example.com", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("", assert.AnError)
			},
			req: &request.SignInAdminRequest{
				Key:      "test@example.com",
				Password: "password",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get admin by cognito id",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(gomock.Any(), "test@example.com", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(nil, assert.AnError)
			},
			req: &request.SignInAdminRequest{
				Key:      "test@example.com",
				Password: "password",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/auth"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestSignOutAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(mocks *mocks)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().SignOut(gomock.Any(), "access-token").Return(nil)
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to sign out",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().SignOut(gomock.Any(), "access-token").Return(assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/auth"
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetAdminAuth(t *testing.T) {
	t.Parallel()
	admin := &entity.Admin{
		ID:           "admin-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeEmail,
		CreatedAt:    current,
		UpdatedAt:    current,
	}
	tests := []struct {
		name   string
		setup  func(mocks *mocks)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(admin, nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.SignInAdminResponse{
					AdminAuth: &response.AdminAuth{
						AdminID:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "",
						ExpiresIn:    0,
					},
				},
			},
		},
		{
			name: "failed to get username",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("", assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get admin by cognito id",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(nil, assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/auth"
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestRefreshAdminToken(t *testing.T) {
	t.Parallel()
	result := &cognito.AuthResult{
		IDToken:      "id-token",
		AccessToken:  "access-token",
		RefreshToken: "",
		ExpiresIn:    3600,
	}
	admin := &entity.Admin{
		ID:           "admin-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeOAuth,
		CreatedAt:    current,
		UpdatedAt:    current,
	}
	tests := []struct {
		name   string
		setup  func(mocks *mocks)
		req    *request.RefreshAdminTokenRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(gomock.Any(), "refresh-token").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(admin, nil)
			},
			req: &request.RefreshAdminTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.RefreshAdminTokenResponse{
					AdminAuth: &response.AdminAuth{
						AdminID:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "",
						ExpiresIn:    3600,
					},
				},
			},
		},
		{
			name:  "invalid argument",
			setup: func(mocks *mocks) {},
			req:   &request.RefreshAdminTokenRequest{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to refresh token",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(gomock.Any(), "refresh-token").Return(nil, assert.AnError)
			},
			req: &request.RefreshAdminTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get username",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(gomock.Any(), "refresh-token").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("", assert.AnError)
			},
			req: &request.RefreshAdminTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get admin by cognito id",
			setup: func(mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(gomock.Any(), "refresh-token").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(gomock.Any(), "access-token").Return("cognito-id", nil)
				mocks.db.admin.EXPECT().GetByCognitoID(gomock.Any(), "cognito-id").Return(nil, assert.AnError)
			},
			req: &request.RefreshAdminTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/admin/auth/refresh"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}
