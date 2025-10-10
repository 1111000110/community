package types

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

}

type Effector struct {
	NowCommand Command // 当前命令
	Level      int64   // 当前级别
}
