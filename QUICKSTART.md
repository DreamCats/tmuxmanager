# 快速开始

## 1. 编译和安装

```bash
# 方式1: 使用安装脚本（推荐）
./install.sh

# 方式2: 手动安装
go build -o tmx ./cmd/tmx
sudo mv tmx /usr/local/bin/
```

## 2. 配置 tmux

```bash
# 启动 tmux（如果未运行）
tmux

# 在 tmux 中安装配置
tmx --install

# 重新加载 tmux 配置
tmux source-file ~/.tmux.conf
```

### 卸载配置

如果不再使用 tmx，可以卸载配置：

```bash
tmx --uninstall

# 然后重新加载配置
tmux source-file ~/.tmux.conf
```

## 3. 开始使用

### 第一次使用

如果你还没有运行 tmux，直接运行：

```bash
tmx
```

tmx 会提示是否自动启动 tmux：
```
📝 tmux 未运行

💡 tmx 可以自动启动 tmux 并创建默认会话
是否自动启动? [Y/n]:
```

按 `Enter` 或输入 `Y`，tmx 会：
1. 自动启动 tmux
2. 创建名为 `default` 的会话
3. 在新窗口中打开会话管理器

### 后续使用

现在你可以按 `Ctrl+b t` 随时打开会话管理器了！

### 创建测试会话

```bash
# 创建几个测试会话
tmux new-session -d -s dev
tmux new-session -d -s test
tmux new-session -d -s prod

# 打开管理器
tmx
```

### 常用操作

- 在管理器中选择会话，按 `Enter` 进入
- 按 `n` 创建新会话
- 按 `d` 断开会话
- 按 `x` 删除会话
- 按 `q` 退出管理器

## 故障排除

### 问题: 提示 "tmux 未运行"

**好消息**: tmx 会询问是否自动启动 tmux！

**手动启动**:
```bash
tmux
```

### 问题: 快捷键不生效

**解决方案 1**: 检查 tmx 是否在 PATH 中
```bash
which tmx
# 应该显示: /usr/local/bin/tmx 或其他路径

# 如果没有显示，需要安装到 PATH：
sudo mv tmx /usr/local/bin/
```

**解决方案 2**: 重新加载 tmux 配置
```bash
tmux source-file ~/.tmux.conf
```

**解决方案 3**: 检查快捷键是否绑定
```bash
tmux list-keys | grep tmx
# 应该显示: bind-key t run-shell "tmx"
```

**解决方案 4**: 使用诊断脚本
```bash
./diagnose.sh
```

### 问题: 看不到状态栏提示

**解决方案**: 检查 `~/.tmux.conf` 文件，确保包含以下内容：
```bash
set -g status-right '#[fg=green][Ctrl+B T] 管理器#[default] | %H:%M %Y-%m-%d'
```

### 问题: Ctrl+b t 没反应

**可能原因**:
1. tmx 不在 PATH 中 - 运行 `which tmx` 检查
2. tmux 配置未加载 - 运行 `tmux source-file ~/.tmux.conf`
3. 快捷键冲突 - 检查是否有其他配置覆盖了 `t` 键

**调试方法**:
```bash
# 在 tmux 中手动运行 tmx
tmx

# 如果手动运行可以，说明是快捷键绑定问题
# 检查绑定
tmux list-keys | grep "bind-key t"
```

## 下一步

查看 [README.md](./README.md) 了解更多功能
