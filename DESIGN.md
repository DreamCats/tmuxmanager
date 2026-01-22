# tmx - Tmux 会话管理器设计文档

## 项目简介

**tmx** 是一个简单易用的 tmux 会话管理 TUI 工具，帮助用户不记住复杂的 tmux 快捷键也能轻松管理会话。

## 核心设计原则

1. **零配置** - 开箱即用，不需要修改 tmux 配置也能工作
2. **快捷键可视化** - 底部永久显示快捷键提示，不需要记忆
3. **简单直观** - 所有操作都有明确的提示和反馈
4. **单一入口** - 只需记住 `Ctrl+b t`，其他都在界面上

## 功能需求

### 核心功能

1. **会话列表展示**
   - 显示所有 tmux 会话
   - 标记当前活跃会话
   - 显示会话创建时间/最后活跃时间

2. **会话操作**
   - 进入会话（attach）
   - 断开会话（detach）
   - 新建会话（新建并提示输入名称）
   - 删除会话

3. **TUI 界面**
   - 简洁的列表展示
   - 底部快捷键提示栏
   - 支持方向键导航
   - 支持vim键位（j/k）

### 快捷键设计

| 快捷键 | 功能 | 说明 |
|--------|------|------|
| `↑`/`↓` 或 `k`/`j` | 导航 | 上下移动选择 |
| `Enter` | 进入会话 | 连接到选中的会话 |
| `n` | 新建会话 | 提示输入会话名称 |
| `d` | 断开会话 | 分离选中的会话 |
| `Del` | 删除会话 | 永久删除选中的会话 |
| `q` / `Esc` | 退出 | 关闭管理器 |

## 技术选型

### 编程语言
- **Go** - 本项目就是 Go 项目，使用 Go 便于维护

### TUI 库选择
推荐使用 **bubbletea** + **lipogloss**：
- `bubbletea` - 强大的 TUI 框架（基于 Elm 架构）
- `lipogloss` - 样式库，用于美化界面
- `bubbletea` 是目前 Go 生态中最流行的 TUI 框架

### tmux 交互
- 通过执行 shell 命令与 tmux 交互
- 主要命令：
  - `tmux list-sessions` - 列出会话
  - `tmux attach-session -t <name>` - 进入会话
  - `tmux detach-session -t <name>` - 断开会话
  - `tmux new-session -d -s <name>` - 新建会话
  - `tmux kill-session -t <name>` - 删除会话

## 项目结构

```
tmuxmanager/
├── cmd/
│   └── tmx/
│       └── main.go           # 入口文件
├── internal/
│   ├── tmux/
│   │   ├── session.go        # 会话管理
│   │   └── command.go        # tmux 命令封装
│   ├── ui/
│   │   ├── tui.go            # TUI 主逻辑
│   │   ├── styles.go         # 样式定义
│   │   └── components.go     # UI 组件
│   └── config/
│       └── install.go        # 配置安装脚本
├── DESIGN.md                 # 本文档
├── README.md                 # 使用说明
├── go.mod
└── go.sum
```

## TUI 界面设计

```
┌─ Tmux 会话管理 ──────────────────────┐
│                                      │
│  ▶ dev-server      (活跃)            │
│    backend-api     (2小时前)          │
│    frontend        (昨天)             │
│    test-env                           │
│                                      │
│ [Enter]进入 [d]断开 [n]新建 [q]退出  │
└──────────────────────────────────────┘
```

- **顶部**：标题栏
- **中间**：会话列表（当前选中高亮）
- **底部**：快捷键提示栏

## tmux 集成方案

### 方式1：手动运行
```bash
tmx
```

### 方式2：tmux 快捷键（推荐）
在 `~/.tmux.conf` 中添加：
```bash
# 绑定 Ctrl+b t 打开会话管理器
bind-key t run-shell "tmx"

# 在状态栏显示提示
set -g status-right '#[fg=green][Ctrl+B T] 管理器#[default] | %H:%M %Y-%m-%d'
```

### 自动配置
提供 `tmx --install` 命令，自动修改 `~/.tmux.conf`：
- 添加快捷键绑定
- 添加状态栏提示
- 保留原有配置，只追加内容

## 使用流程示例

### 场景1：切换会话
```bash
1. 按 Ctrl+b t（打开管理器）
2. 选择目标会话
3. 按 Enter（自动切换）
```

### 场景2：新建会话
```bash
1. 按 Ctrl+b t
2. 按 n（新建）
3. 输入会话名称
4. 自动进入新会话
```

### 场景3：断开会话
```bash
1. 按 Ctrl+b t
2. 选中要断开的会话
3. 按 d（断开）
```

## 实现优先级

### Phase 1: MVP（最小可用版本）
- [x] 设计文档
- [ ] Go 项目初始化
- [ ] tmux 会话列表获取
- [ ] 基础 TUI 界面
- [ ] 进入会话功能

### Phase 2: 核心功能
- [ ] 新建会话
- [ ] 断开会话
- [ ] 删除会话
- [ ] 快捷键完整实现

### Phase 3: 优化
- [ ] 会话时间显示
- [ ] 自动配置脚本（--install）
- [ ] 错误处理和提示
- [ ] 彩色样式优化

## 依赖库

```go
require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/lipgloss v0.7.0
)
```

## 注意事项

1. **跨平台兼容性** - 确保在 Linux/macOS 上都能正常工作
2. **错误处理** - 如果 tmux 未运行，给出友好提示
3. **权限问题** - 修改 `~/.tmux.conf` 需要用户确认
4. **备份配置** - 自动配置前备份原有配置文件

## 未来扩展

- [ ] 支持会话重命名
- [ ] 支持会话搜索/过滤
- [ ] 支持保存常用会话配置
- [ ] 支持会话分组（window/pane 管理）
- [ ] 支持主题切换

## 参考资料

- [bubbletea 官方文档](https://github.com/charmbracelet/bubbletea)
- [lipgloss 官方文档](https://github.com/charmbracelet/lipgloss)
- [tmux 官方文档](https://github.com/tmux/tmux/wiki)
