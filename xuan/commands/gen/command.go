package gen

import (
	"community/xuan/types"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GenCommand 代码生成命令
type GenCommand struct {
	*types.BaseCommand
}

func NewGenCommand() *GenCommand {
	cmd := &GenCommand{
		BaseCommand: types.NewBaseCommand("gen", "代码生成命令，用于生成API或RPC相关代码", 1, nil),
	}

	// 添加子命令
	apiCmd := NewApiCommand(cmd)
	rpcCmd := NewRpcCommand(cmd)

	cmd.AddSon(apiCmd)
	cmd.AddSon(rpcCmd)

	return cmd
}

// ApiCommand API生成命令
type ApiCommand struct {
	*types.BaseCommand
}

func NewApiCommand(father *GenCommand) *ApiCommand {
	cmd := &ApiCommand{
		BaseCommand: types.NewBaseCommand("api", "在指定目录下递归处理 .api 文件，生成API代码", 2, father),
	}
	cmd.SetEnd(true) // 设置为最终命令
	return cmd
}

func (a *ApiCommand) Exec() error {
	// 这里需要处理参数，暂时显示帮助
	a.Help()
	return nil
}

// 处理API命令的具体执行
func (a *ApiCommand) ExecuteWithArgs(args []string) error {
	if len(args) == 0 {
		args = []string{"."} // 默认当前目录
	}

	path := args[0]
	return handleGenAPI(path)
}

// RpcCommand RPC生成命令
type RpcCommand struct {
	*types.BaseCommand
}

func NewRpcCommand(father *GenCommand) *RpcCommand {
	cmd := &RpcCommand{
		BaseCommand: types.NewBaseCommand("rpc", "在指定目录下递归处理 .proto 文件，生成RPC代码", 2, father),
	}
	cmd.SetEnd(true) // 设置为最终命令
	return cmd
}

func (r *RpcCommand) Exec() error {
	// 这里需要处理参数，暂时显示帮助
	r.Help()
	return nil
}

// 处理RPC命令的具体执行
func (r *RpcCommand) ExecuteWithArgs(args []string) error {
	if len(args) == 0 {
		args = []string{"."} // 默认当前目录
	}

	path := args[0]
	multiple := false

	// 检查是否有-m参数
	for i, arg := range args {
		if arg == "-m" {
			multiple = true
			// 移除-m参数
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	return handleGenRPC(path, multiple)
}

// 从oldmain.go移植的API生成功能
func handleGenAPI(rootDir string) error {
	if !isProjectRoot() {
		return fmt.Errorf("错误：xuan gen api 只能在项目根目录执行")
	}
	// 检查目录是否存在
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return fmt.Errorf("目录 %s 不存在", rootDir)
	}

	fmt.Printf("开始在 %s 目录下递归生成 API 代码...\n", rootDir)

	// 递归遍历目录
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果是目录，检查是否包含 .api 文件
		if info.IsDir() {
			apiFiles, err := filepath.Glob(filepath.Join(path, "*.api"))
			if err != nil {
				return fmt.Errorf("查找 API 文件出错: %v", err)
			}
			if len(apiFiles) > 0 {
				fmt.Printf("\n在目录 %s 中发现 %d 个 API 文件，开始生成代码...\n", path, len(apiFiles))

				args := []string{"api", "go", "-dir", ".", "--api"}
				for _, file := range apiFiles {
					relPath, err := filepath.Rel(path, file)
					if err != nil {
						relPath = file
					}
					args = append(args, relPath)
				}

				cmd := exec.Command("goctl", args...)
				cmd.Dir = path
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("执行 goctl 命令失败: %v\n输出: %s", err, string(output))
				}
				fmt.Printf("API 代码生成成功:\n%s\n", string(output))
			}
		}
		return nil
	})
}

// 从oldmain.go移植的RPC生成功能
func handleGenRPC(rootDir string, multiple bool) error {
	if !isProjectRoot() {
		return fmt.Errorf("错误：xuan gen rpc 只能在项目根目录执行")
	}

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return fmt.Errorf("目录 %s 不存在", rootDir)
	}

	fmt.Printf("开始在 %s 目录下递归生成 RPC 代码...\n", rootDir)

	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			protoFiles, err := filepath.Glob(filepath.Join(path, "*.proto"))
			if err != nil {
				return fmt.Errorf("查找 proto 文件出错: %v", err)
			}

			if len(protoFiles) > 0 {
				fmt.Printf("\n在目录 %s 中发现 %d 个 proto 文件，开始生成代码...\n", path, len(protoFiles))

				// 1. 构建基础参数（不含 -m 和 proto 文件）
				args := []string{
					"rpc", "protoc",
					"--go_out=./pb",
					"--go-grpc_out=./pb",
					"--zrpc_out=.",
					"--client=true",
				}

				// 2. 添加所有 proto 文件（放在 -m 之前）
				protoArgs := make([]string, 0, len(protoFiles))
				for _, file := range protoFiles {
					relPath, err := filepath.Rel(path, file)
					if err != nil {
						relPath = file
					}
					protoArgs = append(protoArgs, relPath)
				}
				args = append(args, protoArgs...)

				// 3. 最后添加 -m 参数（关键修复：与手动执行命令顺序一致）
				if multiple {
					args = append(args, "-m")
				}

				// 执行命令
				cmd := exec.Command("goctl", args...)
				cmd.Dir = path
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("执行 goctl 命令失败: %v\n输出: %s", err, string(output))
				}
				fmt.Printf("RPC 代码生成成功:\n%s\n", string(output))
			}
		}
		return nil
	})
}

func isProjectRoot() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}
