package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumane/internal/auth/database"
	mock_database "github.com/and-period/furumane/mock/auth/database"
	mock_cognito "github.com/and-period/furumane/mock/pkg/cognito"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/jst"
	"github.com/and-period/furumane/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var current = jst.Now()

type mocks struct {
	db        *dbmocks
	adminAuth *mock_cognito.MockClient
	userAuth  *mock_cognito.MockClient
}

type dbmocks struct {
	admin *mock_database.MockAdmin
}

type testResponse struct {
	code codes.Code
	body proto.Message
}

type testOptions struct {
	now  func() time.Time
	uuid func() string
}

type testOption func(opts *testOptions)

func withNow(now time.Time) testOption {
	return func(opts *testOptions) {
		opts.now = func() time.Time {
			return now
		}
	}
}

func withUUID(uuid string) testOption {
	return func(opts *testOptions) {
		opts.uuid = func() string {
			return uuid
		}
	}
}

type grpcCaller func(ctx context.Context, service *service) (proto.Message, error)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db:        newDBMocks(ctrl),
		adminAuth: mock_cognito.NewMockClient(ctrl),
		userAuth:  mock_cognito.NewMockClient(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbmocks {
	return &dbmocks{
		admin: mock_database.NewMockAdmin(ctrl),
	}
}

func newService(mocks *mocks, opts ...testOption) *service {
	dopts := &testOptions{
		now:  jst.Now,
		uuid: uuid.New,
	}
	for i := range opts {
		opts[i](dopts)
	}
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Database: &database.Database{
			Admin: mocks.db.admin,
		},
		AdminAuth: mocks.adminAuth,
		UserAuth:  mocks.userAuth,
	}
	service := NewService(params).(*service)
	service.now = func() time.Time {
		return dopts.now()
	}
	service.uuid = func() string {
		return dopts.uuid()
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
		{
			name:   "cognito not found",
			input:  cognito.ErrNotFound,
			expect: codes.Unauthenticated,
		},
		// database error
		{
			name:   "database deadline exceeded",
			input:  database.ErrDeadlineExceeded,
			expect: codes.DeadlineExceeded,
		},
		{
			name:   "database not found",
			input:  database.ErrNotFound,
			expect: codes.NotFound,
		},
		{
			name:   "database already exists",
			input:  database.ErrAlreadyExists,
			expect: codes.AlreadyExists,
		},
		{
			name:   "database failed precondition",
			input:  database.ErrFailedPrecondition,
			expect: codes.FailedPrecondition,
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
