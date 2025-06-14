package proto_mapper

import (
	"github.com/google/uuid"
)

func ValidateID(s string) bool {
	err := uuid.Validate(s)
	return err == nil
}

func ValidateIDs(s []string) bool {
	valid := true
	for _, v := range s {
		err := uuid.Validate(v)
		if err != nil {
			valid = false
			break
		}
	}
	return valid
}

func FromIDProto(s *string) *uuid.UUID {
	if s == nil || *s == "" {
		return nil
	}

	parsed, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}

	return &parsed
}

func FromIDsProto(s []string) []uuid.UUID {
	result := make([]uuid.UUID, len(s))
	for i, v := range s {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return nil
		}
		result[i] = parsed
	}

	return result
}

func ToIDProto(u *uuid.UUID) string {
	if u == nil {
		return ""
	}

	return u.String()
}

func ToOptionalIDProto(u *uuid.UUID) *string {
	if u == nil {
		return nil
	}

	res := u.String()

	return &res
}

func ToIDsProto(u []uuid.UUID) []string {
	result := make([]string, len(u))
	for i, v := range u {
		result[i] = ToIDProto(&v)
	}
	return result
}
