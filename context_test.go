package lo_test

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestContextString(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	ctx := lo.ContextWith(context.Background(), "some string")
	is.NotNil(ctx)

	original, ok := lo.FromContext[string](ctx)

	is.True(ok)
	is.Equal("some string", original)
}
