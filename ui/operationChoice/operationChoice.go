package operationChoice

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styling
var (
	focusedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	titleStyle            = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
)

// A Selection represents the choice made in a single-select step
type Selection struct {
	Choice string
}

// A model contains the data for the single-selection step.
// It implements the bubbletea.Model interface.
type model struct {
	cursor   int
	options  []string
	selection *Selection
	header   string
	exit     *bool
}

func (m model) Init() tea.Cmd {
	return nil
}

// InitialModel initializes a single-selection step with the given data
func InitialModelSelectionInput(options []string, selection *Selection, header string) model {
	return model{
		options:   options,
		selection: selection,
		header:    titleStyle.Render(header),
	}
}

// Update handles user input and updates the model state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			*m.exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter", "tab":
			m.selection.Choice = m.options[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the single-selection step
func (m model) View() string {
	s := m.header + "\n\n"

	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
		}

		option := focusedStyle.Render(option)

		s += fmt.Sprintf("%s %s\n", cursor, option)
	}

	s += fmt.Sprintf("Press %s to select an option.\n", focusedStyle.Render("enter"))
	return s
}
