# tmux 快捷键速查

## 🔑 最常用的快捷键

### 退出会话（不杀死）

```
Ctrl+b d
```

**这是什么**：Detach（分离）- 从会话中退出，但会话继续在后台运行

**什么时候用**：
- ✅ 想离开会话但保持其运行
- ✅ 想切换到其他会话
- ✅ 想关闭终端但保留工作环境

**对比**：
- ❌ `exit` - 会杀死会话，所有工作丢失
- ✅ `Ctrl+b d` - 只分离，会话继续运行

## 📋 常用操作

| 操作 | 快捷键 | 说明 |
|------|--------|------|
| **退出会话** | `Ctrl+b d` | 分离会话（保持运行） |
| **打开管理器** | `Ctrl+b t` | 打开 tmx 会话管理器 |
| **创建新窗口** | `Ctrl+b c` | 在当前会话中创建新窗口 |
| **切换窗口** | `Ctrl+b 0-9` | 切换到指定编号的窗口 |
| **列出窗口** | `Ctrl+b w` | 显示窗口列表 |
| **重命名窗口** | `Ctrl+b ,` | 重命名当前窗口 |
| **垂直分屏** | `Ctrl+b %` | 上下分屏 |
| **水平分屏** | `Ctrl+b "` | 左右分屏 |
| **切换面板** | `Ctrl+b 方向键` | 在分屏间切换 |

## 🎯 典型使用流程

### 场景 1：启动工作环境

```bash
# 1. 打开 tmx
./tmx

# 2. 创建会话
按 n → 输入 "work" → Enter

# 3. 进入会话
选中 "work" → Enter

# 4. 在会话中工作
运行你的程序...

# 5. 暂时离开（保持会话运行）
Ctrl+b d  ← 重要！

# 6. 下次继续工作
./tmx 或 Ctrl+b t → 选中 "work" → Enter
```

### 场景 2：多个项目

```bash
# 创建多个会话
./tmx → n → "project1" → Enter
./tmx → n → "project2" → Enter
./tmx → n → "project3" → Enter

# 快速切换
Ctrl+b t → 选择项目 → Enter

# 离开时
Ctrl+b d  ← 保持所有项目运行
```

## ⚠️ 常见错误

### 错误 1：用 exit 退出

```bash
# ❌ 错误做法
$ exit
# 结果：会话被杀死，所有程序关闭

# ✅ 正确做法
Ctrl+b d
# 结果：会话保持运行，下次可以继续
```

### 错误 2：忘记快捷键

```bash
# 💡 记住这两个就够了：
Ctrl+b t  # 打开会话管理器
Ctrl+b d  # 退出会话（保持运行）
```

## 🎓 记忆技巧

### Ctrl+b 是什么意思？

`Ctrl+b` 是 tmux 的**前缀键**（prefix key），所有 tmux 快捷键都要先按它。

### 为什么是 d？

**d** = **d**etach（分离）

### 怎么记 Ctrl+b d？

```
Ctrl+b   = 唤醒 tmux
d        = detach（分离）
```

## 💡 进阶技巧

### 永久保存会话

结合 tmux 的持久化功能，即使重启服务器，会话也能恢复：

```bash
# 在 ~/.tmux.conf 中添加
set -g @plugin 'tmux-plugins/tpm'
set -g @plugin 'tmux-plugins/tmux-resurrect'
```

### 自动启动会话

```bash
# 在 ~/.tmux.conf 中添加
new-session -n work  # 自动创建 work 会话
```

## 📚 更多资源

- tmux 官方文档：`man tmux`
- tmx 快捷键：按 `Ctrl+b t` 打开管理器
- 列出所有快捷键：`Ctrl+b ?`

## 🆘 遇到问题？

```bash
# 查看所有会话
tmux list-sessions

# 杀死指定会话
tmux kill-session -t session-name

# 杀死所有会话
tmux kill-server

# 查看实时日志
tmux show-options -g
```
