package authtui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
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

type authtui struct {
	focusIndex int
	inputs     []textinput.Model

	offline     bool
	errorString string
}

func NewAuthTUI() authtui {
	m := authtui{
		inputs: make([]textinput.Model, 2),
	}
	var t textinput.Model
	for i := range m.inputs {
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

		m.inputs[i] = t
	}

	return m
}

