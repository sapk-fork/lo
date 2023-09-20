package lo

import (
	"context"
)

// ContextWith return a new context with the value attached by type.
func ContextWith[T any](ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, (*T)(nil), val)
}

// FromContext returns the entry in context using type as the context key.
func FromContext[T any](ctx context.Context) (val T, ok bool) {
	val, ok = ctx.Value((*T)(nil)).(T)

	return val, ok
}

// FromContextOr returns the entry in context using type as the context key otherwise return default value.
func FromContextOr[T any](ctx context.Context, def T) T {
	val, ok := FromContext[T](ctx)

	if ok {
		return val
	}

	return def
}
