package types

import (
	"fmt"
)

type Command interface {
	Exec() error             // 执行
	Father() Command         // 父命令
	Son() map[string]Command // 子命令列表
	Help()                   // 帮助
	IsEnd() bool             // 是否是最终命令
	Level() int64            // 级别
	String() string          // 命令名称
}

var VEffectors = map[string]Command{}

func AddCommand(cmd Command) {
	VEffectors[cmd.String()] = cmd
}

// BaseCommand 基础命令实现
type BaseCommand struct {
	name        string
	description string
	level       int64
	father      Command
	sons        map[string]Command
	isEnd       bool
}

func NewBaseCommand(name, description string, level int64, father Command) *BaseCommand {
	return &BaseCommand{
		name:        name,
		description: description,
		level:       level,
		father:      father,
		sons:        make(map[string]Command),
		isEnd:       false,
	}
}

func (b *BaseCommand) Exec() error {
	if b.isEnd {
		return fmt.Errorf("命令 %s 是最终命令，需要提供子命令", b.name)
	}
	b.Help()
	return nil
}

func (b *BaseCommand) Father() Command {
	return b.father
}

func (b *BaseCommand) Son() map[string]Command {
	return b.sons
}

func (b *BaseCommand) Help() {
	fmt.Printf("用法: %s [子命令]\n\n", b.name)
	fmt.Printf("描述: %s\n\n", b.description)

	if len(b.sons) > 0 {
		fmt.Println("子命令:")
		for name, cmd := range b.sons {
			fmt.Printf("  %-10s %s\n", name, getCommandDescription(cmd))
		}
	}
}

func (b *BaseCommand) IsEnd() bool {
	return b.isEnd
}

func (b *BaseCommand) Level() int64 {
	return b.level
}

func (b *BaseCommand) String() string {
	return b.name
}

func (b *BaseCommand) AddSon(cmd Command) {
	b.sons[cmd.String()] = cmd
}

func (b *BaseCommand) SetEnd(isEnd bool) {
	b.isEnd = isEnd
}

func getCommandDescription(cmd Command) string {
	if baseCmd, ok := cmd.(*BaseCommand); ok {
		return baseCmd.description
	}
	return "无描述"
}

// Effector 命令执行器
type Effector struct {
	NowCommand Command // 当前命令
	Level      int64   // 当前级别
}

func NewEffector() *Effector {
	return &Effector{
		Level: 0,
	}
}

func (e *Effector) Execute(args []string) error {
	if len(args) == 0 {
		e.printMainUsage()
		return nil
	}

	// 从根命令开始查找
	cmdName := args[0]
	cmd, exists := VEffectors[cmdName]
	if !exists {
		return fmt.Errorf("未知命令: %s", cmdName)
	}

	// 递归查找子命令
	currentCmd := cmd
	for i := 1; i < len(args); i++ {
		subCmdName := args[i]
		sons := currentCmd.Son()
		subCmd, exists := sons[subCmdName]
		if !exists {
			// 如果没有找到子命令，检查是否是最终命令
			if currentCmd.IsEnd() {
				// 执行最终命令，传递剩余参数
				return e.executeEndCommand(currentCmd, args[i:])
			}
			return fmt.Errorf("未知子命令: %s", subCmdName)
		}
		currentCmd = subCmd
	}

	// 执行命令
	return currentCmd.Exec()
}

func (e *Effector) executeEndCommand(cmd Command, args []string) error {
	// 检查是否是帮助请求
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		cmd.Help()
		return nil
	}

	// 使用接口方法执行命令，避免循环导入
	// 具体命令需要实现ExecuteWithArgs方法
	if executor, ok := cmd.(interface{ ExecuteWithArgs([]string) error }); ok {
		return executor.ExecuteWithArgs(args)
	}

	// 如果没有ExecuteWithArgs方法，则执行默认的Exec方法
	return cmd.Exec()
}

func (e *Effector) printMainUsage() {
	fmt.Println("用法: xuan [命令]")
	fmt.Println()
	fmt.Println("一级命令:")
	for name, cmd := range VEffectors {
		if baseCmd, ok := cmd.(*BaseCommand); ok {
			fmt.Printf("  %-10s %s\n", name, baseCmd.description)
		}
	}
	fmt.Println()
	fmt.Println("使用 'xuan <命令> -h' 查看具体命令的帮助信息")
}
