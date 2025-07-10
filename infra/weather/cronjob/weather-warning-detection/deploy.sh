#!/bin/bash

# 天气预警检测系统部署脚本

echo "=== 天气预警检测系统部署脚本 ==="

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go环境，请先安装Go"
    exit 1
fi

echo "1. 检查Go环境... ✓"

# 编译程序
echo "2. 编译程序..."
go build -o weather-warning-detection main.go

if [ $? -eq 0 ]; then
    echo "   编译成功 ✓"
else
    echo "   编译失败 ✗"
    exit 1
fi

# 检查配置文件
if [ ! -f "config.json" ]; then
    echo "3. 创建默认配置文件..."
    cat > config.json << EOF
{
  "weather_api": {
    "key": "YOUR_AMAP_API_KEY",
    "city_code": "110105",
    "city_name": "北京市朝阳区",
    "send_name": "用户",
    "open_id": "YOUR_OPEN_ID"
  },
  "feishu": {
    "webhook_url": ""
  }
}
EOF
    echo "   配置文件已创建，请编辑 config.json 设置您的API密钥和城市信息 ⚠️"
else
    echo "3. 配置文件已存在 ✓"
fi

# 设置执行权限
chmod +x weather-warning-detection
echo "4. 设置执行权限 ✓"

# 测试运行
echo "5. 测试运行程序..."
echo "   如果配置正确，程序将获取天气信息并显示结果"
echo "   按 Ctrl+C 可以中断测试"
echo ""
echo "=== 开始测试 ==="
./weather-warning-detection

echo ""
echo "=== 部署完成 ==="
echo "程序已编译完成: ./weather-warning-detection"
echo "配置文件位置: ./config.json"
echo ""
echo "使用方法:"
echo "  ./weather-warning-detection                    # 使用默认配置文件"
echo "  ./weather-warning-detection -config custom.json  # 使用自定义配置文件"
echo ""
echo "定时任务设置示例 (每天8:00执行):"
echo "  0 8 * * * $(pwd)/weather-warning-detection -config $(pwd)/config.json"
echo ""
echo "注意事项:"
echo "1. 请在 config.json 中设置正确的高德API密钥"
echo "2. 如需飞书通知，请设置 feishu.webhook_url"
echo "3. 确保网络连接正常，能访问高德API和Moonshot AI"