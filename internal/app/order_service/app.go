package order_service

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/FlyKarlik/orderService/config"
	grpc_async_handler "github.com/FlyKarlik/orderService/internal/delivery/grpc/async"
	grpc_interceptor "github.com/FlyKarlik/orderService/internal/delivery/grpc/interceptor"
	grpc_sync_handler "github.com/FlyKarlik/orderService/internal/delivery/grpc/sync"
	"github.com/FlyKarlik/orderService/internal/driver"
	"github.com/FlyKarlik/orderService/internal/repository"
	"github.com/FlyKarlik/orderService/internal/usecase"
	grpc_client "github.com/FlyKarlik/orderService/pkg/client/grpc"
	"github.com/FlyKarlik/orderService/pkg/logger"
	"github.com/FlyKarlik/orderService/pkg/metric"
	"github.com/FlyKarlik/orderService/pkg/tracer"
	pb "github.com/FlyKarlik/proto/order_service/gen/order_service/proto"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type OrderService struct {
	cfg        *config.Config
	grpcServer *grpc.Server
	logger     logger.Logger
}

func New(cfg *config.Config, logger logger.Logger) *OrderService {
	return &OrderService{
		cfg:    cfg,
		logger: logger,
	}
}

func (o *OrderService) Start() error {
	const layer = "app"
	const method = "Start"

	o.logger.Info(layer, method, "starting service")

	o.mustSetupTracer()

	driver, clients, err := o.mustSetupDriver(o.cfg, o.logger)
	if err != nil {
		o.logger.Error(layer, method, "failed to init driver client", err)
		return err
	}

	repo := o.mustSetupRepo()
	usecase := o.mustSetupUsecase(driver, repo)

	go func() {
		o.logger.Infof(
			layer,
			method,
			"starting prometheus",
			"address: %s", o.cfg.Infrastructure.Prometheus.Address)
		if err := o.mustStartPrometheus(); err != nil {
			o.logger.Error(layer, method, "failed to start prometheus", err)
			os.Exit(1)
		}
		o.logger.Info(layer, method, "prometheus started successfully")
	}()

	go func() {
		o.logger.Infof(
			layer,
			method,
			"starting gRPC server",
			"address: %s", o.cfg.GRPCServer.Address)
		if err := o.mustStartGRPCServer(usecase); err != nil {
			o.logger.Error(layer, method, "failed to start grpc server", err)
			os.Exit(1)
		}
		o.logger.Info(layer, method, "grpc server stopped")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	o.logger.Info(layer, method, "waiting for shutdown signal")
	<-quit
	o.logger.Info(layer, method, "shutdown signal received")

	err = o.mustCloseConnectionWithGRPCClients(clients)
	if err != nil {
		o.logger.Error(layer, method, "failed to close connection with grpc clients", err)
	}
	o.mustStopGRPCServer()
	o.logger.Info(layer, method, "service stopped gracefully")

	return nil
}

func (o *OrderService) mustStartPrometheus() error {
	return metric.StartPrometheus(o.cfg)
}

func (o *OrderService) mustSetupTracer() {
	const method = "mustSetupTracer"
	const layer = "app"

	tracer.New(context.Background(), o.cfg)

	o.logger.Info(layer, method, "setting up tracing")
}

func (o *OrderService) mustSetupRepo() repository.Repository {
	const method = "mustSetupRepo"
	const layer = "app"

	o.logger.Info(layer, method, "setting up repository")
	return repository.New(o.logger)
}

func (o *OrderService) mustSetupDriver(cfg *config.Config, l logger.Logger) (driver.Driver, []grpc_client.IGRPCClient, error) {
	const method = "mustSetuDriver"
	const layer = "app"

	o.logger.Info(layer, method, "setting up driver")
	return driver.New(cfg, l)
}

func (o *OrderService) mustSetupUsecase(driver driver.Driver, repo repository.Repository) usecase.Usecase {
	const method = "mustSetuUsecase"
	const layer = "app"

	o.logger.Info(layer, method, "setting up usecase")
	return usecase.New(o.logger, driver, repo)
}

func (o *OrderService) mustStartGRPCServer(usecase usecase.Usecase) error {
	const layer = "app"
	const method = "mustStartGRPCServer"

	grpcInterceptor := grpc_interceptor.New(o.logger)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcInterceptor.XRequestIDInterceptor(),
			grpcInterceptor.LoggerInterceptor(),
			grpcInterceptor.UnaryPanicRecoveryInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
		),
	)
	o.grpcServer = grpcServer

	grpcSyncHandler := grpc_sync_handler.New(o.logger, usecase)
	grpcAsyncHandler := grpc_async_handler.New(o.logger, usecase)

	pb.RegisterOrderSyncServiceServer(grpcServer, grpcSyncHandler)
	pb.RegisterOrderStreamServiceServer(grpcServer, grpcAsyncHandler)

	grpc_prometheus.Register(grpcServer)

	lis, err := net.Listen("tcp", o.cfg.GRPCServer.Address)
	if err != nil {
		o.logger.Error(layer, method, "failed to listen tcp", err, "address", o.cfg.GRPCServer.Address)
		return err
	}

	o.logger.Info(layer, method, "grpc server listening", "address", o.cfg.GRPCServer.Address)
	err = grpcServer.Serve(lis)
	if err != nil {
		o.logger.Error(layer, method, "grpc server serve error", err)
		return err
	}

	return nil
}

func (o *OrderService) mustStopGRPCServer() {
	const layer = "app"
	const method = "mustStopGRPCServer"

	o.logger.Info(layer, method, "stopping grpc server")
	o.grpcServer.GracefulStop()
	o.logger.Info(layer, method, "grpc server stopped gracefully")
}

func (o *OrderService) mustCloseConnectionWithGRPCClients(clients []grpc_client.IGRPCClient) error {
	const layer = "app"
	const method = "mustCloseConnectionWithGClients"

	o.logger.Info(layer, method, "closing grpc client connections")

	for _, client := range clients {
		if err := client.Close(); err != nil {
			o.logger.Error(layer, method, "failed to close grpc client", err)
		}
	}

	o.logger.Info(layer, method, "all grpc clients closed")
	return nil
}
