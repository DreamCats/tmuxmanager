package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/DreamCats/tmuxmanager/internal/config"
	"github.com/DreamCats/tmuxmanager/internal/tmux"
	"github.com/DreamCats/tmuxmanager/internal/ui"
)

func main() {
	// 先检查命令行参数（不需要 tmux 运行）
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		case "-v", "--version":
			printVersion()
			os.Exit(0)
		case "--install":
			installConfig()
			os.Exit(0)
		case "--uninstall":
			uninstallConfig()
			os.Exit(0)
		default:
			fmt.Printf("未知参数: %s\n", os.Args[1])
			fmt.Println("使用 -h 查看帮助")
			os.Exit(1)
		}
	}

	// 检查是否在 tmux 会话中
	if os.Getenv("TMUX") == "" {
		fmt.Println("📝 tmx 需要在 tmux 会话中运行")
		fmt.Println("\n💡 使用方法：")
		fmt.Println("   tmux                          # 启动 tmux")
		fmt.Println("   tmx                           # 在 tmux 中运行管理器")
		fmt.Println("\n或者：")
		fmt.Println("   tmux attach-session -t default  # 连接到现有会话")
		fmt.Println("   tmx                           # 然后运行 tmx")
		fmt.Println("\n💡 提示：运行 ./tmx --install 可配置 Ctrl+b t 快捷键")
		os.Exit(1)
	}

	// 检查 tmux 是否运行
	manager := tmux.NewManager()
	if !manager.IsTmuxRunning() {
		// tmux 未运行，询问是否自动启动
		fmt.Println("📝 tmux 未运行")
		fmt.Println("\n💡 tmx 可以自动启动 tmux 并创建默认会话")
		fmt.Print("是否自动启动? [Y/n]: ")

		var answer string
		fmt.Scanln(&answer)

		// 默认是 Y，或者用户输入 y/Y
		if answer == "" || answer == "y" || answer == "Y" {
			fmt.Println("\n🚀 正在启动 tmux...")

			// 检查是否已有 default 会话
			sessions, _ := manager.ListSessions()
			hasDefault := false
			for _, s := range sessions {
				if s.Name == "default" {
					hasDefault = true
					break
				}
			}

			if hasDefault {
				// default 会话已存在，直接附加
				fmt.Println("✓ 找到现有会话 'default'，正在连接...")
				attachCmd := exec.Command("tmux", "attach-session", "-t", "default")
				attachCmd.Stdin = os.Stdin
				attachCmd.Stdout = os.Stdout
				attachCmd.Stderr = os.Stderr

				if err := attachCmd.Run(); err != nil {
					fmt.Printf("❌ 附加到 tmux 会话失败: %v\n", err)
					os.Exit(1)
				}
				return
			} else {
				// 创建新会话
				createCmd := exec.Command("tmux", "new-session", "-d", "-s", "default")
				if err := createCmd.Run(); err != nil {
					fmt.Printf("❌ 创建 tmux 会话失败: %v\n", err)
					fmt.Println("\n你可以手动启动 tmux：")
					fmt.Println("  tmux")
					os.Exit(1)
				}

				// 设置 tmux 在附加后运行 tmx
				execCmd := exec.Command("tmux", "send-keys", "-t", "default", "tmx", "C-m")
				if err := execCmd.Run(); err != nil {
					fmt.Printf("⚠️  警告: 无法自动启动 tmx: %v\n", err)
				}

				// 附加到会话
				attachCmd := exec.Command("tmux", "attach-session", "-t", "default")
				attachCmd.Stdin = os.Stdin
				attachCmd.Stdout = os.Stdout
				attachCmd.Stderr = os.Stderr

				if err := attachCmd.Run(); err != nil {
					fmt.Printf("❌ 附加到 tmux 会话失败: %v\n", err)
					os.Exit(1)
				}
				return
			}
		} else {
			// 用户选择不自动启动
			fmt.Println("\n请先启动 tmux：")
			fmt.Println("  tmux")
			fmt.Println("\n或者创建新会话：")
			fmt.Println("  tmux new")
			os.Exit(1)
		}
	}

	// 启动 TUI
	model := ui.NewModel()
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // 使用备用屏幕
		tea.WithMouseCellMotion(), // 启用鼠标支持
	)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 检查是否需要附加到会话
	if m, ok := finalModel.(ui.Model); ok && m.AttachSessionName() != "" {
		if err := manager.AttachSession(m.AttachSessionName()); err != nil {
			fmt.Printf("错误: 无法连接到会话: %v\n", err)
			os.Exit(1)
		}
	}
}

func printHelp() {
	fmt.Println("tmx - Tmux 会话管理器")
	fmt.Println("\n用法:")
	fmt.Println("  tmx                打开会话管理器（TUI）")
	fmt.Println("  tmx --install      安装 tmux 配置（快捷键 + 状态栏提示）")
	fmt.Println("  tmx --uninstall    卸载 tmux 配置")
	fmt.Println("  tmx -h             显示帮助")
	fmt.Println("  tmx -v             显示版本")
	fmt.Println("\n注意: tmx 需要在 tmux 会话中运行")
	fmt.Println("\n💡 使用方法：")
	fmt.Println("   tmux              # 启动 tmux")
	fmt.Println("   tmx               # 在 tmux 中运行管理器")
	fmt.Println("   或运行 ./tmx --install 配置 Ctrl+b t 快捷键")
	fmt.Println("\nTUI 快捷键:")
	fmt.Println("  Enter           进入选中的会话")
	fmt.Println("  n               新建会话")
	fmt.Println("  d               断开会话")
	fmt.Println("  x               删除会话")
	fmt.Println("  ↑/↓ 或 j/k      导航")
	fmt.Println("  q/Esc           退出")
	fmt.Println("\n退出 tmux 会话:")
	fmt.Println("  Ctrl+b d        分离会话（保持运行）")
	fmt.Println("  quit            分离会话（需要先运行 ./tmx --install）")
}

func printVersion() {
	fmt.Println("tmx version 1.0.0")
}

func installConfig() {
	if err := config.InstallConfig(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}

func uninstallConfig() {
	if err := config.UninstallConfig(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}
