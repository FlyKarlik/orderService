package spot_instrument_driver

import (
	"github.com/FlyKarlik/orderService/config"
	grpc_interceptor "github.com/FlyKarlik/orderService/internal/delivery/grpc/interceptor"
	grpc_client "github.com/FlyKarlik/orderService/pkg/client/grpc"
	"google.golang.org/grpc"
)

func SetupSpotInstrumentClient(
	cfg *config.Config,
	interceptor *grpc_interceptor.GRPCInterceptor,
) (grpc_client.IGRPCClient, error) {
	conn, err := grpc_client.New(
		cfg.GRPCApi.SpotInstrumentServiceHost,
		cfg,
		grpc.WithUnaryInterceptor(interceptor.XRequestIDUnaryClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
