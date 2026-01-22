package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/DreamCats/tmuxmanager/internal/tmux"
)

// Styles 定义 UI 样式
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#86AAEC")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("#EEEDFF")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true)

	activeIndicator = "▶ "

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(0, 1)
)

// Model 是 TUI 的状态模型
type Model struct {
	sessions         []tmux.Session
	selected         int
	manager          *tmux.Manager
	quitting         bool
	width            int
	height           int
	inputMode        bool
	inputBuffer      string
	newSessionName   string    // 新创建的会话名称
	attachSessionName string    // 要附加的会话名称
}

// Messages
type sessionsLoadedMsg []tmux.Session
type sessionAttachedMsg struct {
	err  error
	name string
}
type sessionDetachedMsg struct{ err error }
type sessionCreatedMsg struct {
	name string
	err  error
}
type sessionKilledMsg struct{ err error }

// Init 初始化 TUI
func (m Model) Init() tea.Cmd {
	return m.loadSessions()
}

// Update 处理事件
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 输入模式下处理
		if m.inputMode {
			return m.handleInput(msg)
		}

		// 正常模式
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}

		case "down", "j":
			if m.selected < len(m.sessions)-1 {
				m.selected++
			}

		case "enter":
			return m, m.attachSession()

		case "n":
			m.inputMode = true
			m.inputBuffer = ""
			return m, nil

		case "d":
			return m, m.detachSession()

		case "x":
			return m, m.killSession()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case sessionsLoadedMsg:
		m.sessions = msg
		// 如果刚创建了新会话，选中它
		if m.newSessionName != "" {
			for i, session := range m.sessions {
				if session.Name == m.newSessionName {
					m.selected = i
					m.newSessionName = "" // 清空标记
					break
				}
			}
		}
		// 确保选中项有效
		if m.selected >= len(m.sessions) {
			m.selected = len(m.sessions) - 1
		}
		return m, nil

	case sessionAttachedMsg:
		if msg.err != nil {
			// 显示错误
			m.quitting = true
			return m, tea.Quit
		}
		// 保存要附加的会话名，然后退出 TUI
		m.attachSessionName = msg.name
		m.quitting = true
		return m, tea.Quit

	case sessionDetachedMsg:
		return m, m.loadSessions()

	case sessionCreatedMsg:
		if msg.err != nil {
			// 创建失败，显示错误并返回列表
			fmt.Printf("\n创建会话失败: %v\n", msg.err)
			return m, tea.Quit
		}
		// 创建成功，保存会话名并刷新列表
		m.newSessionName = msg.name
		return m, m.loadSessions()

	case sessionKilledMsg:
		return m, m.loadSessions()
	}

	return m, nil
}

// handleInput 处理输入模式
func (m Model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.inputBuffer != "" {
			m.inputMode = false
			return m, m.createSession(m.inputBuffer)
		}
		m.inputMode = false
		return m, nil

	case "esc":
		m.inputMode = false
		m.inputBuffer = ""
		return m, nil

	case "ctrl+h", "backspace":
		if len(m.inputBuffer) > 0 {
			m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
		}

	default:
		// 添加字符到缓冲区
		if len(msg.String()) == 1 {
			m.inputBuffer += msg.String()
		}
	}

	return m, nil
}

// View 渲染 UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// 输入模式
	if m.inputMode {
		return m.renderInput()
	}

	// 正常模式
	return m.renderNormal()
}

// renderNormal 渲染正常模式界面
func (m Model) renderNormal() string {
	var b strings.Builder

	// 标题
	title := titleStyle.Render("Tmux 会话管理")
	b.WriteString(title)
	b.WriteString("\n\n")

	// 会话列表
	if len(m.sessions) == 0 {
		b.WriteString(itemStyle.Render("没有会话，按 n 新建会话"))
		b.WriteString("\n")
	} else {
		for i, session := range m.sessions {
			style := itemStyle
			if i == m.selected {
				style = selectedStyle
			}

			// 构建会话信息
			indicator := "  "
			if session.Attached {
				indicator = activeIndicator
			}

			timeInfo := formatTime(session.Created)

			line := fmt.Sprintf("%s%s%s (%s)",
				indicator,
				session.Name,
				strings.Repeat(" ", 40-len(session.Name)),
				timeInfo,
			)

			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	// 快捷键提示
	hints := "[Enter]进入 [d]断开 [n]新建 [x]删除 [q]退出"
	b.WriteString(hintStyle.Render(hints))
	b.WriteString("\n")

	// 额外提示：如何退出 tmux 会话
	tip := "💡 提示：进入会话后按 Ctrl+b d 可退出但保持会话运行"
	b.WriteString(hintStyle.Render(tip))

	return b.String()
}

// renderInput 渲染输入模式界面
func (m Model) renderInput() string {
	var b strings.Builder

	// 标题
	title := titleStyle.Render("新建会话")
	b.WriteString(title)
	b.WriteString("\n\n")

	// 输入提示
	b.WriteString(itemStyle.Render("请输入会话名称:"))
	b.WriteString("\n\n")

	// 输入框
	inputStyle := selectedStyle
	inputLine := "> " + m.inputBuffer + "_"
	b.WriteString(inputStyle.Render(inputLine))
	b.WriteString("\n\n")

	// 快捷键提示
	hints := "[Enter]确认 [Esc]取消"
	b.WriteString(hintStyle.Render(hints))

	return b.String()
}

// Commands

func (m Model) loadSessions() tea.Cmd {
	return func() tea.Msg {
		sessions, err := m.manager.ListSessions()
		if err != nil {
			return sessionsLoadedMsg{}
		}
		return sessionsLoadedMsg(sessions)
	}
}

func (m Model) attachSession() tea.Cmd {
	return func() tea.Msg {
		if m.selected >= len(m.sessions) {
			return sessionAttachedMsg{err: fmt.Errorf("no session selected")}
		}
		session := m.sessions[m.selected]

		// 注意：我们不在这里直接调用 AttachSession
		// 因为它会阻塞并接管终端
		// 我们只返回会话名，让 main 函数处理
		return sessionAttachedMsg{
			err:  nil,
			name: session.Name,
		}
	}
}

func (m Model) detachSession() tea.Cmd {
	return func() tea.Msg {
		if m.selected >= len(m.sessions) {
			return sessionDetachedMsg{nil}
		}
		session := m.sessions[m.selected]
		err := m.manager.DetachSession(session.Name)
		return sessionDetachedMsg{err}
	}
}

func (m Model) createSession(name string) tea.Cmd {
	return func() tea.Msg {
		// 只创建会话，不自动进入
		err := m.manager.NewSession(name)
		return sessionCreatedMsg{name: name, err: err}
	}
}

func (m Model) killSession() tea.Cmd {
	return func() tea.Msg {
		if m.selected >= len(m.sessions) {
			return sessionKilledMsg{nil}
		}
		session := m.sessions[m.selected]
		err := m.manager.KillSession(session.Name)
		return sessionKilledMsg{err}
	}
}

// AttachSessionName 返回要附加的会话名称
func (m Model) AttachSessionName() string {
	return m.attachSessionName
}

// NewModel 创建新的 Model
func NewModel() Model {
	return Model{
		sessions: make([]tmux.Session, 0),
		selected: 0,
		manager:  tmux.NewManager(),
		quitting: false,
		inputMode: false,
		inputBuffer: "",
	}
}

// formatTime 格式化时间显示
func formatTime(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "刚刚"
	} else if duration < time.Hour {
		return fmt.Sprintf("%d分钟前", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d小时前", int(duration.Hours()))
	} else if duration < 30*24*time.Hour {
		return fmt.Sprintf("%d天前", int(duration.Hours()/24))
	}
	return t.Format("2006-01-02")
}
