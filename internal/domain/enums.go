package domain

type OrderTypeEnum string

const (
	OrderTypeEnumUnspecified OrderTypeEnum = "UNSPECIFIED"
	OrderTypeEnumLimit       OrderTypeEnum = "LIMIT"
	OrderTypeEnumMarket      OrderTypeEnum = "MARKET"
)

func (o OrderTypeEnum) String() string {
	return string(o)
}

type OrderStatusEnum string

const (
	OrderStatusEnumUnspecified OrderStatusEnum = "UNSPECIFIED"
	OrderStatusEnumCreated     OrderStatusEnum = "CREATED"
	OrderStatusEnumPending     OrderStatusEnum = "PENDING"
	OrderStatusEnumFilled      OrderStatusEnum = "FILLED"
	OrderStatusEnumRejected    OrderStatusEnum = "REJECTED"
)

func (o OrderStatusEnum) String() string {
	return string(o)
}

type UserRoleEnum string

const (
	UserRoleEnumUnspecified UserRoleEnum = "UNSPECIFIED"
	UserRoleEnumTrader      UserRoleEnum = "TRADER"
	UserRoleEnumViewer      UserRoleEnum = "VIEWER"
	UserRoleEnumAdmin       UserRoleEnum = "ADMIN"
)

func (u UserRoleEnum) String() string {
	return string(u)
}

type UserRolesEnum []UserRoleEnum

func (u UserRolesEnum) Strings() []string {
	res := make([]string, 0, len(u))
	for _, role := range u {
		res = append(res, role.String())
	}
	return res
}
