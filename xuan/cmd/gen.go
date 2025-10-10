package cmd

import (
	"fmt"

	"community/xuan/commands/gen"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code (api/rpc)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use: xuan gen [api|rpc] <dir> [flags]")
		cmd.Help()
	},
}

var genAPICmd = &cobra.Command{
	Use:   "api [dir]",
	Short: "Generate API code from .api files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 { dir = args[0] }
		return gen.API(dir)
	},
}

var (
	multiple bool
)

var genRPCCmd = &cobra.Command{
	Use:   "rpc [dir]",
	Short: "Generate RPC code from .proto files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 { dir = args[0] }
		return gen.RPC(dir, multiple)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(genAPICmd)
	genRPCCmd.Flags().BoolVarP(&multiple, "multiple", "m", false, "Enable multi-service mode")
	genCmd.AddCommand(genRPCCmd)
}
