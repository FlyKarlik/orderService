package driver

import (
	"context"

	"github.com/FlyKarlik/orderService/config"
	grpc_interceptor "github.com/FlyKarlik/orderService/internal/delivery/grpc/interceptor"
	"github.com/FlyKarlik/orderService/internal/domain"
	spot_instrument_driver "github.com/FlyKarlik/orderService/internal/driver/spot_instrument"
	grpc_client "github.com/FlyKarlik/orderService/pkg/client/grpc"
	"github.com/FlyKarlik/orderService/pkg/logger"
	pb "github.com/FlyKarlik/proto/spot_instrument_service/gen/spot_instrument_service/proto"
)

type IMarketDriver interface {
	ViewMarkets(ctx context.Context, req domain.ViewMarketsRequest) (domain.ViewMarketsResponse, error)
}

type Driver interface {
	IMarketDriver
}

type driverImpl struct {
	IMarketDriver
}

func New(
	cfg *config.Config,
	logger logger.Logger,
	interceptor *grpc_interceptor.GRPCInterceptor,
) (*driverImpl, []grpc_client.IGRPCClient, error) {
	grpcConns := make([]grpc_client.IGRPCClient, 0)
	spotInstrumentConn, err := spot_instrument_driver.SetupSpotInstrumentClient(cfg, interceptor)
	if err != nil {
		return nil, nil, err
	}
	grpcConns = append(grpcConns, spotInstrumentConn)

	return &driverImpl{
		IMarketDriver: spot_instrument_driver.NewMarketDriver(
			logger,
			pb.NewSpotInstrumentServiceClient(spotInstrumentConn),
		),
	}, grpcConns, nil
}
