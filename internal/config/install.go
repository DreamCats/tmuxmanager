package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const tmuxConfigMarker = "# ========== tmx 配置 =========="

const tmuxConfigContent = `
# ========== tmx 配置 ==========
# 按 Ctrl+b t 打开会话管理器
bind-key t run-shell "tmx"

# 在状态栏显示快捷键提示
set -g status-right '#[fg=green][Ctrl+B T] 管理器#[default] | %H:%M %Y-%m-%d'
# ========== tmx 配置结束 ==========
`

const shellConfigMarker = "# ========== tmx quit 命令 =========="

const bashQuitConfig = `
# ========== tmx quit 命令 ==========
# 退出 tmux 会话但保持运行（相当于 Ctrl+b d）
quit() {
    if [ -n "$TMUX" ]; then
        tmux detach-client
    else
        echo "不在 tmux 会话中"
    fi
}
# ========== tmx quit 命令结束 ==========
`

const zshQuitConfig = `
# ========== tmx quit 命令 ==========
# 退出 tmux 会话但保持运行（相当于 Ctrl+b d）
quit() {
    if [ -n "$TMUX" ]; then
        tmux detach-client
    else
        echo "不在 tmux 会话中"
    fi
}
# ========== tmx quit 命令结束 ==========
`

// InstallConfig 安装 tmux 配置
func InstallConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户目录: %w", err)
	}

	configPath := filepath.Join(homeDir, ".tmux.conf")

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 配置文件不存在，创建新文件
		if err := os.WriteFile(configPath, []byte(tmuxConfigContent), 0644); err != nil {
			return fmt.Errorf("无法创建配置文件: %w", err)
		}
		fmt.Printf("✓ 已创建配置文件: %s\n", configPath)
	} else {
		// 配置文件已存在，检查是否已包含 tmx 配置
		content, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("无法读取配置文件: %w", err)
		}

		// 检查是否已包含配置
		if contains(string(content), "tmx") {
			fmt.Println("⚠ tmx 配置已存在，无需重复安装")
		} else {
			// 追加配置
			file, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("无法打开配置文件: %w", err)
			}
			defer file.Close()

			if _, err := file.WriteString(tmuxConfigContent); err != nil {
				return fmt.Errorf("无法写入配置: %w", err)
			}

			fmt.Printf("✓ 已添加配置到: %s\n", configPath)
		}
	}

	// 安装 quit 命令到 shell 配置
	if err := installQuitCommand(homeDir); err != nil {
		fmt.Printf("⚠️  安装 quit 命令失败: %v\n", err)
	}

	fmt.Println("\n请运行以下命令重新加载 tmux 配置：")
	fmt.Println("  tmux source-file ~/.tmux.conf")
	fmt.Println("\n或重启 tmux")

	return nil
}

// installQuitCommand 安装 quit 命令到 shell 配置
func installQuitCommand(homeDir string) error {
	// 检测使用的 shell
	shellConfigPath := ""
	quitConfig := ""

	// 优先检查 zsh
	zshrc := filepath.Join(homeDir, ".zshrc")
	if _, err := os.Stat(zshrc); err == nil {
		shellConfigPath = zshrc
		quitConfig = zshQuitConfig
	}

	// 如果 zsh 不存在，检查 bash
	if shellConfigPath == "" {
		bashrc := filepath.Join(homeDir, ".bashrc")
		if _, err := os.Stat(bashrc); err == nil {
			shellConfigPath = bashrc
			quitConfig = bashQuitConfig
		}
	}

	// 如果都不存在，创建 .bashrc
	if shellConfigPath == "" {
		shellConfigPath = filepath.Join(homeDir, ".bashrc")
		quitConfig = bashQuitConfig
	}

	// 读取现有配置
	content, err := os.ReadFile(shellConfigPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("无法读取 shell 配置文件: %w", err)
	}

	// 检查是否已包含 quit 命令
	if content != nil && contains(string(content), shellConfigMarker) {
		fmt.Println("✓ quit 命令已安装")
		return nil
	}

	// 追加 quit 命令
	file, err := os.OpenFile(shellConfigPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("无法打开 shell 配置文件: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(quitConfig); err != nil {
		return fmt.Errorf("无法写入 quit 命令: %w", err)
	}

	fmt.Printf("✓ 已添加 quit 命令到: %s\n", shellConfigPath)
	return nil
}

// UninstallConfig 卸载 tmux 配置
func UninstallConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户目录: %w", err)
	}

	configPath := filepath.Join(homeDir, ".tmux.conf")

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("⚠ tmux 配置文件不存在，无需卸载")
	} else {
		// 读取配置文件
		content, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("无法读取配置文件: %w", err)
		}

		// 检查是否包含 tmx 配置
		if !contains(string(content), tmuxConfigMarker) {
			fmt.Println("⚠ 未找到 tmx 配置，无需卸载")
		} else {
			// 移除 tmx 配置块
			lines := strings.Split(string(content), "\n")
			newLines := []string{}
			inTmxBlock := false

			for _, line := range lines {
				if strings.Contains(line, tmuxConfigMarker) {
					inTmxBlock = !inTmxBlock
					continue
				}
				if !inTmxBlock {
					newLines = append(newLines, line)
				}
			}

			// 移除多余的空行
			newContent := strings.Join(newLines, "\n")
			newContent = strings.TrimSpace(newContent) + "\n"

			// 写回文件
			if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
				return fmt.Errorf("无法写入配置文件: %w", err)
			}

			fmt.Printf("✓ 已从 %s 移除 tmx 配置\n", configPath)
		}
	}

	// 卸载 quit 命令
	if err := uninstallQuitCommand(homeDir); err != nil {
		fmt.Printf("⚠️  卸载 quit 命令失败: %v\n", err)
	}

	fmt.Println("\n请运行以下命令重新加载 tmux 配置：")
	fmt.Println("  tmux source-file ~/.tmux.conf")
	fmt.Println("\n或重启 tmux")
	fmt.Println("\n并重新加载 shell 配置：")
	fmt.Println("  source ~/.bashrc  # 或 ~/.zshrc")

	return nil
}

// uninstallQuitCommand 卸载 quit 命令
func uninstallQuitCommand(homeDir string) error {
	shellConfigs := []string{
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".bashrc"),
	}

	removed := false
	for _, configPath := range shellConfigs {
		content, err := os.ReadFile(configPath)
		if err != nil {
			continue // 文件不存在，跳过
		}

		// 检查是否包含 quit 命令
		if !contains(string(content), shellConfigMarker) {
			continue
		}

		// 移除 quit 命令块
		lines := strings.Split(string(content), "\n")
		newLines := []string{}
		inQuitBlock := false

		for _, line := range lines {
			if strings.Contains(line, shellConfigMarker) {
				inQuitBlock = !inQuitBlock
				continue
			}
			if !inQuitBlock {
				newLines = append(newLines, line)
			}
		}

		// 移除多余的空行
		newContent := strings.Join(newLines, "\n")
		newContent = strings.TrimSpace(newContent) + "\n"

		// 写回文件
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			return fmt.Errorf("无法写入配置文件: %w", err)
		}

		fmt.Printf("✓ 已从 %s 移除 quit 命令\n", configPath)
		removed = true
	}

	if !removed {
		fmt.Println("⚠ 未找到 quit 命令配置")
	}

	return nil
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
