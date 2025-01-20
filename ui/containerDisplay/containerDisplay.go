package containerDisplay

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
	_, containers := RetrieveContainers(source)

	// if err != nil{
	// 	pass
	// }

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

// environment variable must be set, should control environment var for remote vs local config from .json
// configure based on source input
func RetrieveContainers(source string) (error, ContainerList) {
	// host := "unix:///var/run/docker.sock"
	// version := ""
	// cert_path := ""
	// tls_verify := ""

	// if source == "local"{
	// 	pass
	// } else{
	// 	pass
	// }

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err, nil
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return err, nil
	}

	var list ContainerList

	for _, ctr := range containers {
		details := ContainerDetails{
			PortMappings: formatPorts(ctr.Ports),
			CreatedTime: time.Unix(ctr.Created, 0).Format(time.RFC3339),
			Names: formatContainerNames(ctr.Names),
			ID: ctr.ID[:8],
			Image: ctr.Image,
			Status: ctr.Status,
			Command: ctr.Command,
		}

		list = append(list, details)
	}
	return nil, list
}

// formatPorts converts Docker port mappings to a readable string
func formatPorts(ports []types.Port) string {
	if len(ports) == 0 {
		return "-"
	}

	portStr := ""
	for _, port := range ports {
		if port.PublicPort > 0 {
			portStr += fmt.Sprintf("%d->%d/%s, ", port.PublicPort, port.PrivatePort, port.Type)
		} else {
			portStr += fmt.Sprintf("%d/%s, ", port.PrivatePort, port.Type)
		}
	}
	return portStr[:len(portStr)-2] // Remove trailing comma and space
}

// formatContainerNames removes leading '/' from container names
func formatContainerNames(names []string) string {
	formattedNames := make([]string, len(names))
	for i, name := range names {
		formattedNames[i] = name[1:] // Remove leading '/'
	}
	return formatSlice(formattedNames)
}

// formatSlice converts a slice to a comma-separated string
func formatSlice(slice []string) string {
	if len(slice) == 0 {
		return "-"
	}
	result := ""
	for _, s := range slice {
		result += s + ", "
	}
	return result[:len(result)-2] // Remove trailing comma and space
}