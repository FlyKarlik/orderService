package proto_mapper

func ToInt64Proto(value *int64) int64 {
	if value == nil {
		return 0
	}

	return *value
}

func FromInt64Proto(value int64) *int64 {
	return &value
}

func ToOptionalInt64Proto(value *int64) *int64 {
	return value
}

func FromOptionalInt64Proto(value *int64) *int64 {
	return value
}
