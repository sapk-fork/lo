package lo_test

import (
	"context"
	"fmt"

	"github.com/samber/lo"
)

// ExampleContext set and retrieve a custom user type from context
func ExampleContext() {
	ctx := context.Background()

	type user struct {
		id, name string
	}

	ctx = lo.ContextWith(ctx, &user{id: "42", name: "John Doe"})

	userfromContext, ok := lo.FromContext[*user](ctx)
	fmt.Printf("%v\n%#v", ok, userfromContext)

	// Output:
	// true
	// &lo_test.user{id:"42", name:"John Doe"}
}
