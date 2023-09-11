package service

import (
	"context"
	"testing"

	"github.com/and-period/furumane/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func TestSignUpAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.SignUpAdminRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			req: &auth.SignUpAdminRequest{
				Email:                "test@example.com",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.SignUpAdminResponse{},
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.SignUpAdmin(ctx, tt.req)
		}))
	}
}

func TestVerifyAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.VerifyAdminRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
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
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.SignUpAdminWithOAuthRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			req: &auth.SignUpAdminWithOAuthRequest{
				AccessToken: "access-token",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &auth.SignUpAdminWithOAuthResponse{},
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.SignUpAdminWithOAuth(ctx, tt.req)
		}))
	}
}

func TestUpdateAdminEmail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.UpdateAdminEmailRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
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
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.VerifyAdminEmailRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
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
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.UpdateAdminPasswordRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			req: &auth.UpdateAdminPasswordRequest{
				Email:                "test@example.com",
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
				Email:                "test@example.com",
				OldPassword:          "password",
				NewPassword:          "password",
				PasswordConfirmation: "password-confirmation",
			},
			expect: &testResponse{
				code: codes.InvalidArgument,
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
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.ForgotAdminPasswordRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
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
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.ResetAdminPasswordRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
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
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		req    *auth.DeleteAdminRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *service) (proto.Message, error) {
			return service.DeleteAdmin(ctx, tt.req)
		}))
	}
}
