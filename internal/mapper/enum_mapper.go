package mapper

import "github.com/FlyKarlik/orderService/internal/domain"

type ProtoUserRole interface {
	~int32 | ~int
}

func FromProtoUserRole[E ProtoUserRole](protoRole E) domain.UserRoleEnum {
	switch int32(protoRole) {
	case int32(1):
		return domain.UserRoleEnumTrader
	case int32(2):
		return domain.UserRoleEnumViewer
	case int32(3):
		return domain.UserRoleEnumAdmin
	default:
		return domain.UserRoleEnumUnspecified
	}
}

func ToProtoUserRole[E ProtoUserRole](userRole domain.UserRoleEnum, _ E) E {
	switch userRole {
	case domain.UserRoleEnumTrader:
		return E(1)
	case domain.UserRoleEnumViewer:
		return E(2)
	case domain.UserRoleEnumAdmin:
		return E(3)
	default:
		return E(0)
	}
}

func FromProtoUserRoles[E ProtoUserRole](protoRoles []E) domain.UserRolesEnum {
	userRoles := make(domain.UserRolesEnum, 0, len(protoRoles))
	for _, protoRole := range protoRoles {
		userRoles = append(userRoles, FromProtoUserRole(protoRole))
	}
	return userRoles
}

func ToProtoUserRoles[E ProtoUserRole](userRoles domain.UserRolesEnum, _ E) []E {
	protoRoles := make([]E, 0, len(userRoles))
	for _, role := range userRoles {
		protoRoles = append(protoRoles, ToProtoUserRole(role, E(0)))
	}
	return protoRoles
}
