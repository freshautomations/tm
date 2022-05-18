package startstop

import (
	"tm/tm/v2/context"
	"tm/tm/v2/execute"
	"tm/tm/v2/ux"
)

func Reset(ctx context.Context) {
	for _, fullNodename := range ctx.Input {
		pid := execute.GetPid(ctx.Config.GetHome(fullNodename))
		if pid != nil {
			err := execute.Stop(ctx.Config.GetHome(fullNodename))
			if err != nil {
				ux.Info("✘ %s not stopped: %s.", fullNodename, err)
				continue
			}
		}
		execute.Reset(ctx.Config.GetBinary(fullNodename), ctx.Config.GetHome(fullNodename))
		if pid != nil {
			_, err := execute.Start(ctx.Config.GetBinary(fullNodename), ctx.Config.GetHome(fullNodename))
			if err != nil {
				ux.Info("✘ %s not started, %s.", fullNodename, err)
				continue
			}
		}
		ux.Info("✔ %s reset.", fullNodename)
	}
}
