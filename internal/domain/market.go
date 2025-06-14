package domain

import (
	"time"

	"github.com/google/uuid"
)

type Market struct {
	ID           *uuid.UUID
	Name         *string
	Enabled      *bool
	DeletedAt    *time.Time
	AllowedRoles UserRolesEnum
}

type ViewMarketsRequest struct {
	UserRoles UserRolesEnum
}

type ViewMarketsResponse struct {
	Markets []Market
}
