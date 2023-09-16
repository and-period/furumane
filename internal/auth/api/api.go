package api

import (
	"sync"
	"time"

	"github.com/and-period/furumane/internal/auth/database"
	"github.com/and-period/furumane/internal/auth/response"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/jst"
	"github.com/and-period/furumane/pkg/uuid"
	"github.com/and-period/furumane/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller interface {
	Routes(ctx *gin.RouterGroup) // エンドポイント一覧の定義
}

type Params struct {
	WaitGroup *sync.WaitGroup
	Database  *database.Database
	AdminAuth cognito.Client
	UserAuth  cognito.Client
}

type controller struct {
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	sharedGroup *singleflight.Group
	db          *database.Database
	validator   validator.Validator
	adminAuth   cognito.Client
	userAuth    cognito.Client
	uuid        func() string
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

func NewController(params *Params, opts ...Option) Controller {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &controller{
		now:         jst.Now,
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		sharedGroup: &singleflight.Group{},
		db:          params.Database,
		validator:   validator.NewValidator(),
		adminAuth:   params.AdminAuth,
		userAuth:    params.UserAuth,
		uuid:        uuid.New,
	}
}

func (c *controller) Routes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	{
		c.adminAuthRoutes(admin)
		c.adminRoutes(admin)
	}
}

func (c *controller) bind(ctx *gin.Context, req interface{}) error {
	if err := ctx.BindJSON(req); err != nil {
		return err
	}
	return c.validator.Struct(req)
}

func httpError(ctx *gin.Context, err error) {
	res, status := response.NewErrorResponse(err)
	ctx.AbortWithStatusJSON(status, res)
}

func badRequest(ctx *gin.Context, format string, args ...interface{}) {
	httpError(ctx, status.Errorf(codes.InvalidArgument, format, args...))
}

func unauthorized(ctx *gin.Context, format string, args ...interface{}) {
	httpError(ctx, status.Errorf(codes.Unauthenticated, format, args...))
}

func preconditionFailed(ctx *gin.Context, format string, args ...interface{}) {
	httpError(ctx, status.Errorf(codes.FailedPrecondition, format, args...))
}
