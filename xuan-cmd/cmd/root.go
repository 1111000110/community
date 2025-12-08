package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xuan",                                             // 跟命令名称
	Short: "Project CLI",                                      // 简短描述
	Long:  "Xuan CLI - generate code, run and stop services.", // 长描述
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { // 在所有子命令执行前执行
		if !isProjectRoot() {
			return fmt.Errorf("xuan 命令只能在项目根目录执行")
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isProjectRoot() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}
