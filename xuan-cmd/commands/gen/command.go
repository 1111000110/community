package gen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// API generates API code in the given directory (recursively scans for .api files)
func API(dir string) error { return handleGenAPI(dir) }

// RPC generates RPC code in the given directory (recursively scans for .proto files)
func RPC(dir string, multiple bool) error { return handleGenRPC(dir, multiple) }

// --- internal implementation (migrated from old main) ---

func handleGenAPI(rootDir string) error {
	if !isProjectRoot() {
		return fmt.Errorf("错误：xuan gen api 只能在项目根目录执行")
	}
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return fmt.Errorf("目录 %s 不存在", rootDir)
	}

	fmt.Printf("开始在 %s 目录下递归生成 API 代码...\n", rootDir)

	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			apiFiles, err := filepath.Glob(filepath.Join(path, "*.api"))
			if err != nil {
				return fmt.Errorf("查找 API 文件出错: %v", err)
			}
			if len(apiFiles) > 0 {
				fmt.Printf("\n在目录 %s 中发现 %d 个 API 文件，开始生成代码...\n", path, len(apiFiles))

				templateHome := filepath.Join(getProjectRoot(), "xuan", "goctltemplate")
				// 注意：goctl 解析顺序要求在提供 --api 的文件列表之前放置 --home
				args := []string{"api", "go", "-dir", ".", "--home", templateHome, "--api"}
				for _, file := range apiFiles {
					relPath, err := filepath.Rel(path, file)
					if err != nil {
						relPath = file
					}
					args = append(args, relPath)
				}

				cmd := exec.Command("goctl", args...)
				// 默认使用项目内模板目录作为 GOCTL_HOME（若未显式设置）
				if _, exists := os.LookupEnv("GOCTL_HOME"); !exists {
					projectTemplateDir := filepath.Join(getProjectRoot(), "xuan", "goctltemplate")
					cmd.Env = append(os.Environ(), fmt.Sprintf("GOCTL_HOME=%s", projectTemplateDir))
				}
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

				templateHome := filepath.Join(getProjectRoot(), "xuan", "goctltemplate")
				args := []string{
					"rpc", "protoc",
					"--go_out=./pb",
					"--go-grpc_out=./pb",
					"--zrpc_out=.",
					"--client=true",
					"--home", templateHome,
				}
				protoArgs := make([]string, 0, len(protoFiles))
				for _, file := range protoFiles {
					relPath, err := filepath.Rel(path, file)
					if err != nil {
						relPath = file
					}
					protoArgs = append(protoArgs, relPath)
				}
				args = append(args, protoArgs...)
				if multiple {
					args = append(args, "-m")
				}

				cmd := exec.Command("goctl", args...)
				// 默认使用项目内模板目录作为 GOCTL_HOME（若未显式设置）
				if _, exists := os.LookupEnv("GOCTL_HOME"); !exists {
					projectTemplateDir := filepath.Join(getProjectRoot(), "xuan", "goctltemplate")
					cmd.Env = append(os.Environ(), fmt.Sprintf("GOCTL_HOME=%s", projectTemplateDir))
				}
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

// getProjectRoot 返回当前项目根目录的绝对路径（基于 go.mod 存在性判定）
func getProjectRoot() string {
	// 从当前工作目录向上查找包含 go.mod 的目录
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir { // 到达根目录
			break
		}
		dir = parent
	}
	return wd
}
