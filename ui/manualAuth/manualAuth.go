package manualAuth

import (
	"fmt"
	"os"

	"SSH-Client/ui/authChoice"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type Model struct {
	sources []string
	form    *huh.Form
	exit    *bool
	width   int
	height  int
}

func NewModel() Model {
	m := Model{}
	exit := false
	m.exit = &exit

	var choice string

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Remote Sources").
				Key("sources").
				Title("Error: ..."). // fix this
				Options(
					huh.NewOption("Change Selection", "Change Selection"),                     // go back
					huh.NewOption("Enter Credentials Manually", "Enter Credentials Manually"), // go to manual entry
					huh.NewOption("Exit", "Exit"),                                             // exit
				).
				Value(&choice),
		),
		// huh.NewGroup(
		// 	huh.NewInput().
		// )
	).
		WithShowHelp(false).
		WithShowErrors(false)
		// WithKeyMap()

	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// make sure these keys always quit
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
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

func RunForm(failed []authChoice.ConfigStatus) (string, string, error) {
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

// configList []FailedAttempt
