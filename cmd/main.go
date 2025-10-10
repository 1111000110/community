package main

import (
	"fmt"
	"os"
)

func main() {
	if !isProjectRoot() {
		fmt.Println("错误：xuan 命令只能在项目根目录执行")
		os.Exit(1)
	}
	if len(os.Args) < 2 {

	}
}
func isProjectRoot() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}
