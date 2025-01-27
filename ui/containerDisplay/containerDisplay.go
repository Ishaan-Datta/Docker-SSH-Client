package containerDisplay

// pass in variable for list of containers, whether multi or single...
// choose containers from var list earlier, should define choice type from multi or single based on entering command
// could also be no select, so displays for both operations, just change help and what keybinds do depending on input
// multiple columns to display? maybe details in side-view

// manually implementing table but multi selectable

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "SSH-Client/utils"
)

// add filter at top later (copy from fancy list)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	Table table.Model
}

type ContainerDetails struct {
    ID          string
    Image       string
    Status      string
    Command     string
    CreatedTime string
    Names       string
    PortMappings string
}

type ContainerList []ContainerDetails

func (m model) Init() tea.Cmd { 
	return nil 
}

// format for use by bubbletea table
// sort by column...

func InitialTable(source string) model {
	// pull container list from utility function
	// err, containers := utils.RetrieveContainers(utils.RetrieveLocalConfiguration())


	// if err != nil{
	// 	os.Exit(1)
	// }

	containers := []ContainerDetails{}

	columns := []table.Column{
		{Title: "Container ID", Width: 10},
		{Title: "Image", Width: 10},
		{Title: "Command", Width: 10},
		{Title: "Created", Width: 10},
		{Title: "Status", Width: 10},
		{Title: "Ports", Width: 10},
		{Title: "Names", Width: 10},
	}

	var rows []table.Row

	for _, ctr := range containers {
		rows = append(rows, table.Row{ctr.ID, ctr.Image, ctr.Command, ctr.CreatedTime, ctr.Status, ctr.PortMappings, ctr.Names})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{
		Table: t,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.Table.View()) + "\n"
}

// configure based on source input and restore environment variable afterwards

// table: sort by filter w/ fuzzy search for field selected