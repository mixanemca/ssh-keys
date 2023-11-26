/*
Copyright © 2023 Michael Bruskov <mixanemca@yandex.ru>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mixanemca/ssh-keys/internal/keys"
)

type Model struct {
	// Keys stores the private keys paths.
	Keys []string
	// selectedIndex stores index of current selected private key.
	selectedIndex int
}

// NewModel is an initializer which creates a new model for rendering
// our Bubbletea app.
func NewModel() (*Model, error) {
	return &Model{}, nil
}

// Ensure that model fulfils the tea.Model interface at compile time.
var _ tea.Model = (*Model)(nil)

// View renders output to the CLI.
func (m *Model) View() string {
	var keys []string
	for i, k := range m.Keys {
		if i == m.selectedIndex {
			keys = append(keys, fmt.Sprintf("-> %s", k))
		} else {
			keys = append(keys, fmt.Sprintf("   %s", k))
		}
	}

	return fmt.Sprintf(`Found private keys:
%s

Press enter/return to select a list item, arrow keys to move, Ctrl+C or q to exit.`,
		strings.Join(keys, "\n"))
}

// Update is called with a tea.Msg, representing something that happened within
// our application.
//
// This can be things like terminal resizing, keypresses, or custom IO.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Let's figure out what is in tea.Msg, and what we need to do.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// The terminal was resized.  We can access the new size with:
		_, _ = msg.Width, msg.Height
	case tea.KeyMsg:
		// msg is a keypress. We can handle each key combo uniquely, and update
		// our state:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "down":
			return m.moveCursor(msg), nil
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			// In this case, ctrl+c quits the app by sending a
			// tea.Quit cmd. This is a Bubbletea builtin which terminates the
			// overall framework which renders our model.
			//
			// Unfortunately, if you don't include this quitting can be, uh,
			// frustrating, as bubbletea catches every key combo by default.
			return m, tea.Quit
		}
	}
	// We return an updated model to Bubbletea for rendering here.  This allows
	// us to mutate state so that Bubbletea can render an updated view.
	//
	// We also return "commands".  A command is something that you need to do
	// after rendering.  Each command produces a tea.Msg which is its *result*.
	// Bubbletea calls this Update function again with the tea.Msg - this is our
	// render loop.
	//
	// For now, we have no commands to run given the message is not a keyboard
	// quit combo.
	return m, nil
}

// Init is called to kick off the render cycle.  It allows you to
// perform IO after the app has loaded and rendered once, asynchronously.
// The tea.Cmd can return a tea.Msg which will be passed into Update() in order
// to update the model's state.
func (m *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, findPrivateKeys(m))

	return tea.Batch(cmds...)
}

func (m *Model) moveCursor(msg tea.KeyMsg) *Model {
	switch msg.String() {
	case "up":
		m.selectedIndex--
	case "down":
		m.selectedIndex++
	default:
		// do nothing
	}

	keysCount := len(m.Keys)
	if keysCount != 0 {
		m.selectedIndex = (m.selectedIndex + keysCount) % keysCount
	}

	return m
}

// findPrivateKeys finds the SSH private keys in user's home directory.
func findPrivateKeys(m *Model) tea.Cmd {
	return func() tea.Msg {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("Failed to get user home dir: ", err)
		}

		sshDir := filepath.Join(home, ".ssh")

		m.Keys, err = keys.LoadPrivateKeys(sshDir)
		if err != nil {
			log.Fatal("Failed to load private keys: ", err)
		}

		return m
	}
}
