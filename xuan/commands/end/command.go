package end

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// StopAllServicesWrapper stops all services and infra
func StopAllServicesWrapper() error { return stopAllServices() }

// CleanAllDataWrapper cleans all logs, data, and binaries
func CleanAllDataWrapper() error { return cleanAllData() }

func CleanLogsWrapper() error { return cleanLogs() }
func CleanDataDirectoriesWrapper() error { return cleanDataDirectories() }
func CleanBinariesWrapper() error { return cleanBinaries() }
func RunCmdStreamingWrapper(name string, args ...string) error { return runCmdStreaming(name, args...) }

// --- internal implementation (refactored for flexibility) ---

func stopAllServices() error {
	fmt.Println("停止所有服务...")
	fmt.Println("停止 Docker Compose 服务...")
	if err := runCmdStreaming("docker-compose", "down"); err != nil {
		fmt.Printf("停止 Docker 服务失败: %v\n", err)
	}
	fmt.Println("停止 etcd...")
	if err := runCmdSilent("pkill", "-f", "etcd"); err != nil {
		fmt.Printf("停止 etcd 失败: %v\n", err)
	}
	fmt.Println("停止业务服务...")
	if err := stopAllBusinessServices(); err != nil {
		fmt.Printf("停止业务服务失败: %v\n", err)
	}
	fmt.Println("所有服务已停止")
	return nil
}

func stopAllBusinessServices() error {
	serviceDir := "service"
	if _, err := os.Stat(serviceDir); os.IsNotExist(err) { return nil }
	return filepath.Walk(serviceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if info.IsDir() && (info.Name() == "api" || info.Name() == "rpc") {
			if err := stopServicePath(path); err != nil { fmt.Printf("停止服务 %s 失败: %v\n", path, err) }
		}
		return nil
	})
}

func stopServicePath(servicePath string) error {
	serviceName := deriveServiceName(servicePath)
	fmt.Printf("停止服务: %s (%s)\n", serviceName, servicePath)
	_ = runCmdSilent("pkill", "-f", serviceName)
	binaryPath := filepath.Join(servicePath, serviceName)
	if _, err := os.Stat(binaryPath); err == nil {
		if err := os.Remove(binaryPath); err != nil { fmt.Printf("删除二进制文件失败: %v\n", err) }
	}
	return nil
}

func deriveServiceName(servicePath string) string {
	dirName := filepath.Base(servicePath)
	parent := filepath.Base(filepath.Dir(servicePath))
	return parent + "-" + dirName
}

func cleanAllData() error {
	fmt.Println("清理所有数据...")
	if err := cleanLogs(); err != nil { fmt.Printf("清理日志失败: %v\n", err) }
	if err := cleanDataDirectories(); err != nil { fmt.Printf("清理数据目录失败: %v\n", err) }
	if err := cleanBinaries(); err != nil { fmt.Printf("清理二进制文件失败: %v\n", err) }
	fmt.Println("数据清理完成")
	return nil
}

func cleanLogs() error {
	fmt.Println("清理日志文件...")
	_ = runCmdSilent("rm", "-f", "*.log")
	_ = runCmdSilent("rm", "-f", "etcd.log")
	_ = runCmdSilent("find", "service", "-name", "*.log", "-delete")
	return nil
}

func cleanDataDirectories() error {
	fmt.Println("清理数据目录...")
	if _, err := os.Stat("etcd-data"); err == nil { _ = os.RemoveAll("etcd-data") }
	if _, err := os.Stat("docker-data"); err == nil { _ = os.RemoveAll("docker-data") }
	return nil
}

func cleanBinaries() error {
	fmt.Println("清理二进制文件...")
	return filepath.Walk("service", func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if !info.IsDir() && info.Mode()&0111 != 0 {
			fileName := info.Name()
			if strings.Contains(fileName, "-api") || strings.Contains(fileName, "-rpc") {
				if err := os.Remove(path); err != nil { fmt.Printf("删除二进制文件 %s 失败: %v\n", path, err) } else { fmt.Printf("已删除二进制文件: %s\n", path) }
			}
		}
		return nil
	})
}

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
