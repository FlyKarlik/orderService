package proto_mapper

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromTimestampProto(value *timestamppb.Timestamp) *time.Time {
	if value == nil || !value.IsValid() {
		return nil
	}

	_t := value.AsTime()

	return &_t
}

func ToTimestampProto(value *time.Time) *timestamppb.Timestamp {
	if value == nil || value.IsZero() {
		return nil
	}

	return timestamppb.New(*value)
}
