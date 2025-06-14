package mapper

import (
	"github.com/FlyKarlik/orderService/internal/domain"
	"github.com/FlyKarlik/orderService/pkg/proto_mapper"
	spotPb "github.com/FlyKarlik/proto/spot_instrument_service/gen/spot_instrument_service/proto"
)

func ToProtoViewMarketsRequest(req domain.ViewMarketsRequest) *spotPb.ViewMarketsRequest {
	return &spotPb.ViewMarketsRequest{
		UserRoles: ToProtoUserRoles(req.UserRoles, spotPb.UserRole(0)),
	}
}

func FromProtoViewMarketsResponse(pb *spotPb.ViewMarketsResponse) domain.ViewMarketsResponse {
	return domain.ViewMarketsResponse{
		Markets: FromProtoMarkets(pb.Markets),
	}
}

func FromProtoMarket(pb *spotPb.Market) domain.Market {
	return domain.Market{
		ID:           proto_mapper.FromIDProto(&pb.Id),
		Name:         proto_mapper.FromStringProto(pb.Name),
		Enabled:      proto_mapper.FromBoolProto(pb.Enabled),
		DeletedAt:    proto_mapper.FromTimestampProto(pb.DeletedAt),
		AllowedRoles: FromProtoUserRoles(pb.AllowedRoles),
	}
}

func FromProtoMarkets(pb []*spotPb.Market) []domain.Market {
	markets := make([]domain.Market, 0, len(pb))
	for _, market := range pb {
		markets = append(markets, FromProtoMarket(market))
	}
	return markets
}
