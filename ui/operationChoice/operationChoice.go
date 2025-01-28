package operationChoice

import (
	"fmt"
	"os"

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
	// check json config
	// docker engine running
	// containers list isnt empty -> log containers based on type
	model := NewModel()

	_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	formSource := model.form.GetBool("source")
	fmt.Printf("formSource value: %v\n", formSource)
	source := fmt.Sprintf("%v", formSource)

	fmt.Printf("\n%v", *model.exit)

	// if the source is remote, then trigger auth, else can do as normal

	return source, "remote", nil
}

// 	// m.form = huh.NewForm(
// 	// 	huh.NewGroup(
// 	// 		huh.NewSelect[string]().
// 	// 		Key("operation").
// 	// 		Title("Select the operation you would like to peform on the container(s):").
// 	// 		Options(
// 	// 			huh.NewOption("Log into a container", "Log into a container"), // docker exec
// 	// 			huh.NewOption("Send commands to container(s)", "Send commands to container(s)"), // binary ad hoc or script file
// 	// 			huh.NewOption("Push a local file to container(s)", "Push a local file to container(s)"), // file selector or text input, prompt to overwrite, progress bar/spinner
// 	// 			huh.NewOption("Pull a remote file from a container", "Pull a remote file from a container"), // file selector or text input, prompt to overwrite, file not found error
// 	// 			huh.NewOption("View available containers", "View available containers"), // 
// 	// 		).
// 	// 		Value(&option),
// 	// 	),
// 	// )

// 	// err := m.form.Run()
// 	// if err != nil {
// 	// 	return m, tea.Quit
// 	// }
// 	m.operation = option

// 	// check 