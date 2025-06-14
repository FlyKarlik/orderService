package shared_context

type ContextKeyEnum string

const (
	ContextKeyEnumXRequestID ContextKeyEnum = "X_REQUEST_ID"
)

func (c ContextKeyEnum) String() string {
	return string(c)
}
