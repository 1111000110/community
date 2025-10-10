package cmd

import (
	"community/xuan/commands/run"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run or restart services",
}

var runAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Start all infra and services",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run.RunAll()
	},
}

var runDockerCmd = &cobra.Command{
	Use:   "docker [service]",
	Short: "Restart docker-compose or one service",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return run.DockerComposeDownUp()
		}
		return run.DockerComposeRestartService(args[0])
	},
}

var runApiCmd = &cobra.Command{
	Use:   "api [path]",
	Short: "Restart API under path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return run.RestartServices(args[0], true, false)
	},
}

var runRpcCmd = &cobra.Command{
	Use:   "rpc [path]",
	Short: "Restart RPC under path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return run.RestartServices(args[0], false, true)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.AddCommand(runAllCmd)
	runCmd.AddCommand(runDockerCmd)
	runCmd.AddCommand(runApiCmd)
	runCmd.AddCommand(runRpcCmd)
}
