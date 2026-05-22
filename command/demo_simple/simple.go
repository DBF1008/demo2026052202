package demo_simple

import (
	"ginskeleton/app/global/variable"
	"github.com/spf13/cobra"
	"time"
)

var (
	LogAction string
	Date      string
)

var DemoSimple = &cobra.Command{
	Use:     "demo_simple",
	Aliases: []string{"demo_simple"},
	Short:   "这是一个最简单的demo示例",
	Long: `调用方法：
			1.进入项目根目录（Ginkeleton）。
			2.执行 go  run  cmd/cli/main.go  demo_simple -h
			3.执行 go  run  cmd/cli/main.go  demo_simple  -A insert
		`,

	Run: func(cmd *cobra.Command, args []string) {

		start(LogAction, Date)
	},
}

func init() {
	DemoSimple.Flags().StringVarP(&LogAction, "logAction", "A", "insert", "-A 指定参数动作,例如：-A insert ")
	DemoSimple.Flags().StringVarP(&Date, "date", "D", time.Now().Format("2006-01-02"), "-D 指定日期,例如：-D  2021-09-13")
}

func start(actionName, Date string) {
	switch actionName {
	case "insert":
		variable.ZapLog.Info("insert 参数执行对应业务逻辑,Date参数值：" + Date)
	case "update":
		variable.ZapLog.Info("update 参数执行对应业务逻辑,Date参数值：" + Date)
	}

}
