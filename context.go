package lo

import (
	"context"
)

// ContextWith return a new context with the value attached by type.
func ContextWith[T comparable](ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, *new(T), val)
}

// FromContext returns the entry in context using type as the context key.
func FromContext[T comparable](ctx context.Context) (val T, ok bool) {
	entry := ctx.Value(*new(T))
	val, ok = entry.(T)

	return val, ok
}
