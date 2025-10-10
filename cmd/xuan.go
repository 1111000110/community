package main

import "fmt"

type xuan struct{}

func (x *xuan) Help() {
	fmt.Println("xuan 为项目帮助路径，支持以下子命令：\n" +
		"生成代码：xuan gen <代码类型> < 目录路径 >  \n" +
		"示例：xuan gen api ./service/chatgroup\n" +
		"详细介绍请查看./cmd/README.md")
}
