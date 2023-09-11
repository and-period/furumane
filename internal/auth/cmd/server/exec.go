package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	apgrpc "github.com/and-period/furumane/pkg/grpc"
	aphttp "github.com/and-period/furumane/pkg/http"
	"github.com/and-period/furumane/pkg/log"
	"github.com/and-period/furumane/proto/auth"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func (a *app) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 環境変数の読み込み
	conf, err := newConfig()
	if err != nil {
		return err
	}

	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(conf.LogLevel), log.WithOutput(conf.LogPath))
	if err != nil {
		return err
	}
	defer logger.Sync() //nolint:errcheck

	// 依存関係の解決
	reg, err := newRegistry(ctx, conf, logger)
	if err != nil {
		logger.Error("Failed to new registry", zap.Error(err))
		return err
	}

	// gRPC Serverの設定
	gopts := apgrpc.NewGRPCOptions(apgrpc.WithLogger(logger))

	s := grpc.NewServer(gopts...)
	auth.RegisterAuthServiceServer(s, reg.service)

	gs, err := apgrpc.NewGRPCServer(s, conf.Port)
	if err != nil {
		logger.Error("Failed to new grpc server", zap.Error(err))
		return err
	}

	// Metrics Serverの設定
	ms := aphttp.NewMetricsServer(conf.MetricsPort)

	// Serverの起動
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if err = ms.Serve(); err != nil {
			logger.Error("Failed to run metrics server", zap.Error(err))
		}
		return
	})
	eg.Go(func() (err error) {
		if err = gs.Serve(); err != nil {
			logger.Error("Failed to run grpc server", zap.Error(err))
		}
		return
	})
	logger.Info("Started server", zap.Int64("port", conf.Port))

	// シグナル検知設定
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ectx.Done():
		logger.Error("Done context", zap.Error(ectx.Err()))
	case signal := <-signalCh:
		logger.Info("Received signal", zap.Any("signal", signal))
		delay := time.Duration(conf.ShutdownDelaySec) * time.Second
		logger.Info("Pre-shutdown", zap.Duration("delay", delay))
		time.Sleep(delay)
	}

	// Serverの停止
	logger.Info("Shutdown...")
	gs.Stop()
	if err = ms.Stop(ectx); err != nil {
		logger.Error("Failed to stopeed metrics server", zap.Error(err))
		return err
	}
	reg.waitGroup.Wait()
	return eg.Wait()
}
