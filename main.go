package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"SSH-Client/ui/containerDisplay"
)

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

// error header text?

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
	source string
	operation string
	authorization string
	table table.Model
	quitting bool
	showTable bool
	state string
}

// type Options struct {
// 	ProjectName *textInput.Output //pointer reference since we will be modifying the value, not just passing in a copy value
// 	ProjectType *multiInput.Selection
// }

func NewModel() Model {
	m := Model{}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.source = ""
	m.operation = ""
	m.authorization = ""
	m.quitting = false
	m.showTable = false

	var source bool
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
			Key("source").
			Title("Would you like to connect to remote or local container(s)?").
			Value(&source).
			Affirmative("Remote").
			Negative("Local"),
		),
	).
		WithShowHelp(false)
		// WithShowErrors(false)

	err := m.form.Run()

	if err != nil {
		fmt.Errorf("ruh roh")
	}

	m.source = fmt.Sprintf("%T", source)
	m.state = "operation"
	return m
}

// we never started that shit...
func (m Model) Init() tea.Cmd {
    return m.form.Init()
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// make sure these keys always quit
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	}

	// hands off message and model to the appropriate update function for
	// appropriate view based on current state
	outerLoop: 
		for {
			switch state := m.state; state {
			case "source":
				return updateSource(msg, m)
			case "operation":
				return updateOperation(msg, m)
			case "authorization":
				return updateAuthorization(msg, m)
			case "container table":
				return updateTable(msg, m)
			case "done":
				break outerLoop
			}
			// if m.quitting {break}
		}

	return m, tea.Quit
}

func (m Model) View() string {
	var s string

	switch state := m.state; state {
	case "source":
		s = m.form.View()
	case "operation":
		s = m.form.View()
	case "authorization":
		s = m.form.View()
	case "container table":
		s = m.table.View()
	}

	return s
}

func updateSource(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var source bool

	m.form = nil

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
			Key("source").
			Title("Would you like to connect to remote or local container(s)?").
			Value(&source).
			Affirmative("Remote").
			Negative("Local"),
		),
	)

	m.form.Init()
	
	err := m.form.Run()
	if err != nil {
		return m, tea.Quit
	}

	m.source = fmt.Sprintf("%T", source)

	// check json config
	// docker engine running
	// containers list isnt empty -> log containers based on type

	m.state = "operation"
	return m, nil
}

func updateOperation(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var option string

	m.form = nil

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
			Key("operation").
			Title("Select the operation you would like to peform on the container(s):").
			Options(
				huh.NewOption("Log into a container", "Log into a container"), // docker exec
				huh.NewOption("Send commands to container(s)", "Send commands to container(s)"), // binary ad hoc or script file
				huh.NewOption("Push a local file to container(s)", "Push a local file to container(s)"), // file selector or text input, prompt to overwrite, progress bar/spinner
				huh.NewOption("Pull a remote file from a container", "Pull a remote file from a container"), // file selector or text input, prompt to overwrite, file not found error
				huh.NewOption("View available containers", "View available containers"), // 
			).
			Value(&option),
		),
	)

	err := m.form.Run()
	if err != nil {
		return m, tea.Quit
	}
	m.operation = option

	// check 

	if m.operation == "View available containers"{
		m.table = containerDisplay.InitialTable(m.source).Table
		m.state = "container table"
	}else {
		m.state = "authorization"
	}

	return m, nil
}

func updateTable(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

func updateAuthorization(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var option string

	m.form = nil

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
			Key("auth").
			Title("Select the authentication method you would like to use:").
			Options(
				huh.NewOption("Kerberos", "Kerberos"),
				huh.NewOption("OAuth 2.0", "OAuth 2.0"),
				huh.NewOption("SAML", "SAML"),
				huh.NewOption("SSH Key", "SSH Key"),
			).
			Value(&option),
		),
	)

	err := m.form.Run()
	if err != nil {
		return m, tea.Quit
	}
	m.authorization = option

	// verify config has credentials, use auth gRPC stuff using the credentials

	m.state = "done"

	return m, nil
}

func main() {
	// tea.WithAltScreen()

	// containerDisplay.RetrieveContainers()

	// os.Exit(1)

	_, err := tea.NewProgram(NewModel()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}