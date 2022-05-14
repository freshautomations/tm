package startstop

import (
	"tm/m/v2/context"
	"tm/m/v2/execute"
	"tm/m/v2/ux"
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
