package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if !isProjectRoot() {
		fmt.Println("错误：xuan 命令只能在项目根目录执行")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printMainUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "gen":
		handleGenCommand()
	default:
		printMainUsage()
		os.Exit(1)
	}
}

func handleGenCommand() {
	if len(os.Args) < 3 {
		printGenUsage()
		os.Exit(1)
	}

	subCmd := os.Args[2]
	switch subCmd {
	case "api":
		path := "."
		if len(os.Args) > 3 {
			path = os.Args[3]
		}
		path = os.Args[3]
		if err := handleGenAPI(path); err != nil {
			fmt.Printf("执行命令失败: %v\n", err)
			os.Exit(1)
		}
	case "rpc":
		// 解析 -m 参数（支持在路径前后）
		rpcFlagSet := flag.NewFlagSet("rpc", flag.ExitOnError)
		multiple := rpcFlagSet.Bool("m", false, "启用多服务模式")

		// 从第3个参数开始解析（跳过 xuan gen rpc）
		if err := rpcFlagSet.Parse(os.Args[3:]); err != nil {
			fmt.Printf("参数解析错误: %v\n", err)
			os.Exit(1)
		}

		// 获取路径参数（默认当前目录）
		path := "."
		args := rpcFlagSet.Args()
		if len(args) > 0 {
			path = args[0]
		}

		if err := handleGenRPC(path, *multiple); err != nil {
			fmt.Printf("执行命令失败: %v\n", err)
			os.Exit(1)
		}
	default:
		printGenUsage()
		os.Exit(1)
	}
}

// 处理 gen rpc 命令，关键修复：将 -m 放在 proto 文件之后
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

// 处理 gen api 命令
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

func isProjectRoot() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}

// 一级命令帮助信息
func printMainUsage() {
	usage := `用法: xuan [命令]

一级命令:
  gen    代码生成命令，用于生成API或RPC相关代码
         可通过 xuan gen 查看二级命令详情

示例:
  xuan gen api    生成API代码
  xuan gen rpc    生成RPC代码
`
	fmt.Print(usage)
}

// 二级命令(gen)帮助信息
func printGenUsage() {
	usage := `用法: xuan gen [子命令] [参数]

子命令:
  api <目录路径>    在指定目录下递归处理 .api 文件，生成API代码
  rpc <目录路径>    在指定目录下递归处理 .proto 文件，生成RPC代码
                    可选参数: -m  启用多服务模式（解决 multiple services 问题）

示例:
  xuan gen api ./service/message        生成API代码
  xuan gen rpc ./service/user           生成RPC代码（单服务模式）
  xuan gen rpc ./service/user -m        生成RPC代码（多服务模式）
`
	fmt.Print(usage)
}
