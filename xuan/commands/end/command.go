package end

import (
	"community/xuan/types"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EndCommand 停止命令
type EndCommand struct {
	*types.BaseCommand
}

func NewEndCommand() *EndCommand {
	cmd := &EndCommand{
		BaseCommand: types.NewBaseCommand("end", "停止相关服务", 1, nil),
	}

	// 添加子命令
	allCmd := NewAllCommand(cmd)
	dockerCmd := NewDockerCommand(cmd)
	apiCmd := NewApiCommand(cmd)
	rpcCmd := NewRpcCommand(cmd)
	etcdCmd := NewEtcdCommand(cmd)
	cleanCmd := NewCleanCommand(cmd)

	cmd.AddSon(allCmd)
	cmd.AddSon(dockerCmd)
	cmd.AddSon(apiCmd)
	cmd.AddSon(rpcCmd)
	cmd.AddSon(etcdCmd)
	cmd.AddSon(cleanCmd)

	return cmd
}

// AllCommand 停止所有服务命令
type AllCommand struct {
	*types.BaseCommand
}

func NewAllCommand(father *EndCommand) *AllCommand {
	cmd := &AllCommand{
		BaseCommand: types.NewBaseCommand("all", "停止所有进程、停止 etcd，并 docker-compose down", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (a *AllCommand) Exec() error {
	return stopAllServices()
}

// 实现ExecuteWithArgs方法以支持参数
func (a *AllCommand) ExecuteWithArgs(args []string) error {
	cleanData := false

	// 解析参数
	for _, arg := range args {
		if arg == "--clean-data" || arg == "-c" {
			cleanData = true
		}
	}

	// 停止所有服务
	if err := stopAllServices(); err != nil {
		return err
	}

	// 如果需要清理数据
	if cleanData {
		return cleanAllData()
	}

	return nil
}

// DockerCommand Docker停止命令
type DockerCommand struct {
	*types.BaseCommand
}

func NewDockerCommand(father *EndCommand) *DockerCommand {
	cmd := &DockerCommand{
		BaseCommand: types.NewBaseCommand("docker", "docker-compose down 或停止单个服务", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (d *DockerCommand) Exec() error {
	// 这里需要处理参数，暂时显示帮助
	d.Help()
	return nil
}

// ApiCommand API停止命令
type ApiCommand struct {
	*types.BaseCommand
}

func NewApiCommand(father *EndCommand) *ApiCommand {
	cmd := &ApiCommand{
		BaseCommand: types.NewBaseCommand("api", "停止 API 服务", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (a *ApiCommand) Exec() error {
	// 这里需要处理参数，暂时显示帮助
	a.Help()
	return nil
}

// RpcCommand RPC停止命令
type RpcCommand struct {
	*types.BaseCommand
}

func NewRpcCommand(father *EndCommand) *RpcCommand {
	cmd := &RpcCommand{
		BaseCommand: types.NewBaseCommand("rpc", "停止 RPC 服务", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (r *RpcCommand) Exec() error {
	// 这里需要处理参数，暂时显示帮助
	r.Help()
	return nil
}

// EtcdCommand Etcd停止命令
type EtcdCommand struct {
	*types.BaseCommand
}

func NewEtcdCommand(father *EndCommand) *EtcdCommand {
	cmd := &EtcdCommand{
		BaseCommand: types.NewBaseCommand("etcd", "停止 etcd（通过 pkill etcd）", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (e *EtcdCommand) Exec() error {
	return stopEtcd()
}

// CleanCommand 清理命令
type CleanCommand struct {
	*types.BaseCommand
}

func NewCleanCommand(father *EndCommand) *CleanCommand {
	cmd := &CleanCommand{
		BaseCommand: types.NewBaseCommand("clean", "清理数据、日志和二进制文件", 2, father),
	}
	cmd.SetEnd(true)
	return cmd
}

func (c *CleanCommand) Exec() error {
	return cleanAllData()
}

// 实现ExecuteWithArgs方法以支持参数
func (c *CleanCommand) ExecuteWithArgs(args []string) error {
	if len(args) == 0 {
		// 默认清理所有数据
		return cleanAllData()
	}

	// 根据参数清理特定内容
	for _, arg := range args {
		switch arg {
		case "logs":
			return cleanLogs()
		case "data":
			return cleanDataDirectories()
		case "binaries":
			return cleanBinaries()
		case "all":
			return cleanAllData()
		default:
			return fmt.Errorf("未知的清理选项: %s", arg)
		}
	}

	return nil
}

// 核心功能函数 - 重新设计为更灵活的方式

// stopAllServices 停止所有服务
func stopAllServices() error {
	fmt.Println("停止所有服务...")

	// 1. 停止Docker服务
	fmt.Println("停止 Docker Compose 服务...")
	if err := runCmdStreaming("docker-compose", "down"); err != nil {
		fmt.Printf("停止 Docker 服务失败: %v\n", err)
	}

	// 2. 停止etcd
	fmt.Println("停止 etcd...")
	if err := stopEtcd(); err != nil {
		fmt.Printf("停止 etcd 失败: %v\n", err)
	}

	// 3. 停止所有业务服务
	fmt.Println("停止业务服务...")
	if err := stopAllBusinessServices(); err != nil {
		fmt.Printf("停止业务服务失败: %v\n", err)
	}

	fmt.Println("所有服务已停止")
	return nil
}

// stopAllBusinessServices 停止所有业务服务（动态发现）
func stopAllBusinessServices() error {
	// 扫描service目录，动态发现服务
	serviceDir := "service"
	if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
		fmt.Println("service 目录不存在，跳过业务服务停止")
		return nil
	}

	return filepath.Walk(serviceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是api或rpc目录
		if info.IsDir() && (info.Name() == "api" || info.Name() == "rpc") {
			if err := stopServicePath(path); err != nil {
				fmt.Printf("停止服务 %s 失败: %v\n", path, err)
			}
		}

		return nil
	})
}

// stopServicePath 停止指定路径的服务
func stopServicePath(servicePath string) error {
	serviceName := deriveServiceName(servicePath)
	fmt.Printf("停止服务: %s (%s)\n", serviceName, servicePath)

	// 停止进程
	_ = runCmdSilent("pkill", "-f", serviceName)

	// 清理二进制文件
	binaryPath := filepath.Join(servicePath, serviceName)
	if _, err := os.Stat(binaryPath); err == nil {
		if err := os.Remove(binaryPath); err != nil {
			fmt.Printf("删除二进制文件失败: %v\n", err)
		}
	}

	return nil
}

// deriveServiceName 根据路径生成服务名
func deriveServiceName(servicePath string) string {
	dirName := filepath.Base(servicePath)
	parent := filepath.Base(filepath.Dir(servicePath))
	return parent + "-" + dirName
}

// stopEtcd 停止etcd
func stopEtcd() error {
	return runCmdSilent("pkill", "-f", "etcd")
}

// cleanAllData 清理所有数据
func cleanAllData() error {
	fmt.Println("清理所有数据...")

	// 清理日志
	if err := cleanLogs(); err != nil {
		fmt.Printf("清理日志失败: %v\n", err)
	}

	// 清理数据目录
	if err := cleanDataDirectories(); err != nil {
		fmt.Printf("清理数据目录失败: %v\n", err)
	}

	// 清理二进制文件
	if err := cleanBinaries(); err != nil {
		fmt.Printf("清理二进制文件失败: %v\n", err)
	}

	fmt.Println("数据清理完成")
	return nil
}

// cleanLogs 清理日志文件
func cleanLogs() error {
	fmt.Println("清理日志文件...")

	// 清理根目录日志
	_ = runCmdSilent("rm", "-f", "*.log")

	// 清理etcd日志
	_ = runCmdSilent("rm", "-f", "etcd.log")

	// 清理服务日志
	_ = runCmdSilent("find", "service", "-name", "*.log", "-delete")

	return nil
}

// cleanDataDirectories 清理数据目录
func cleanDataDirectories() error {
	fmt.Println("清理数据目录...")

	// 清理etcd数据
	if _, err := os.Stat("etcd-data"); err == nil {
		if err := os.RemoveAll("etcd-data"); err != nil {
			fmt.Printf("删除 etcd-data 失败: %v\n", err)
		} else {
			fmt.Println("已删除 etcd-data 目录")
		}
	}

	// 清理docker数据
	if _, err := os.Stat("docker-data"); err == nil {
		if err := os.RemoveAll("docker-data"); err != nil {
			fmt.Printf("删除 docker-data 失败: %v\n", err)
		} else {
			fmt.Println("已删除 docker-data 目录")
		}
	}

	return nil
}

// cleanBinaries 清理二进制文件
func cleanBinaries() error {
	fmt.Println("清理二进制文件...")

	// 扫描service目录，清理所有二进制文件
	return filepath.Walk("service", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是二进制文件（可执行文件）
		if !info.IsDir() && info.Mode()&0111 != 0 {
			// 检查是否是服务二进制文件（根据命名规则）
			fileName := info.Name()
			if strings.Contains(fileName, "-api") || strings.Contains(fileName, "-rpc") {
				if err := os.Remove(path); err != nil {
					fmt.Printf("删除二进制文件 %s 失败: %v\n", path, err)
				} else {
					fmt.Printf("已删除二进制文件: %s\n", path)
				}
			}
		}

		return nil
	})
}

// 工具函数
func runCmdStreaming(name string, args ...string) error {
	fmt.Printf("$ %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmdSilent(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}
