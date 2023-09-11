package server

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumane/internal/auth/service"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/jst"
	"github.com/and-period/furumane/proto/auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
)

type registry struct {
	appName   string
	env       string
	waitGroup *sync.WaitGroup
	service   auth.AuthServiceServer
}

type params struct {
	config    *config
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
	aws       aws.Config
	adminAuth cognito.Client
	userAuth  cognito.Client
	now       func() time.Time
}

//nolint:funlen
func newRegistry(ctx context.Context, conf *config, logger *zap.Logger) (*registry, error) {
	params := &params{
		config:    conf,
		logger:    logger,
		now:       jst.Now,
		waitGroup: &sync.WaitGroup{},
	}

	// AWS SDKの設定
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(conf.AWSRegion))
	if err != nil {
		return nil, err
	}
	params.aws = awscfg

	// Amazon Cognitoの設定
	adminAuthParams := &cognito.Params{
		UserPoolID:  conf.CognitoAdminPoolID,
		AppClientID: conf.CognitoAdminClientID,
	}
	params.adminAuth = cognito.NewClient(awscfg, adminAuthParams)
	userAuthParams := &cognito.Params{
		UserPoolID:  conf.CognitoUserPoolID,
		AppClientID: conf.CognitoUserClientID,
	}
	params.userAuth = cognito.NewClient(awscfg, userAuthParams, cognito.WithLogger(params.logger))

	// Serviceの設定
	srvParams := &service.Params{
		WaitGroup: params.waitGroup,
		AdminAuth: params.adminAuth,
		UserAuth:  params.userAuth,
	}
	return &registry{
		appName:   conf.AppName,
		env:       conf.Environment,
		waitGroup: params.waitGroup,
		service:   service.NewService(srvParams, service.WithLogger(logger)),
	}, nil
}
