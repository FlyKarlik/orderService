package proto_mapper

func ToStringProto(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func FromStringProto(value string) *string {
	return &value
}

func ToOptionalStringProto(value *string) *string {
	return value
}

func FromOptionalStringProto(value *string) *string {
	return value
}
