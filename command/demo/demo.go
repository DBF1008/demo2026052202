package demo

import (
	"ginskeleton/app/global/variable"
	"github.com/spf13/cobra"
)

var (

	SearchEngines string

	SearchType string

	KeyWords string
)

var Demo1 = &cobra.Command{
	Use:     "sousuo",
	Aliases: []string{"sou", "ss", "s"},
	Short:   "这是一个Demo，以搜索内容进行演示业务逻辑...",
	Long: `调用方法：
			1.进入项目根目录（Ginkeleton）。
			2.执行 go  run  cmd/cli/main.go sousuo -h
			3.执行 go  run  cmd/cli/main.go sousuo 百度
			4.执行 go  run  cmd/cli/main.go  sousuo 百度 -K 关键词  -E  baidu -T img
		`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		variable.ZapLog.Sugar().Infof("Run函数子命令的前置方法，位置参数：%v ，flag参数：%s, %s, %s \n", args[0], SearchEngines, SearchType, KeyWords)
	},

	PreRun: func(cmd *cobra.Command, args []string) {
		variable.ZapLog.Sugar().Infof("Run函数的前置方法，位置参数：%v ，flag参数：%s, %s, %s \n", args[0], SearchEngines, SearchType, KeyWords)

	},

	Run: func(cmd *cobra.Command, args []string) {

		start(SearchEngines, SearchType, KeyWords)
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		variable.ZapLog.Sugar().Infof("Run函数的后置方法，位置参数：%v ，flag参数：%s, %s, %s \n", args[0], SearchEngines, SearchType, KeyWords)
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {

		variable.ZapLog.Sugar().Infof("Run函数子命令的后置方法，位置参数：%v ，flag参数：%s, %s, %s \n", args[0], SearchEngines, SearchType, KeyWords)
	},
}

func init() {
	Demo1.AddCommand(subCmd)
	Demo1.Flags().StringVarP(&SearchEngines, "Engines", "E", "baidu", "-E 或者 --Engines 选择搜索引擎，例如：baidu、sogou")
	Demo1.Flags().StringVarP(&SearchType, "Type", "T", "img", "-T 或者 --Type 选择搜索的内容类型，例如：图片类")
	Demo1.Flags().StringVarP(&KeyWords, "KeyWords", "K", "关键词", "-K 或者 --KeyWords 搜索的关键词")

}

func start(SearchEngines, SearchType, KeyWords string) {

	variable.ZapLog.Sugar().Infof("您输入的搜索引擎：%s， 搜索类型：%s, 关键词：%s\n", SearchEngines, SearchType, KeyWords)

}
