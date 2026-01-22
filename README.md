# tmx - Tmux 会话管理器

一个简单易用的 tmux 会话管理 TUI 工具，让你不再需要记住复杂的 tmux 快捷键。

## 特性

- ✅ **零配置** - 开箱即用，不需要修改 tmux 配置
- ✅ **快捷键可视化** - 底部永久显示快捷键提示，不需要记忆
- ✅ **简单直观** - TUI 界面，所有操作都有明确提示
- ✅ **单一入口** - 只需记住 `Ctrl+b t`，其他都在界面上
- ✅ **quit 命令** - 自动安装 `quit` 命令，优雅退出 tmux 会话
- ✅ **Go 语言编写** - 单一二进制文件，方便部署

## 设计理念

**tmx 是 tmux 的内建工具**，类似浏览器扩展。

tmux 会话是容器，tmx 是其中的管理界面。这种设计：
- 避免终端状态问题
- 逻辑简单，不会出错
- 符合 tmux 插件生态

## 安装

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/DreamCats/tmuxmanager.git
cd tmuxmanager

# 编译
go build -o tmx ./cmd/tmx

# 安装到 PATH
sudo mv tmx /usr/local/bin/
```

### 配置 tmux 快捷键（推荐）

```bash
# 自动安装配置
tmx --install

# 重新加载 shell 配置（安装 quit 命令）
source ~/.zshrc  # 或 source ~/.bashrc

# 重新加载 tmux 配置
tmux source-file ~/.tmux.conf
```

这将：
- 绑定 `Ctrl+b t` 快捷键
- 在状态栏显示快捷键提示
- 安装 `quit` 命令到你的 shell

## 使用方法

### ⚠️ 重要提示

**tmx 必须在 tmux 会话内运行！**

### 基本用法

```bash
# 1. 启动 tmux
tmux

# 2. 在 tmux 中运行 tmx
tmx

# 3. 或者配置快捷键后
# 按 Ctrl+b t 直接打开管理器
```

### 命令行参数

| 参数 | 功能 | 使用位置 |
|------|------|----------|
| `tmx` | 打开管理器 | tmux 内 |
| `tmx --install` | 安装配置 | 任何地方 |
| `tmx --uninstall` | 卸载配置 | 任何地方 |
| `tmx -h` | 显示帮助 | 任何地方 |
| `tmx -v` | 显示版本 | 任何地方 |

### TUI 界面

```
┌─ Tmux 会话管理 ──────────────────────┐
│                                      │
│  ▶ dev-server      (活跃)            │
│    backend-api     (2小时前)          │
│    frontend        (昨天)             │
│    test-env                           │
│                                      │
│ [Enter]进入 [d]断开 [n]新建 [x]删除  │
│ 💡 提示：进入会话后按 Ctrl+b d 可    │
│          退出但保持会话运行          │
└──────────────────────────────────────┘
```

### TUI 快捷键

| 快捷键 | 功能 |
|--------|------|
| `↑` / `↓` 或 `k` / `j` | 导航会话列表 |
| `Enter` | 进入选中的会话 |
| `n` | 新建会话（会提示输入名称） |
| `d` | 断开选中的会话 |
| `x` | 删除选中的会话 |
| `q` / `Esc` | 退出管理器 |

### 退出 tmux 会话

**推荐方式**（保持会话运行）：
- `Ctrl+b d` - tmux 标准快捷键
- `quit` - tmx 提供的便捷命令（需先运行 `tmx --install`）

**不推荐**：
- `exit` - 会杀死会话，所有程序停止

## 典型使用流程

### 第一次使用

```bash
# 1. 编译并安装
go build -o tmx ./cmd/tmx
sudo mv tmx /usr/local/bin/

# 2. 配置 tmux
tmx --install
source ~/.zshrc              # 或 ~/.bashrc
tmux source-file ~/.tmux.conf

# 3. 启动 tmux
tmux

# 4. 在 tmux 中运行 tmx
tmx
```

### 日常使用

```bash
# 1. 进入 tmux
tmux                          # 或 tmux attach-session -t default

# 2. 按 Ctrl+b t 打开管理器
#    或直接输入: tmx

# 3. 选择会话，按 Enter 进入

# 4. 工作完成后，退出会话
quit                          # 或 Ctrl+b d
```

### 创建多个会话

```bash
# 在 tmux 中
tmx                           # 打开管理器
# 按 n 创建 "work" 会话
# 按 n 创建 "study" 会话
# 按 n 创建 "personal" 会话

# 快速切换
Ctrl+b t → 选择会话 → Enter

# 离开时
quit                          # 会话保持运行
```

## 项目结构

```
tmuxmanager/
├── cmd/
│   └── tmx/
│       └── main.go           # 入口文件
├── internal/
│   ├── tmux/
│   │   └── session.go        # tmux 会话管理
│   ├── ui/
│   │   └── tui.go            # TUI 界面
│   └── config/
│       └── install.go        # 配置安装脚本
├── DESIGN.md                 # 设计文档
├── README.md                 # 本文档
├── QUICKSTART.md             # 快速开始指南
├── TMUX_GUIDE.md             # tmux 快捷键指南
├── KEYBINDINGS.md            # 快捷键速查表
└── go.mod
```

## 技术栈

- **Go 1.24+**
- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI 框架
- [lipgloss](https://github.com/charmbracelet/lipgloss) - 样式库

## 为什么选择 tmx？

### 传统的 tmux 使用痛点

- ❌ 快捷键太多，记不住
- ❌ 查看会话列表需要记住 `Ctrl+b s`
- ❌ 切换、断开会话操作复杂
- ❌ 容易误用 `exit` 杀死会话
- ❌ 没有可视化的会话管理界面

### tmx 的解决方案

- ✅ 只需记住一个快捷键：`Ctrl+b t`
- ✅ 所有操作都有可视化提示
- ✅ 简单直观的列表界面
- ✅ 底部固定显示快捷键说明
- ✅ 提供 `quit` 命令，避免误杀会话

## 常见问题

### Q: 为什么 tmx 必须在 tmux 内运行？

A: tmx 是 tmux 的内建工具，类似浏览器扩展。在 tmux 内使用可以：
- 避免终端状态问题
- 提供一致的用户体验
- 符合 tmux 插件生态

### Q: 如何在 tmux 外快速启动 tmux 和 tmx？

A: 运行 `tmx --install` 配置快捷键后：
- 在 tmux 内按 `Ctrl+b t` 即可打开管理器
- 或手动：`tmux` → `tmx`

### Q: `exit` 和 `quit` 有什么区别？

A:
- `exit` - 杀死会话，所有程序停止
- `quit` - 分离会话，程序继续运行

## 开发

```bash
# 安装依赖
go mod download

# 编译
go build -o tmx ./cmd/tmx

# 运行测试
go test ./...
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 作者

买峰 <maifeng@bytedance.com>

## 许可证

MIT License

## 相关资源

- [tmux 官方文档](https://github.com/tmux/tmux/wiki)
- [DESIGN.md](./DESIGN.md) - 详细设计文档
- [QUICKSTART.md](./QUICKSTART.md) - 快速开始指南
- [TMUX_GUIDE.md](./TMUX_GUIDE.md) - tmux 快捷键完整指南
