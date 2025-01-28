package authChoice

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"SSH-Client/utils"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	failedMark          = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")
)

type Model struct {
	exit      *bool
	failed    []ConfigStatus
	succeeded []ConfigStatus
	configs   []utils.Config
	index     int
	width     int
	height    int
	spinner   spinner.Model
	progress  progress.Model
	done      bool
}

type FailedAttempt struct {
	hostname string
	reason   error
}

func NewModel(configs []utils.Config) Model {
	m := Model{}
	exit := false
	m.exit = &exit
	m.failed = []ConfigStatus{}
	m.succeeded = []ConfigStatus{}
	m.configs = configs

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	m.progress = p
	m.spinner = s

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(callAuthRequest(m.configs[m.index]), m.spinner.Tick)
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
	case ConfigStatus:
		var style lipgloss.Style
		config := m.configs[m.index]
		if msg.err != nil {
			m.failed = append(m.failed, msg)
			style = failedMark
		} else {
			m.succeeded = append(m.succeeded, msg)
			style = checkMark
		}

		if m.index >= len(m.configs)-1 {
			m.done = true
			return m, tea.Sequence(
				tea.Printf("%s %s", style, config.Host), // print the last success message
				tea.Quit,                                // exit the program
			)
		}

		// Update progress bar
		m.index++
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.configs)))

		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", style, config.Host), // print success message above our program
			callAuthRequest(m.configs[m.index]),     // authenticate the next config
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}

	*m.exit = false
	return m, nil
}

func (m Model) View() string {
	n := len(m.configs)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("Done! Completed authentication for %d configurations.\n", n))
	}

	configCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+configCount))

	configName := currentPkgNameStyle.Render(m.configs[m.index].Host)
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Authenticating " + configName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+configCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + configCount
}

func RunForm(configFilePath string, remoteSources []string) ([]ConfigStatus, []ConfigStatus, string, error) {
	remoteConfigs, err := utils.RetrieveRemoteConfiguration(configFilePath)
	if err != nil {
		return nil, nil, "done", err
	}

	var toAuthenticate []utils.Config

	for _, config := range remoteConfigs {
		if contains(remoteSources, config.Host) && config.TLSVerify == "1" {
			toAuthenticate = append(toAuthenticate, config)
		}
	}

	var succeeded []ConfigStatus
	var failed []ConfigStatus

	if len(toAuthenticate) > 0 {
		model := NewModel(toAuthenticate)
		_, err = tea.NewProgram(model).Run() // tea.WithAltScreen()
		if err != nil {
			return nil, nil, "done", err
		}

		if *model.exit {
			return nil, nil, "done", nil
		}

		succeeded = model.succeeded
		failed = model.failed
	}

	// decide based on failed list > 0 to assign state

	return succeeded, failed, "done", nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type ConfigStatus struct {
	hostname string
	config   utils.Config
	err      error
}

// implement request calls here from auth.go, should
func callAuthRequest(config utils.Config) tea.Cmd {
	// return func() tea.Msg {
	// 	// var err error
	// 	switch config.Type {
	// 	case "Kerberos":
	// 		// Call Kerberos auth request function
	// 		fmt.Printf("Calling Kerberos auth request with CertPath: %s\n", config.CertPath)
	// 	case "SAML":
	// 		// Call SAML auth request function
	// 		fmt.Printf("Calling SAML auth request with CertPath: %s\n", config.CertPath)
	// 	case "OAuth2":
	// 		// Call OAuth2 auth request function
	// 		fmt.Printf("Calling OAuth2 auth request with CertPath: %s\n", config.CertPath)
	// 	default:
	// 		fmt.Printf("Unsupported auth type: %s\n", config.Type)
	// 	}
	// 	return ConfigStatus{hostname: config.Host, config: config, err: nil}
	// }

	d := time.Millisecond * time.Duration(rand.Intn(2000)) //nolint:gosec
	return tea.Tick(d, func(t time.Time) tea.Msg {
		// return ConfigStatus{hostname: config.Host, config: config, err: fmt.Errorf("Failed to authenticate %s", config.Host)}
		return ConfigStatus{hostname: config.Host, config: config, err: nil}
	})
}
