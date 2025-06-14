package mapper

import (
	"github.com/FlyKarlik/orderService/internal/domain"
	"github.com/FlyKarlik/orderService/pkg/proto_mapper"
	pb "github.com/FlyKarlik/proto/order_service/gen/order_service/proto"
)

func FromProtoCreateOrderRequest(pb *pb.CreateOrderRequest) domain.CreateOrderRequest {
	return domain.CreateOrderRequest{
		UserID:    proto_mapper.FromIDProto(&pb.UserId),
		MarketID:  proto_mapper.FromIDProto(&pb.MarketId),
		OrderType: (*domain.OrderTypeEnum)(proto_mapper.FromStringProto(pb.OrderType.String())),
		Price:     proto_mapper.FromStringProto(pb.Price),
		Quantity:  proto_mapper.FromInt64Proto(pb.Quantity),
		UserRoles: FromProtoUserRoles(pb.UserRoles),
	}
}

func FromProtoGetOrderStatusRequest(pb *pb.GetOrderStatusRequest) domain.GetOrderStatusRequest {
	return domain.GetOrderStatusRequest{
		OrderID: proto_mapper.FromIDProto(&pb.OrderId),
		UserID:  proto_mapper.FromIDProto(&pb.UserId),
	}
}

func FromProtoStreamOrderUpdatesRequest(pb *pb.StreamOrderUpdatesRequest) domain.StreamOrderUpdatesRequest {
	return domain.StreamOrderUpdatesRequest{
		OrderID: proto_mapper.FromIDProto(&pb.OrderId),
		UserID:  proto_mapper.FromIDProto(&pb.UserId),
	}
}

func ToProtoStreamOrderUpdatesResponse(domain domain.StreamOrderUpdatesResponse) *pb.OrderUpdate {
	return &pb.OrderUpdate{
		OrderId:   proto_mapper.ToIDProto(domain.OrderID),
		Status:    MapEnumToOrderStatus(domain.OrderStatus),
		UpdatedAt: proto_mapper.ToTimestampProto(domain.UpdatedAt),
	}
}

func ToProtoCreateOrderResponse(domain domain.CreateOrderResponse) *pb.CreateOrderResponse {
	return &pb.CreateOrderResponse{
		OrderId: proto_mapper.ToIDProto(domain.OrderID),
		Status:  MapEnumToOrderStatus(domain.OrderStatus),
	}
}

func ToProtoGetOrderStatusResponse(domain domain.GetOrderStatusResponse) *pb.GetOrderStatusResponse {
	return &pb.GetOrderStatusResponse{
		Status: MapEnumToOrderStatus(domain.Status),
	}
}

func MapOrderStatusToEnum(status *pb.OrderStatus) domain.OrderStatusEnum {
	if status == nil {
		return domain.OrderStatusEnumUnspecified
	}

	switch *status {
	case pb.OrderStatus_CREATED:
		return domain.OrderStatusEnumCreated
	case pb.OrderStatus_PENDING:
		return domain.OrderStatusEnumPending
	case pb.OrderStatus_FILLED:
		return domain.OrderStatusEnumFilled
	case pb.OrderStatus_REJECTED:
		return domain.OrderStatusEnumRejected
	default:
		return domain.OrderStatusEnumUnspecified
	}
}

func MapEnumToOrderStatus(enum *domain.OrderStatusEnum) pb.OrderStatus {
	if enum == nil {
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}

	switch *enum {
	case domain.OrderStatusEnumCreated:
		return pb.OrderStatus_CREATED
	case domain.OrderStatusEnumPending:
		return pb.OrderStatus_PENDING
	case domain.OrderStatusEnumFilled:
		return pb.OrderStatus_FILLED
	case domain.OrderStatusEnumRejected:
		return pb.OrderStatus_REJECTED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}
