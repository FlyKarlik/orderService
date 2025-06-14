package spot_instrument_driver

import (
	"github.com/FlyKarlik/orderService/config"
	grpc_client "github.com/FlyKarlik/orderService/pkg/client/grpc"
)

func SetupSpotInstrumentClient(cfg *config.Config) (grpc_client.IGRPCClient, error) {
	conn, err := grpc_client.New(cfg.GRPCApi.SpotInstrumentServiceHost, cfg)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
