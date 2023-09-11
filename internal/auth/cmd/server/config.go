package server

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	AppName              string `envconfig:"APP_NAME" default:"auth"`
	Environment          string `envconfig:"ENV" default:"none"`
	Port                 int64  `envconfig:"PORT" default:"8080"`
	MetricsPort          int64  `envconfig:"METRICS_PORT" default:"9090"`
	ShutdownDelaySec     int64  `envconfig:"SHUTDOWN_DELAY_SEC" default:"20"`
	LogPath              string `envconfig:"LOG_PATH" default:""`
	LogLevel             string `envconfig:"LOG_LEVEL" default:"info"`
	NewRelicLicense      string `envconfig:"NEW_RELIC_LICENSE" default:""`
	NewRelicSecretName   string `envconfig:"NEW_RELIC_SECRET_NAME" default:""`
	AWSRegion            string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	CognitoAdminPoolID   string `envconfig:"COGNITO_Admin_POOL_ID" default:""`
	CognitoAdminClientID string `envconfig:"COGNITO_Admin_CLIENT_ID" default:""`
	CognitoUserPoolID    string `envconfig:"COGNITO_USER_POOL_ID" default:""`
	CognitoUserClientID  string `envconfig:"COGNITO_USER_CLIENT_ID" default:""`
	SlackAPIToken        string `envconfig:"SLACK_API_TOKEN" default:""`
	SlackChannelID       string `envconfig:"SLACK_CHANNEL_ID" default:""`
	SlackSecretName      string `envconfig:"SLACK_SECRET_NAME" default:""`
}

func newConfig() (*config, error) {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return conf, fmt.Errorf("config: failed to new config: %w", err)
	}
	return conf, nil
}
