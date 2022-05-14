package startstop

import (
	"tm/tm/v2/context"
	"tm/tm/v2/execute"
	"tm/tm/v2/ux"
)

func Stop(ctx context.Context) {
	for _, fullNodename := range ctx.Input {
		pid := execute.GetPid(ctx.Config.GetHome(fullNodename))
		if pid != nil {
			err := execute.Stop(ctx.Config.GetHome(fullNodename))
			if err != nil {
				ux.Info("✘ %s not stopped: %s.", fullNodename, err)
				continue
			}
		}
		ux.Info("✔ %s stopped.", fullNodename)
	}
}
