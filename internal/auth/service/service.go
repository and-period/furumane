package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/jst"
	"github.com/and-period/furumane/proto/auth"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Params struct {
	WaitGroup *sync.WaitGroup
	AdminAuth cognito.Client
	UserAuth  cognito.Client
}

type service struct {
	auth.UnimplementedAuthServiceServer
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	sharedGroup *singleflight.Group
	adminAuth   cognito.Client
	userAuth    cognito.Client
}

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewService(params *Params, opts ...Option) auth.AuthServiceServer {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &service{
		now:         jst.Now,
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		sharedGroup: &singleflight.Group{},
		adminAuth:   params.AdminAuth,
		userAuth:    params.UserAuth,
	}
}

func gRPCError(err error) error {
	e, ok := status.FromError(err)
	if ok || e == nil {
		return err
	}
	switch {
	case errors.Is(err, context.Canceled),
		errors.Is(err, cognito.ErrCanceled):
		return status.Error(codes.Canceled, err.Error())
	case errors.Is(err, cognito.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, cognito.ErrTimeout):
		return status.Error(codes.DeadlineExceeded, err.Error())
	case errors.Is(err, cognito.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, cognito.ErrResourceExhausted):
		return status.Error(codes.ResourceExhausted, err.Error())
	case errors.Is(err, cognito.ErrUnauthenticated),
		errors.Is(err, cognito.ErrNotFound):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func invalidArgument(err error) error {
	return status.Error(codes.InvalidArgument, err.Error())
}
