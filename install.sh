#!/bin/bash

# tmx 安装脚本

set -e

echo "🚀 安装 tmx..."

# 编译
echo "📦 编译中..."
go build -o tmx ./cmd/tmx

# 安装到 /usr/local/bin
if [ "$EUID" -ne 0 ]; then
    echo "⚠️  需要 sudo 权限来安装到 /usr/local/bin"
    sudo mv tmx /usr/local/bin/
else
    mv tmx /usr/local/bin/
fi

echo "✅ 安装完成！"
echo ""
echo "接下来："
echo "1. 如果 tmux 未运行，请先启动: tmux"
echo "2. 在 tmux 中运行: tmx --install"
echo "3. 重新加载配置: tmux source-file ~/.tmux.conf"
echo ""
echo "然后按 Ctrl+b t 就能打开会话管理器了！"
