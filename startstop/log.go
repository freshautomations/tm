package startstop

import (
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
	"io"
	"tm/tm/v2/consts"
	"tm/tm/v2/context"
)

func Log(ctx context.Context) {
	follow := viper.GetBool("follow")
	followAndRetry := viper.GetBool("follow-and-retry")
	var location *tail.SeekInfo
	if follow || followAndRetry {
		location = &tail.SeekInfo{
			Offset: -15 * 120, // about 120 char per line, 15 lines
			Whence: io.SeekEnd,
		}
	}
	for _, fullNodename := range ctx.Input {
		home := ctx.Config.GetHome(fullNodename)
		t, _ := tail.TailFile(consts.GetLog(home), tail.Config{
			Location: location,
			ReOpen:   followAndRetry,
			Follow:   follow || followAndRetry,
		})
		for line := range t.Lines {
			fmt.Println(line.Text)
		}
	}
}
