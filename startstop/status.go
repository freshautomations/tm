package startstop

import (
	"strconv"
	"tm/m/v2/context"
	"tm/m/v2/execute"
	"tm/m/v2/ux"
)

func Status(ctx context.Context) {
	for _, fullNodename := range ctx.Input {
		pid := execute.GetPid(ctx.Config.GetHome(fullNodename))
		if pid != nil {
			ux.Info("✔ %s running, PID %s.", fullNodename, strconv.Itoa(*pid))
		} else {
			ux.Info("✘ %s stopped.", fullNodename)
		}
	}
}
