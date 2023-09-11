package service

import (
	"context"
	"testing"

	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/proto/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func TestSignInAdmin(t *testing.T) {
	t.Parallel()
	result := &cognito.AuthResult{
		IDToken:      "id-token",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.SignInAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "test@example.com", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("admin-id", nil)
			},
			req: &auth.SignInAdminRequest{
				Key:      "test@example.com",
				Password: "password",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.SignInAdminResponse{
					Auth: &auth.AdminAuth{
						AdminId:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    3600,
					},
				},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.SignInAdminRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to sign in",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "test@example.com", "password").Return(nil, assert.AnError)
			},
			req: &auth.SignInAdminRequest{
				Key:      "test@example.com",
				Password: "password",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "test@example.com", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			req: &auth.SignInAdminRequest{
				Key:      "test@example.com",
				Password: "password",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.SignInAdmin(ctx, tt.req)
		}))
	}
}

func TestSignInAdminWithOAuth(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.SignInAdminWithOAuthRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("admin-id", nil)
			},
			req: &auth.SignInAdminWithOAuthRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.SignInAdminWithOAuthResponse{
					Auth: &auth.AdminAuth{
						AdminId:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "",
						ExpiresIn:    0,
					},
				},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.SignInAdminWithOAuthRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			req: &auth.SignInAdminWithOAuthRequest{
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
			return service.SignInAdminWithOAuth(ctx, tt.req)
		}))
	}
}

func TestSignOutAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.SignOutAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignOut(ctx, "access-token").Return(nil)
			},
			req: &auth.SignOutAdminRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.SignOutAdminResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.SignOutAdminRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to sign out",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignOut(ctx, "access-token").Return(assert.AnError)
			},
			req: &auth.SignOutAdminRequest{
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
			return service.SignOutAdmin(ctx, tt.req)
		}))
	}
}

func TestGetAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.GetAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("admin-id", nil)
			},
			req: &auth.GetAdminRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.GetAdminResponse{
					Auth: &auth.AdminAuth{
						AdminId:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "",
						ExpiresIn:    0,
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.GetAdmin(ctx, tt.req)
		}))
	}
}

func TestRefreshAdminToken(t *testing.T) {
	t.Parallel()
	result := &cognito.AuthResult{
		IDToken:      "",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.RefreshAdminTokenRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "refresh-token").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("admin-id", nil)
			},
			req: &auth.RefreshAdminTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.RefreshAdminTokenResponse{
					Auth: &auth.AdminAuth{
						AdminId:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    3600,
					},
				},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			req:   &auth.RefreshAdminTokenRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to refresh token",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "refresh-token").Return(nil, assert.AnError)
			},
			req: &auth.RefreshAdminTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "refresh-token").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			req: &auth.RefreshAdminTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.RefreshAdminToken(ctx, tt.req)
		}))
	}
}
