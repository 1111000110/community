package gen

import (
	"community/cmd/types"
	"fmt"
)

type Gen struct {
}

func NewGen() *Gen {
	return &Gen{}
}

func (g *Gen) Help() {
	fmt.Println("gen 用于生成代码，支持以下子命令：\n" +
		"gen api < 目录路径 > -- 在指定目录下寻找.api 文件，生成 go-zero 框架的 api 代码 \n" +
		"示例：xuan gen api ./service/chatgroup\n" +
		"gen rpc < 目录路径 > -- 在指定目录下寻找.proto 文件，生成 go-zero 框架的 rpc 代码 \n" +
		"示例：xuan gen rpc ./service/user\n" +
		"详细介绍请查看./cmd/gen/README.md")
}
func (g *Gen) Exec() error {
	g.Help()
	return nil
}
func (g *Gen) IsEnd() bool {
	return false
}
func (g *Gen) Son() map[string]types.Command {
	return map[string]types.Command{}
}
