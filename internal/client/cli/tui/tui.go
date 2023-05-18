package tui

import (
	"fmt"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedLogin = focusedStyle.Copy().Render("[ Login ]")
	blurredLogin = fmt.Sprintf("[ %s ]", blurredStyle.Render("Login"))

	focusedRegister = focusedStyle.Copy().Render("[ Register ]")
	blurredRegister = fmt.Sprintf("[ %s ]", blurredStyle.Render("Register"))
)

type usecase interface {
	GetLoginsPasswords() []domain.LoginPassword
	AddLoginPassword(domain.LoginPassword) error
	AddTextData(domain.TextData) error
	AddBinaryData(domain.BinaryData) error
	AddCardData(domain.CardData) error
	RegisterUser(user domain.User) error
	LoginUser(user domain.User) error
	CheckSync() (time.Time, error)
	GetLocalSyncTime() time.Time
	SyncData() error
	GetVersion() string
	GetBuildTime() string
}

type tui struct {
	focusIndex int

	inputsAuth []textinput.Model

	lpTable table.Model

	email string

	nonAuth     bool
	usecase     usecase
	offline     bool
	errorString string
}

func NewTUI(uc usecase) tui {
	m := tui{
		inputsAuth: make([]textinput.Model, 2),
		lpTable: table.New(
			table.WithHeight(7),
		),
	}
	m.nonAuth = true
	m.usecase = uc
	var t textinput.Model
	for i := range m.inputsAuth {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Email"
			t.Focus()
			t.CharLimit = 64
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputsAuth[i] = t
	}

	return m
}

func (m tui) Init() tea.Cmd {
	if rand.Int31n(1)%1 > 0 {
		m.offline = true
	}
	return textinput.Blink
}

func (m tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	m.errorString = ""
	if m.nonAuth {
		cmd = m.updateAuth(msg)
	} else {
		cmd = m.updateMain(msg)
	}

	// Handle character input and blinking
	tea.Batch(cmd, m.updateInputs(msg))
	return m, cmd
}

func (m *tui) updateAuth(msg tea.Msg) tea.Cmd {
	var err error
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.focusIndex == 2 { //login btn
				user := domain.User{
					Email:    m.inputsAuth[0].Value(),
					Password: m.inputsAuth[1].Value(),
				}
				err = m.usecase.LoginUser(user)
				if err != nil {
					m.errorString = err.Error()
				} else {
					m.email = user.Email
					m.getData()
					m.lpTable.Focus()
					m.nonAuth = false
				}

			}
			if s == "enter" && m.focusIndex == 3 {
				user := domain.User{
					Email:    m.inputsAuth[0].Value(),
					Password: m.inputsAuth[1].Value(),
				}
				err = m.usecase.RegisterUser(user)
				if err != nil {
					m.errorString = err.Error()
				} else {
					m.email = user.Email
					m.getData()
					m.nonAuth = false
				}

			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 3 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 3
			}

			cmds := make([]tea.Cmd, len(m.inputsAuth))
			for i := 0; i <= len(m.inputsAuth)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputsAuth[i].Focus()
					m.inputsAuth[i].PromptStyle = focusedStyle
					m.inputsAuth[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputsAuth[i].Blur()
				m.inputsAuth[i].PromptStyle = noStyle
				m.inputsAuth[i].TextStyle = noStyle
			}

			return tea.Batch(cmds...)
		}
	}
	return nil
}

func (m *tui) updateMain(msg tea.Msg) tea.Cmd {
	//var err error
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.lpTable.MoveDown(1)
		case "up":
			m.lpTable.MoveUp(1)
		}
	}
	return nil
}

func (m *tui) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputsAuth))

	// Only text inputsAuth with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputsAuth {
		m.inputsAuth[i], cmds[i] = m.inputsAuth[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m tui) View() string {
	var b strings.Builder
	var s string
	//fmt.Fprintf(&b, "\nFocus: %d\n", m.focusIndex)
	if m.nonAuth {
		s = m.viewAuth()
	} else {
		s = m.viewMain()
	}
	b.WriteString(s)
	fmt.Fprintf(&b, "\n%s\n", m.errorString)
	return b.String()
}

func (m tui) viewAuth() string {
	var b strings.Builder

	for i := range m.inputsAuth {
		b.WriteString(m.inputsAuth[i].View())
		if i < len(m.inputsAuth)-1 {
			b.WriteRune('\n')
		}
	}

	loginBtn := &blurredLogin
	if m.focusIndex == 2 {
		loginBtn = &focusedLogin
	}
	registerBtn := &blurredRegister
	if m.focusIndex == 3 {
		registerBtn = &focusedRegister
	}
	fmt.Fprintf(&b, "\n\n%s\n%s\n", *loginBtn, *registerBtn)

	return b.String()
}

func (m tui) viewMain() string {

	return m.lpTable.View()
}

func (m *tui) getData() {
	lps := m.usecase.GetLoginsPasswords()
	rows := make([]table.Row, 0, len(lps))
	for _, lp := range lps {
		rows = append(rows, table.Row{strconv.Itoa(lp.Key), lp.Login, lp.Password, lp.Meta})
	}
	m.lpTable.SetColumns([]table.Column{
		{Title: "ID", Width: 4},
		{Title: "Login", Width: 10},
		{Title: "Password", Width: 10},
		{Title: "Meta", Width: 10}})
	m.lpTable.SetRows(rows)

}

func StartTUI(uc usecase) {
	tea.NewProgram(NewTUI(uc)).Run()
}
