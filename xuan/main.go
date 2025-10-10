package main

import (
	"community/xuan/commands/end"
	"community/xuan/commands/gen"
	"community/xuan/commands/run"
	"community/xuan/types"
	"fmt"
	"os"
)

func main() {
	if !isProjectRoot() {
		fmt.Println("错误：xuan 命令只能在项目根目录执行")
		os.Exit(1)
	}

	// 初始化命令系统
	initCommands()

	// 创建命令执行器
	effector := types.NewEffector()

	// 执行命令
	if err := effector.Execute(os.Args[1:]); err != nil {
		fmt.Printf("执行失败: %v\n", err)
		os.Exit(1)
	}
}

func initCommands() {
	// 注册主要命令
	genCmd := gen.NewGenCommand()
	runCmd := run.NewRunCommand()
	endCmd := end.NewEndCommand()

	types.AddCommand(genCmd)
	types.AddCommand(runCmd)
	types.AddCommand(endCmd)
}

func isProjectRoot() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}
