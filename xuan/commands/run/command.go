package run

import (
	"community/xuan/types"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RunCommand 运行命令
type RunCommand struct {
	*types.BaseCommand
}

func NewRunCommand() *RunCommand {
	cmd := &RunCommand{
		BaseCommand: types.NewBaseCommand("run", "运行/重启相关命令", 1, nil),
	}

	// 添加子命令
	allCmd := NewAllCommand(cmd)
	dockerCmd := NewDockerCommand(cmd)
	apiCmd := NewApiCommand(cmd)
	rpcCmd := NewRpcCommand(cmd)

	cmd.AddSon(allCmd)
	cmd.AddSon(dockerCmd)
	cmd.AddSon(apiCmd)
	cmd.AddSon(rpcCmd)

	return cmd
}

// AllCommand 启动所有服务命令
type AllCommand struct {
	*types.BaseCommand
}

func NewAllCommand(father *RunCommand) *AllCommand {
	cmd := &AllCommand{
		BaseCommand: types.NewBaseCommand("all", "启动所有数据库、etcd并运行所有进程（user/post 的 api 与 rpc）", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (a *AllCommand) Exec() error {
	return runAll()
}

// DockerCommand Docker相关命令
type DockerCommand struct {
	*types.BaseCommand
}

func NewDockerCommand(father *RunCommand) *DockerCommand {
	cmd := &DockerCommand{
		BaseCommand: types.NewBaseCommand("docker", "重启 docker-compose.yml（down 再 up -d）或单个服务", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (d *DockerCommand) Exec() error {
	// 这里需要处理参数，暂时显示帮助
	d.Help()
	return nil
}

// ApiCommand API运行命令
type ApiCommand struct {
	*types.BaseCommand
}

func NewApiCommand(father *RunCommand) *ApiCommand {
	cmd := &ApiCommand{
		BaseCommand: types.NewBaseCommand("api", "仅重启指定目录下的 API 服务", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (a *ApiCommand) Exec() error {
	// 这里需要处理参数，暂时显示帮助
	a.Help()
	return nil
}

// RpcCommand RPC运行命令
type RpcCommand struct {
	*types.BaseCommand
}

func NewRpcCommand(father *RunCommand) *RpcCommand {
	cmd := &RpcCommand{
		BaseCommand: types.NewBaseCommand("rpc", "仅重启指定目录下的 RPC 服务", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (r *RpcCommand) Exec() error {
	// 这里需要处理参数，暂时显示帮助
	r.Help()
	return nil
}

// 从oldmain.go移植的功能函数
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
