package cmd

import (
	"fmt"

	"community/xuan-cmd/commands/gen"

	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code (api/rpc)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use: xuan gen [api|rpc] <dir> [flags]")
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

var genAPICmd = &cobra.Command{
	Use:   "api [dir]",
	Short: "Generate API code from .api files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}
		return gen.API(dir)
	},
}

var (
	multiple  bool
	modelType string
)

var genRPCCmd = &cobra.Command{
	Use:   "rpc [dir]",
	Short: "Generate RPC code from .proto files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}
		return gen.RPC(dir, multiple)
	},
}

var genModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate model code",
	Long: `Generate model code based on the specified type.

Supported types:
  mongo    Generate MongoDB model code
  mysql    Generate MySQL model code (coming soon)

Examples:
  xuan gen model mongo ./service/ai
  xuan gen model mongo ./service/ai user`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var genModelMongoCmd = &cobra.Command{
	Use:   "mongo [dir] [model-name]",
	Short: "Generate MongoDB model code",
	Long: `Generate MongoDB model code in the specified service directory.

Arguments:
  dir          Service directory (e.g., ./service/ai)
  model-name   Model name (optional, defaults to service name)

Examples:
  xuan gen model mongo ./service/ai
  xuan gen model mongo ./service/ai agent`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		modelName := ""
		if len(args) >= 2 {
			modelName = args[1]
		}
		return gen.Model(dir, "mongo", modelName)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(genAPICmd)
	genRPCCmd.Flags().BoolVarP(&multiple, "multiple", "m", false, "Enable multi-service mode")
	genCmd.AddCommand(genRPCCmd)
	genCmd.AddCommand(genModelCmd)
	genModelCmd.AddCommand(genModelMongoCmd)
}
