package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        *uuid.UUID
	UserID    *uuid.UUID
	MarketID  *uuid.UUID
	OrderType *OrderTypeEnum
	Price     *string
	Quantity  *int64
	Status    *OrderStatusEnum
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type CreateOrderRequest struct {
	UserID    *uuid.UUID     `validate:"required"`
	MarketID  *uuid.UUID     `validate:"required"`
	OrderType *OrderTypeEnum `validate:"required"`
	Price     *string        `validate:"required"`
	Quantity  *int64         `validate:"required,gt=0"`
	UserRoles UserRolesEnum  `validate:"required,gt=0"`
}

type CreateOrderResponse struct {
	OrderID     *uuid.UUID
	OrderStatus *OrderStatusEnum
}

type GetOrderStatusRequest struct {
	OrderID *uuid.UUID `validate:"required"`
	UserID  *uuid.UUID `validate:"required"`
}

type GetOrderStatusResponse struct {
	Status *OrderStatusEnum
}

type StreamOrderUpdatesRequest struct {
	OrderID *uuid.UUID `validate:"required"`
	UserID  *uuid.UUID `validate:"required"`
}

type StreamOrderUpdatesResponse struct {
	OrderID     *uuid.UUID
	OrderStatus *OrderStatusEnum
	UpdatedAt   *time.Time
}
