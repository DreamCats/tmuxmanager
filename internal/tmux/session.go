package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Session 表示一个 tmux 会话
type Session struct {
	Name      string
	Created   time.Time
	Active    bool
	Windows   int
	Attached  bool
}

// Manager 管理 tmux 会话
type Manager struct{}

// NewManager 创建一个新的 Manager
func NewManager() *Manager {
	return &Manager{}
}

// ListSessions 获取所有 tmux 会话
func (m *Manager) ListSessions() ([]Session, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}:#{session_created}:#{session_windows}:#{session_attached}")
	output, err := cmd.Output()
	if err != nil {
		// 如果 tmux 没有运行或没有会话
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return []Session{}, nil
			}
		}
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	sessions := make([]Session, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			continue
		}

		name := parts[0]
		createdTimestamp := strings.TrimPrefix(parts[1], ";") // tmux 时间戳格式
		windows := parts[2]
		attached := parts[3]

		// 解析时间戳（tmux 使用 Unix 时间戳，单位是微秒或毫秒）
		var created time.Time
		if ts, err := parseTimestamp(createdTimestamp); err == nil {
			created = ts
		} else {
			created = time.Now() // 如果解析失败，使用当前时间
		}

		sessions = append(sessions, Session{
			Name:     name,
			Created:  created,
			Windows:  parseInt(windows),
			Attached: parseInt(attached) > 0,
		})
	}

	return sessions, nil
}

// AttachSession 连接到指定的会话
func (m *Manager) AttachSession(name string) error {
	// 检查当前是否在 tmux 会话中
	if inTmuxSession() {
		// 在 tmux 中，使用 switch-client
		cmd := exec.Command("tmux", "switch-client", "-t", name)
		return cmd.Run()
	}

	// 不在 tmux 中，使用 attach-session
	cmd := exec.Command("tmux", "attach-session", "-t", name)
	return cmd.Run()
}

// DetachSession 断开指定的会话
func (m *Manager) DetachSession(name string) error {
	cmd := exec.Command("tmux", "detach-session", "-t", name)
	return cmd.Run()
}

// NewSession 创建新会话
func (m *Manager) NewSession(name string) error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", name)
	return cmd.Run()
}

// NewSessionAndAttach 创建新会话并立即进入
func (m *Manager) NewSessionAndAttach(name string) error {
	cmd := exec.Command("tmux", "new-session", "-s", name)
	return cmd.Run()
}

// KillSession 删除指定的会话
func (m *Manager) KillSession(name string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", name)
	return cmd.Run()
}

// IsTmuxRunning 检查 tmux 是否在运行
func (m *Manager) IsTmuxRunning() bool {
	// 检查是否有 tmux 会话存在
	sessions, err := m.ListSessions()
	if err != nil {
		return false
	}
	return len(sessions) > 0
}

// inTmuxSession 检查当前是否在 tmux 会话中
func inTmuxSession() bool {
	// 检查 TMUX 环境变量
	return os.Getenv("TMUX") != ""
}

// 辅助函数

func parseTimestamp(ts string) (time.Time, error) {
	// tmux 时间戳可能是秒或微秒
	var sec int64
	if len(ts) > 10 {
		_, err := fmt.Sscanf(ts, "%d", &sec)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(sec/1000000, 0), nil
	}
	_, err := fmt.Sscanf(ts, "%d", &sec)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
