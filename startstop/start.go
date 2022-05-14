package startstop

import (
	"tm/m/v2/context"
	"tm/m/v2/execute"
	"tm/m/v2/initialize"
	"tm/m/v2/ux"
)

func Start(ctx context.Context) {
	for _, fullNodename := range ctx.Input {
		pid := execute.GetPid(ctx.Config.GetHome(fullNodename))
		if pid != nil {
			ux.Info("⚠ %s skipped, PID %d.", fullNodename, *pid)
			continue
		}
		initialize.ValidateGenesis(ctx, fullNodename)
		pidInt, err := execute.Start(ctx.Config.GetBinary(fullNodename), ctx.Config.GetHome(fullNodename))
		if err != nil {
			ux.Info("✘ %s not started, %s.", fullNodename, err)
			continue
		}
		ux.Info("✔ %s started, PID %d.", fullNodename, pidInt)
	}
}
