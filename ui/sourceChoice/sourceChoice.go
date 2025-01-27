package sourceChoice

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 80

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type Model struct {
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	width  int
	exit   *bool
	// result bool
}

func NewModel() Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	exit := false
	m.exit = &exit

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("source").
				Title("Would you like to connect to remote or local container(s)?").
				Affirmative("Remote").
				Negative("Local"),
				// Value(&m.result),
		),
	).
	WithShowHelp(false).
	WithShowErrors(false)
		// WithKeyMap()
	
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// make sure these keys always quit
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			*m.exit = true
			// os.Exit(1)
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		*m.exit = false
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)

	// shift+tab for back
}

func (m Model) View() string {
	s := m.form.View()
	return s
}

func RunForm() (string, string, error) {
	// containers list isnt empty -> log containers based on type
	model := NewModel()

	_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	
	if err != nil {
		return "", "error", err
	}

	if *model.exit {
		return "", "error", err
	}

	formSource := model.form.GetBool("source")

	if formSource {
		// verify remote config is valid
		return "remote", "remote", nil
	} else {
		// verify local config is valid
		return "local", "operation", nil
	}

	// no valid client configs found
}