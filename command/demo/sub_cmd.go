package demo

import (
	"fmt"
	"github.com/spf13/cobra"
)

var subCmd = &cobra.Command{
	Use:   "subCmd",
	Short: "subCmd 命令简要介绍",
	Long:  `命令使用详细介绍`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", args[0])
	},
}

func init() {
	Demo1.AddCommand(subCmd)

}
