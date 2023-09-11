package service

import (
	"context"
	"sync"
	"testing"
	"time"

	mock_cognito "github.com/and-period/furumane/mock/pkg/cognito"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/jst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type mocks struct {
	adminAuth *mock_cognito.MockClient
	userAuth  *mock_cognito.MockClient
}

type testResponse struct {
	code codes.Code
	body proto.Message
}

type testOptions struct {
	now func() time.Time
}

type testOption func(opts *testOptions)

func withNow(now time.Time) testOption {
	return func(opts *testOptions) {
		opts.now = func() time.Time {
			return now
		}
	}
}

type grpcCaller func(ctx context.Context, service *service) (proto.Message, error)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		adminAuth: mock_cognito.NewMockClient(ctrl),
		userAuth:  mock_cognito.NewMockClient(ctrl),
	}
}

func newService(mocks *mocks, opts ...testOption) *service {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		AdminAuth: mocks.adminAuth,
		UserAuth:  mocks.userAuth,
	}
	service := NewService(params).(*service)
	service.now = func() time.Time {
		return dopts.now()
	}
	return service
}

func testGRPC(
	setup func(ctx context.Context, mocks *mocks),
	expect *testResponse,
	grpcFn grpcCaller,
	opts ...testOption,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mocks := newMocks(ctrl)

		srv := newService(mocks, opts...)
		setup(ctx, mocks)

		res, err := grpcFn(ctx, srv)
		switch expect.code {
		case codes.OK:
			require.NoError(t, err)
		default:
			require.Error(t, err)
			status, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, expect.code, status.Code(), status.Code().String())
		}
		if expect.body != nil {
			require.Equal(t, expect.body, res)
		}
		srv.waitGroup.Wait()
	}
}

func TestService(t *testing.T) {
	t.Parallel()
	srv := NewService(&Params{}, WithLogger(zap.NewNop()))
	assert.NotNil(t, srv)
}

func TestGRPCError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  error
		expect codes.Code
	}{
		// common
		{
			name:   "ok",
			input:  nil,
			expect: codes.OK,
		},
		{
			name:   "gRPC error",
			input:  status.Error(codes.InvalidArgument, "invalid argument"),
			expect: codes.InvalidArgument,
		},
		{
			name:   "unknown error",
			input:  assert.AnError,
			expect: codes.Internal,
		},
		// context error
		{
			name:   "context canceled",
			input:  context.Canceled,
			expect: codes.Canceled,
		},
		{
			name:   "context deadline exceeded",
			input:  context.DeadlineExceeded,
			expect: codes.DeadlineExceeded,
		},
		// cognito error
		{
			name:   "cognito canceled",
			input:  cognito.ErrCanceled,
			expect: codes.Canceled,
		},
		{
			name:   "cognito invalid argument",
			input:  cognito.ErrInvalidArgument,
			expect: codes.InvalidArgument,
		},
		{
			name:   "cognito timeout",
			input:  cognito.ErrTimeout,
			expect: codes.DeadlineExceeded,
		},
		{
			name:   "cognito already exists",
			input:  cognito.ErrAlreadyExists,
			expect: codes.AlreadyExists,
		},
		{
			name:   "cognito resource exhausted",
			input:  cognito.ErrResourceExhausted,
			expect: codes.ResourceExhausted,
		},
		{
			name:   "cognito unauthenticated",
			input:  cognito.ErrUnauthenticated,
			expect: codes.Unauthenticated,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := gRPCError(tt.input)
			assert.Equal(t, tt.expect, status.Code(err), err)
		})
	}
}

func TestCustomError(t *testing.T) {
	t.Parallel()
	t.Run("invalid argument", func(t *testing.T) {
		err := invalidArgument(assert.AnError)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}
