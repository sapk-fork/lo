package lo_test

import (
	"bytes"
	"context"
	"log"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestContextString(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	ctx := lo.ContextWith(context.Background(), "some string")
	is.NotNil(ctx)

	result, ok := lo.FromContext[string](ctx)

	is.True(ok)
	is.Equal("some string", result)
}

func TestContextFunc(t *testing.T) { // not comparable
	t.Parallel()

	is := assert.New(t)

	ctx := lo.ContextWith(context.Background(), func() string { return "ok" })
	is.NotNil(ctx)

	result, ok := lo.FromContext[func() string](ctx)

	is.True(ok)
	is.Equal("ok", result())
}

func TestContextCustom(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	type user struct {
		id, name string
	}

	expected := user{
		id:   lo.RandomString(10, lo.AlphanumericCharset),
		name: lo.RandomString(10, lo.AlphanumericCharset),
	}

	ctx := lo.ContextWith(context.Background(), expected)
	is.NotNil(ctx)

	result, ok := lo.FromContext[user](ctx)

	is.True(ok)
	is.Equal(expected.id, result.id)
	is.Equal(expected.name, result.name)

	// should not be updated
	expected.name = lo.RandomString(10, lo.AlphanumericCharset)
	is.NotEqual(expected.name, result.name)
}

func TestContextCustomPointer(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	type user struct {
		id, name string
	}

	expected := user{
		id:   lo.RandomString(10, lo.AlphanumericCharset),
		name: lo.RandomString(10, lo.AlphanumericCharset),
	}

	ctx := lo.ContextWith(context.Background(), &expected)
	is.NotNil(ctx)

	result, ok := lo.FromContext[*user](ctx)

	is.True(ok)
	is.Equal(expected.id, result.id)
	is.Equal(expected.name, result.name)

	// should be updated by reference
	expected.name = lo.RandomString(10, lo.AlphanumericCharset)
	is.Equal(expected.name, result.name)
}

func TestContextLogger(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	buf := new(bytes.Buffer)

	logger := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}

			return a
		},
	}))

	ctx := lo.ContextWith(context.Background(), logger.With(
		slog.Group("request",
			slog.String("method", http.MethodPost),
			slog.String("url", "http://localhost")),
	))
	is.NotNil(ctx)

	result, ok := lo.FromContext[*slog.Logger](ctx)

	is.True(ok)
	is.NotNil(result)

	result.Info("testing")

	is.Equal(
		"level=INFO msg=testing request.method=POST request.url=http://localhost\n",
		buf.String(),
	)
}

func TestFromContextOrLogger(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	buf := new(bytes.Buffer)

	logger := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}

			return a
		},
	}))

	ctx := lo.ContextWith(context.Background(), logger.With(
		slog.Group("request",
			slog.String("method", http.MethodPost),
			slog.String("url", "http://localhost")),
	))
	is.NotNil(ctx)

	result := lo.FromContextOr[*slog.Logger](ctx, slog.Default())

	is.NotNil(result)

	result.Info("testing")

	is.Equal(
		"level=INFO msg=testing request.method=POST request.url=http://localhost\n",
		buf.String(),
	)
}

func TestFromContextOrLoggerDefault(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	logger := lo.FromContextOr[*log.Logger](context.Background(), log.Default())

	is.NotNil(logger)

	logger.Print("testing")
}

func TestContextMultipleType(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	ctx := context.Background()

	ctx = lo.ContextWith(ctx, "some string")
	is.NotNil(ctx)

	ctx = lo.ContextWith(ctx, time.Date(2023, time.September, 20, 1, 10, 20, 30, time.UTC))
	is.NotNil(ctx)

	ctx = lo.ContextWith(ctx, func() string { return "ok" })
	is.NotNil(ctx)

	type user struct {
		id, name string
	}

	loggedUser := &user{
		id:   lo.RandomString(10, lo.AlphanumericCharset),
		name: lo.RandomString(10, lo.AlphanumericCharset),
	}
	ctx = lo.ContextWith(ctx, loggedUser)
	is.NotNil(ctx)

	resultStr, ok := lo.FromContext[string](ctx)

	is.True(ok)
	is.Equal("some string", resultStr)

	resultTime, ok := lo.FromContext[time.Time](ctx)

	is.True(ok)
	is.Equal(time.Date(2023, time.September, 20, 1, 10, 20, 30, time.UTC), resultTime)

	resultUser, ok := lo.FromContext[*user](ctx)

	is.True(ok)
	is.Equal(loggedUser, resultUser)

	resultFunc, ok := lo.FromContext[func() string](ctx)

	is.True(ok)
	is.Equal("ok", resultFunc())
}
