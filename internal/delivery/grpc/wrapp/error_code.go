package wrapp

import (
	"github.com/FlyKarlik/orderService/internal/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetStatusCodeFromError(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	if customErr, ok := err.(*errs.CustomError); ok {
		switch customErr.Code {
		case errs.CodeMarketNotFound:
			return codes.NotFound
		case errs.CodeInvalidUserID:
			return codes.InvalidArgument
		case errs.CodeInvalidOrderID:
			return codes.InvalidArgument
		default:
			return codes.Internal
		}
	}

	if s, ok := status.FromError(err); ok {
		return s.Code()
	}

	return codes.Unknown
}
