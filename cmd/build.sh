#!/bin/bash

# 进入xuan目录，失败则退出
cd "$(dirname "${BASH_SOURCE[0]}")" || exit

# 编译程序
go build -o xuan main.go || {
    echo "编译失败"
    exit 1
}

# 获取当前目录的绝对路径
current_dir=$(pwd)
echo "程序已编译到: $current_dir"

# 确定Shell配置文件
if [ -n "$ZSH_VERSION" ]; then
    # Zsh环境
    config_file="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ]; then
    # Bash环境
    if [ -f "$HOME/.bash_profile" ]; then
        config_file="$HOME/.bash_profile"  # macOS优先使用
    else
        config_file="$HOME/.bashrc"        # Linux常用
    fi
else
    echo "未识别的Shell环境，无法自动配置PATH"
    exit 1
fi

# 检查当前目录是否已在PATH配置中
if ! grep -qxF "export PATH=\$PATH:$current_dir" "$config_file"; then
    # 路径不存在，添加新配置
    echo "export PATH=\$PATH:$current_dir" >> "$config_file"
    echo "已添加路径配置到 $config_file"
else
    # 路径已存在且正确，无需修改
    echo "路径已在 $config_file 中，无需更新"
fi
echo "如果无效，请执行以下命令使配置生效："
echo "source $config_file"
