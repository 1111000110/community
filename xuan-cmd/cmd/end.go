package cmd

import (
	"community/xuan-cmd/commands/end"
	"github.com/spf13/cobra"
)

var endCmd = &cobra.Command{
	Use:   "end",
	Short: "Stop services and clean resources",
}

var endAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Stop all services (optionally clean data)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cleanDataFlag {
			return end.CleanAllDataWrapper()
		}
		return end.StopAllServicesWrapper()
	},
}

var (
	cleanDataFlag bool
)

var endCleanCmd = &cobra.Command{
	Use:   "clean [logs|data|binaries|all]",
	Short: "Clean logs, data, binaries, or all",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return end.CleanAllDataWrapper()
		}
		switch args[0] {
		case "logs":
			return end.CleanLogsWrapper()
		case "data":
			return end.CleanDataDirectoriesWrapper()
		case "binaries":
			return end.CleanBinariesWrapper()
		case "all":
			return end.CleanAllDataWrapper()
		default:
			return cmd.Help()
		}
	},
}

var endDockerCmd = &cobra.Command{
	Use:   "docker [service]",
	Short: "Stop docker-compose or a specific service",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return end.RunCmdStreamingWrapper("docker-compose", "down")
		}
		return end.RunCmdStreamingWrapper("docker-compose", "stop", args[0])
	},
}

func init() {
	rootCmd.AddCommand(endCmd)
	endCmd.AddCommand(endAllCmd)
	endCmd.AddCommand(endCleanCmd)
	endCmd.AddCommand(endDockerCmd)
	endAllCmd.Flags().BoolVarP(&cleanDataFlag, "clean-data", "c", false, "Clean all data after stopping")
}
