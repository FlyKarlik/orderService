package proto_mapper

func ToBoolProto(value *bool) bool {
	if value == nil {
		return false
	}

	return *value
}

func FromBoolProto(value bool) *bool {
	return &value
}

func ToOptionalBoolProto(value *bool) *bool {
	return value
}

func FromOptionalBoolProto(value *bool) *bool {
	return value
}
