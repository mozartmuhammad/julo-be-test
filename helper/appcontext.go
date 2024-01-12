package helper

import (
	"context"
)

type key string

const (
	CustomerXID key = "customer-xid"
)

func SetCustomerXID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, CustomerXID, value)
}

func GetCustomerXID(ctx context.Context) string {
	if v, ok := ctx.Value(CustomerXID).(string); ok {
		return v
	}
	return ""
}
