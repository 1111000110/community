package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		printMainUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "gen":
		handleGenCommand()
	case "run":
		handleRunCommand()
	case "end":
		handleEndCommand()
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

// 一级命令帮助信息
func printMainUsage() {
	usage := `用法: xuan [命令]

一级命令:
  gen    代码生成命令，用于生成API或RPC相关代码
         可通过 xuan gen 查看二级命令详情
  run    运行/重启相关命令
         可通过 xuan run -h 查看说明
  end    停止相关服务
         可通过 xuan end -h 查看说明
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

// ========= RUN =========

func printRunUsage() {
	usage := `用法: xuan run [子命令] [...参数]

子命令:
  -h                      查看帮助
  all                     启动所有数据库、etcd并运行所有进程（user/post 的 api 与 rpc）
  docker                  重启 docker-compose.yml（down 再 up -d）
  docker <svc>            重启 docker-compose.yml 中的单个服务（如 kafka、mongo 等）
  api rpc <path>          重启指定目录下的 api 与 rpc （例如: ./service/user）
  api <path>              仅重启 api
  rpc <path>              仅重启 rpc

示例:
  xuan run -h
  xuan run all
  xuan run docker
  xuan run docker kafka
  xuan run api rpc ./service/user
  xuan run api ./service/user
  xuan run rpc ./service/user
`
	fmt.Print(usage)
}

func handleRunCommand() {
	if len(os.Args) < 3 || os.Args[2] == "-h" || os.Args[2] == "--help" {
		printRunUsage()
		return
	}

	sub := os.Args[2]
	switch sub {
	case "all":
		if err := runAll(); err != nil {
			fmt.Printf("执行失败: %v\n", err)
			os.Exit(1)
		}
	case "docker":
		// xuan run docker [service]
		if len(os.Args) == 3 {
			if err := dockerComposeDownUp(); err != nil {
				fmt.Printf("执行失败: %v\n", err)
				os.Exit(1)
			}
			return
		}
		if len(os.Args) >= 4 {
			svc := os.Args[3]
			if err := dockerComposeRestartService(svc); err != nil {
				fmt.Printf("执行失败: %v\n", err)
				os.Exit(1)
			}
			return
		}
	case "api":
		// xuan run api [rpc] <path>
		if len(os.Args) < 4 {
			printRunUsage()
			os.Exit(1)
		}
		if os.Args[3] == "rpc" {
			if len(os.Args) < 5 {
				printRunUsage()
				os.Exit(1)
			}
			path := os.Args[4]
			if err := restartServices(path, true, true); err != nil {
				fmt.Printf("执行失败: %v\n", err)
				os.Exit(1)
			}
			return
		}
		path := os.Args[3]
		if err := restartServices(path, true, false); err != nil {
			fmt.Printf("执行失败: %v\n", err)
			os.Exit(1)
		}
	case "rpc":
		if len(os.Args) < 4 {
			printRunUsage()
			os.Exit(1)
		}
		path := os.Args[3]
		if err := restartServices(path, false, true); err != nil {
			fmt.Printf("执行失败: %v\n", err)
			os.Exit(1)
		}
	default:
		printRunUsage()
		os.Exit(1)
	}
}

// ========= END =========

func printEndUsage() {
	usage := `用法: xuan end [子命令] [...参数]

子命令:
  -h                      查看帮助
  all                     停止所有进程、停止 etcd，并 docker-compose down
  docker                  docker-compose down
  docker <svc>            docker-compose stop <svc>
  api rpc <path>          停止指定目录下的 api 与 rpc
  api <path>              停止 api
  rpc <path>              停止 rpc
  etcd                    停止 etcd（通过 pkill etcd）

示例:
  xuan end -h
  xuan end all
  xuan end docker
  xuan end docker kafka
  xuan end api rpc ./service/user
  xuan end api ./service/user
  xuan end rpc ./service/user
  xuan end etcd
`
	fmt.Print(usage)
}

func handleEndCommand() {
	if len(os.Args) < 3 || os.Args[2] == "-h" || os.Args[2] == "--help" {
		printEndUsage()
		return
	}
	sub := os.Args[2]
	switch sub {
	case "all":
		// stop all known services, etcd, and docker-compose down
		_ = stopKnownServices()
		_ = stopEtcd()
		_ = runCmdStreaming("docker-compose", "down")
	case "docker":
		if len(os.Args) == 3 {
			_ = runCmdStreaming("docker-compose", "down")
			return
		}
		svc := os.Args[3]
		_ = runCmdStreaming("docker-compose", "stop", svc)
	case "api":
		if len(os.Args) < 4 {
			printEndUsage()
			os.Exit(1)
		}
		if os.Args[3] == "rpc" {
			if len(os.Args) < 5 {
				printEndUsage()
				os.Exit(1)
			}
			path := os.Args[4]
			_ = stopServices(path, true, true)
			return
		}
		path := os.Args[3]
		_ = stopServices(path, true, false)
	case "rpc":
		if len(os.Args) < 4 {
			printEndUsage()
			os.Exit(1)
		}
		path := os.Args[3]
		_ = stopServices(path, false, true)
	case "etcd":
		_ = stopEtcd()
	default:
		printEndUsage()
		os.Exit(1)
	}
}

// ========= 实现 =========

func runAll() error {
	fmt.Println("重启 docker-compose 服务...")
	if err := dockerComposeDownUp(); err != nil {
		return err
	}

	fmt.Println("等待核心端口可用: mongo(27017), mysql(3306), redis(6379), scylla(9042), kafka(9092)...")
	waitPorts := []string{"127.0.0.1:27017", "127.0.0.1:3306", "127.0.0.1:6379", "127.0.0.1:9042", "127.0.0.1:9092"}
	for _, addr := range waitPorts {
		if err := waitForTCP(addr, 120*time.Second); err != nil {
			return fmt.Errorf("等待端口 %s 失败: %w", addr, err)
		}
	}

	fmt.Println("启动 etcd 并等待 2379...")
	if err := startEtcd(); err != nil {
		return err
	}
	if err := waitForTCP("127.0.0.1:2379", 60*time.Second); err != nil {
		return fmt.Errorf("等待 etcd 2379 失败: %w", err)
	}

	// 启动业务进程
	targets := []string{
		"service/user/rpc",
		"service/user/api",
		"service/post/rpc",
		"service/post/api",
	}
	for _, p := range targets {
		if err := restartPath(p); err != nil {
			return err
		}
	}
	fmt.Println("全部启动完成")
	return nil
}

func dockerComposeDownUp() error {
	if err := runCmdStreaming("docker-compose", "down"); err != nil {
		return err
	}
	return runCmdStreaming("docker-compose", "up", "-d")
}

func dockerComposeRestartService(svc string) error {
	// 优先用 restart；若不支持，可退化为 stop/start
	if err := runCmdStreaming("docker-compose", "restart", svc); err != nil {
		// fallback
		_ = runCmdStreaming("docker-compose", "stop", svc)
		return runCmdStreaming("docker-compose", "up", "-d", svc)
	}
	return nil
}

func restartServices(root string, doAPI, doRPC bool) error {
	if doRPC {
		if err := restartPath(filepath.Join(root, "rpc")); err != nil {
			return err
		}
	}
	if doAPI {
		if err := restartPath(filepath.Join(root, "api")); err != nil {
			return err
		}
	}
	return nil
}

func stopServices(root string, doAPI, doRPC bool) error {
	if doRPC {
		_ = stopPath(filepath.Join(root, "rpc"))
	}
	if doAPI {
		_ = stopPath(filepath.Join(root, "api"))
	}
	return nil
}

func restartPath(servicePath string) error {
	if err := stopPath(servicePath); err != nil {
		// ignore
	}
	return buildAndRunService(servicePath)
}

func stopPath(servicePath string) error {
	serviceName := deriveServiceName(servicePath)
	binPath := filepath.Join(servicePath, serviceName)
	// 使用 pkill -f 来匹配完整路径
	_ = runCmdSilent("pkill", "-f", binPath)
	// 再补一刀，按名字
	_ = runCmdSilent("pkill", "-f", serviceName)
	return nil
}

func buildAndRunService(servicePath string) error {
	if _, err := os.Stat(servicePath); os.IsNotExist(err) {
		return fmt.Errorf("目录不存在: %s", servicePath)
	}
	serviceName := deriveServiceName(servicePath)
	fmt.Printf("构建并启动: %s (%s)\n", serviceName, servicePath)

	// 构建当前包
	if err := runCmdInDirStreaming(servicePath, "go", "build", "-o", serviceName, "."); err != nil {
		return err
	}

	// 运行并将日志输出到上级目录
	logFile := filepath.Join(filepath.Dir(servicePath), serviceName+".log")
	binPath := filepath.Join(servicePath, serviceName)
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}
	cmd := exec.Command(binPath)
	cmd.Stdout = f
	cmd.Stderr = f
	if err := cmd.Start(); err != nil {
		_ = f.Close()
		return fmt.Errorf("启动进程失败: %v", err)
	}
	// 不等待，直接返回
	_ = f.Close()
	return nil
}

func deriveServiceName(servicePath string) string {
	dirName := filepath.Base(servicePath)
	parent := filepath.Base(filepath.Dir(servicePath))
	return parent + "-" + dirName
}

// etcd 管理
func startEtcd() error {
	// 等价: etcd --data-dir=./etcd-data > etcd.log 2>&1 &
	logFile := "etcd.log"
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("打开 etcd 日志失败: %v", err)
	}
	cmd := exec.Command("etcd", "--data-dir=./etcd-data")
	cmd.Stdout = f
	cmd.Stderr = f
	if err := cmd.Start(); err != nil {
		_ = f.Close()
		return fmt.Errorf("启动 etcd 失败: %v", err)
	}
	_ = f.Close()
	return nil
}

func stopEtcd() error {
	// 直接按名字杀
	return runCmdSilent("pkill", "-f", "etcd")
}

// 工具函数
func waitForTCP(address string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("等待超时: %s", address)
}

func runCmdStreaming(name string, args ...string) error {
	fmt.Printf("$ %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmdInDirStreaming(dir string, name string, args ...string) error {
	fmt.Printf("(%s) $ %s %s\n", dir, name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmdSilent(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func stopKnownServices() error {
	// 基于命名规则停止常见进程
	names := []string{
		"user-api", "user-rpc", "post-api", "post-rpc",
	}
	for _, n := range names {
		_ = runCmdSilent("pkill", "-f", n)
	}
	return nil
}
