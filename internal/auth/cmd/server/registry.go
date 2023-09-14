package server

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumane/internal/auth/database/mysql"
	"github.com/and-period/furumane/internal/auth/service"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/jst"
	apmysql "github.com/and-period/furumane/pkg/mysql"
	"github.com/and-period/furumane/pkg/secret"
	"github.com/and-period/furumane/proto/auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/zap"
)

type registry struct {
	appName   string
	env       string
	waitGroup *sync.WaitGroup
	service   auth.AuthServiceServer
}

type params struct {
	config     *config
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	aws        aws.Config
	secret     secret.Client
	db         *apmysql.Client
	adminAuth  cognito.Client
	userAuth   cognito.Client
	now        func() time.Time
	dbHost     string
	dbPort     string
	dbUsername string
	dbPassword string
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

	// AWS Secrets Managerの設定
	params.secret = secret.NewClient(awscfg)
	if err := getSecret(ctx, params); err != nil {
		return nil, err
	}

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

	// Databaseの設定
	params.db, err = newDatabase(params)
	if err != nil {
		return nil, err
	}

	// Serviceの設定
	srvParams := &service.Params{
		WaitGroup: params.waitGroup,
		Database:  mysql.NewDatabase(params.db),
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

func getSecret(ctx context.Context, p *params) error {
	// データベース認証情報の取得
	if p.config.DBSecretName == "" {
		p.dbHost = p.config.DBHost
		p.dbPort = p.config.DBPort
		p.dbUsername = p.config.DBUsername
		p.dbPassword = p.config.DBPassword
		return nil
	}
	secrets, err := p.secret.Get(ctx, p.config.DBSecretName)
	if err != nil {
		return err
	}
	p.dbHost = secrets["host"]
	p.dbPort = secrets["port"]
	p.dbUsername = secrets["username"]
	p.dbPassword = secrets["password"]
	return nil
}

func newDatabase(p *params) (*apmysql.Client, error) {
	params := &apmysql.Params{
		Socket:   p.config.DBSocket,
		Host:     p.dbHost,
		Port:     p.dbPort,
		Database: p.config.DBDatabase,
		Username: p.dbUsername,
		Password: p.dbPassword,
	}
	location, err := time.LoadLocation(p.config.DBTimeZone)
	if err != nil {
		return nil, err
	}
	cli, err := apmysql.NewClient(
		params,
		apmysql.WithLogger(p.logger),
		apmysql.WithNow(p.now),
		apmysql.WithTLS(p.config.DBEnabledTLS),
		apmysql.WithLocation(location),
	)
	if err != nil {
		return nil, err
	}
	tracer := telemetry.NewNrTracer(p.config.DBDatabase, p.dbHost, string(newrelic.DatastoreMySQL))
	if err := cli.DB.Use(tracer); err != nil {
		return nil, err
	}
	return cli, nil
}
