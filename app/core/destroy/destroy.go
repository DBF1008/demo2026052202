package destroy

import (
	"ginskeleton/app/core/event_manage"
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func init() {

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		received := <-c
		variable.ZapLog.Warn(consts.ProcessKilled, zap.String("信号值", received.String()))
		(event_manage.CreateEventManageFactory()).FuzzyCall(variable.EventDestroyPrefix)
		close(c)
		os.Exit(1)
	}()

}
