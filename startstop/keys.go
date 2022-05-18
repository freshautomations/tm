package startstop

import (
	"fmt"
	"tm/tm/v2/context"
)

func Keys(ctx context.Context) {
	for _, fullNodename := range ctx.Input {
		home := ctx.Config.GetHome(fullNodename)
		fmt.Println("Get keys from %s", home)
	}
}
