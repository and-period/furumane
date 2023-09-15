package server

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumane/internal/auth/api"
	"github.com/and-period/furumane/internal/auth/database/mysql"
	"github.com/and-period/furumane/pkg/cognito"
	"github.com/and-period/furumane/pkg/jst"
	apmysql "github.com/and-period/furumane/pkg/mysql"
	"github.com/and-period/furumane/pkg/secret"
	"github.com/and-period/furumane/pkg/slack"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type registry struct {
	appName   string
	env       string
	debugMode bool
	waitGroup *sync.WaitGroup
	service   api.Controller
	newRelic  *newrelic.Application
	slack     slack.Client
}

type params struct {
	config          *config
	logger          *zap.Logger
	waitGroup       *sync.WaitGroup
	aws             aws.Config
	secret          secret.Client
	db              *apmysql.Client
	adminAuth       cognito.Client
	userAuth        cognito.Client
	newRelic        *newrelic.Application
	slack           slack.Client
	now             func() time.Time
	dbHost          string
	dbPort          string
	dbUsername      string
	dbPassword      string
	newRelicLicense string
	slackToken      string
	slackChannelID  string
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

	// New Relicの設定
	if params.newRelicLicense != "" {
		newrelicApp, err := newrelic.NewApplication(
			newrelic.ConfigEnabled(true),
			newrelic.ConfigAppName(conf.AppName),
			newrelic.ConfigLicense(params.newRelicLicense),
			newrelic.ConfigDistributedTracerEnabled(true),
		)
		if err != nil {
			return nil, err
		}
		params.newRelic = newrelicApp
	}

	// Slackの設定
	if params.slackToken != "" {
		slackParams := &slack.Params{
			Token:     params.slackToken,
			ChannelID: params.slackChannelID,
		}
		params.slack = slack.NewClient(slackParams, slack.WithLogger(logger))
	}

	// Serviceの設定
	apiParams := &api.Params{
		WaitGroup: params.waitGroup,
		Database:  mysql.NewDatabase(params.db),
		AdminAuth: params.adminAuth,
		UserAuth:  params.userAuth,
	}
	return &registry{
		appName:   conf.AppName,
		env:       conf.Environment,
		debugMode: conf.LogLevel == "debug",
		waitGroup: params.waitGroup,
		service:   api.NewController(apiParams, api.WithLogger(logger)),
		newRelic:  params.newRelic,
		slack:     params.slack,
	}, nil
}

func getSecret(ctx context.Context, p *params) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// データベース認証情報の取得
		if p.config.DBSecretName == "" {
			p.dbHost = p.config.DBHost
			p.dbPort = p.config.DBPort
			p.dbUsername = p.config.DBUsername
			p.dbPassword = p.config.DBPassword
			return nil
		}
		secrets, err := p.secret.Get(ectx, p.config.DBSecretName)
		if err != nil {
			return err
		}
		p.dbHost = secrets["host"]
		p.dbPort = secrets["port"]
		p.dbUsername = secrets["username"]
		p.dbPassword = secrets["password"]
		return nil
	})
	eg.Go(func() error {
		// New Relic認証情報の取得
		if p.config.NewRelicSecretName == "" {
			p.newRelicLicense = p.config.NewRelicLicense
			return nil
		}
		secrets, err := p.secret.Get(ectx, p.config.NewRelicSecretName)
		if err != nil {
			return err
		}
		p.newRelicLicense = secrets["license"]
		return nil
	})
	eg.Go(func() error {
		// Slack認証情報の取得
		if p.config.SlackSecretName == "" {
			p.slackToken = p.config.SlackAPIToken
			p.slackChannelID = p.config.SlackChannelID
			return nil
		}
		secrets, err := p.secret.Get(ectx, p.config.SlackSecretName)
		if err != nil {
			return err
		}
		p.slackToken = secrets["token"]
		p.slackChannelID = secrets["channelId"]
		return nil
	})
	return eg.Wait()
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
