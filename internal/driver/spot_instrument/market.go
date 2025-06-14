package spot_instrument_driver

import (
	"context"

	"github.com/FlyKarlik/orderService/internal/domain"
	"github.com/FlyKarlik/orderService/internal/mapper"
	"github.com/FlyKarlik/orderService/pkg/logger"
	pb "github.com/FlyKarlik/proto/spot_instrument_service/gen/spot_instrument_service/proto"
)

type marketDriver struct {
	logger logger.Logger
	client pb.SpotInstrumentServiceClient
}

func NewMarketDriver(l logger.Logger, client pb.SpotInstrumentServiceClient) *marketDriver {
	return &marketDriver{
		logger: l,
		client: client,
	}
}

func (s *marketDriver) ViewMarkets(
	ctx context.Context,
	req domain.ViewMarketsRequest,
) (domain.ViewMarketsResponse, error) {
	resp, err := s.client.ViewMarkets(ctx, mapper.ToProtoViewMarketsRequest(req))
	if err != nil {
		return domain.ViewMarketsResponse{}, err
	}
	return mapper.FromProtoViewMarketsResponse(resp), nil
}
