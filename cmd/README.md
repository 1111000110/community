# `xuan` 命令介绍
`xuan` 是本项目的脚手架工具，支持基于文件生成代码、连接数据库、运行代码等核心功能。
## 生成代码
`gen` 为 `xuan` 的子命令，用于生成代码。该功能基于 `go-zero` 的 `goctl` 命令实现，`xuan` 对此进行了封装，详情可参考 community/gen 目录。
```shell
xuan gen <模式> [目录路径|选项]
# 模式说明：
#   api：基于目录下的 API 定义文件（如 .api 格式）生成对应代码
#   rpc：基于目录下的 Proto 定义文件（.proto 格式）生成对应 RPC 代码
# 选项说明：
#   -h：显示对应模式的帮助信息（无需提供目录路径）

# 示例：
# 生成代码（需指定目录路径，必选）
xuan gen api ./service/message  # 根据 ./service/message 目录下的 API 定义文件，在该目录生成对应代码
xuan gen rpc ./service/message  # 根据 ./service/message 目录下的 Proto 定义文件，在该目录生成对应 RPC 代码
xuan gen api ./                 # 递归处理当前目录下所有 API 定义文件，在各文件所在目录生成对应代码
xuan gen rpc ./                 # 递归处理当前目录下所有 Proto 定义文件，在各文件所在目录生成对应 RPC 代码

# 查看帮助（无需目录路径，选项优先）
xuan gen api -h  # 打印 api 模式的帮助信息
xuan gen rpc -h  # 打印 rpc 模式的帮助信息
```
## 执行代码
`run` 为 `xuan` 的子命令，用于执行代码，如果已经启动则为重启。
```shell
xuan run <执行模式|全局参数> [目录路径|镜像名称]
# 执行模式说明：
#   go：执行指定代码目录下的 main 程序（若已启动则重启）
#   docker：通过 docker 启动指定镜像（若已启动则重启）
# 全局参数：
#   -a：启动所有程序（go 模式）和所有镜像（docker 模式），无需指定目录或镜像
#   -h：查看对应模式的帮助信息

# 示例：
# 启动代码
xuan run go ./service/message  # 启动 ./service/message 目录下的 main 程序（若已启动则重启）
xuan run go -a                 # 启动所有 go 模式的程序（若已启动则重启）
# 启动镜像
xuan run docker redis          # 启动 redis 镜像（若已启动则重启）
xuan run docker -a             # 启动所有 docker 镜像（若已启动则重启）
# 启动全部
xuan run -a                    # 启动所有 go 程序和所有 docker 镜像（若已启动则重启）
# 查看帮助
xuan run go -h                 # 查看 go 模式的帮助信息
xuan run docker -h             # 查看 docker 模式的帮助信息
xuan run -h                    # 查看 run 命令的全局帮助信息
```