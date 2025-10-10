package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xuan",
	Short: "Project CLI",
	Long:  "Xuan CLI - generate code, run and stop services.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
