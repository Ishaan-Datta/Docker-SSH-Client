package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"

	"SSH-Client/ui/authChoice"
	"SSH-Client/ui/remoteChoice"
	"SSH-Client/ui/sourceChoice"
)

// type listOptions struct {
// 	options []string
// }

// type Options struct {
// 	Operation *operationChoice.Selection
// }

// func updateOperation(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
// 	var option string

// 	// var selection operationChoice.Selection

// 	options := Options{
// 		Operation: &operationChoice.Selection{},
// 	}

// 	listOfStuff := listOptions{
// 		options: []string{
// 			"Log into a container",
// 			"Send commands to container(s)",
// 			"Push a local file to container(s)",
// 			"Pull a remote file from a container",
// 			"View available containers",
// 		},
// 	}

// 	model := operationChoice.InitialModelSelectionInput(listOfStuff.options, options.Operation, "Select the operation you would like to peform on the container(s):")

// need to store attributed auth and store selected hosts in model variable

type Model struct {
	source        string
	sources       []string
	operation     string
	authorization string
	table         table.Model
	state         string
	succeeded     []authChoice.ConfigStatus
	failed        []authChoice.ConfigStatus
	// containers
}

func main() {
	m := Model{}
	m.state = "source"
	// initialize model w/ pointer and fields

	// pre-check before starting loop, for config file present, docker engine running

controlLoop:
	for {
		switch state := m.state; state {
		case "source":
			var lol bool
			var err error
			m.source, m.state, err = sourceChoice.RunForm()
			fmt.Printf("\n%v", m.source)
			if (err != nil) || lol {
				break controlLoop
			}
		case "remote":
			// var sources []string
			var err error

			m.sources, m.state, err = remoteChoice.RunForm("config.json")
			if err != nil {
				break controlLoop
			}
		// case "operation":
		// 	updateOperation(msg, m)
		case "authentication":
			var err error
			m.succeeded, m.failed, m.state, err = authChoice.RunForm("config.json", m.sources)

			if err != nil {
				break controlLoop
			}
		// case "container table":
		// 	updateTable(msg, m)
		case "done":
			break controlLoop
		}
	}
}
