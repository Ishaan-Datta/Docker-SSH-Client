package remoteChoice

import (
	"fmt"

	"SSH-Client/utils"

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
	sources []string
}

func NewModel(configList []string) Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	exit := false
	m.exit = &exit

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Remote Sources").
				Key("sources").
				Title("Choose the remote sources you would like to access").
				Options(huh.NewOptions(configList...)...).
				Validate(func(t []string) error {
					if len(t) <= 0 {
						return fmt.Errorf("at least one source is required")
					}
					return nil
				}),
				// Filterable(true).
				// Value(&m.sources),
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
		// *m.sources = 
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

func RunForm(configFilePath string) ([]string, string, error) {
	configList, err := utils.RetrieveRemoteConfigurationNames(configFilePath)
	if err != nil {
		return nil, "done", err
	}

	model := NewModel(configList)

	_, err = tea.NewProgram(model, tea.WithAltScreen()).Run()
	if err != nil {
		return nil, "done", err
	}

	if *model.exit {
		return nil, "done", nil
	}

	if sources, ok := model.form.Get("sources").([]string); ok {
		model.sources = sources
	}

	return model.sources, "authentication", nil
}