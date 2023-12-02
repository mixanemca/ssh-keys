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
	"bytes"
	"log"
	"net"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mixanemca/ssh-keys/internal/keys"
	"golang.org/x/crypto/ssh/agent"
)

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

// findAgentKeys finds the SSH keys, added to SSH agent.
func findAgentKeys(m *Model) tea.Cmd {
	return func() tea.Msg {
		// ssh-agent(1) provides a UNIX socket at $SSH_AUTH_SOCK.
		socket := os.Getenv("SSH_AUTH_SOCK")
		conn, err := net.Dial("unix", socket)
		if err != nil {
			log.Fatal("Failed to open SSH_AUTH_SOCK: ", err)
		}

		m.AgentClient = agent.NewClient(conn)
		keys, err := m.AgentClient.List()
		if err != nil {
			log.Fatal("Failed to get list of ssh keys from agent: ", err)
		}

		for _, k := range keys {
			m.AgentKeys = append(m.AgentKeys, k.Blob)
		}

		return m
	}
}

// loadKeyToAgent loads selected key to SSH agent and updates Model.
func loadKeyToAgent(m *Model) tea.Cmd {
	return func() tea.Msg {
		if err := m.AgentClient.Add(agent.AddedKey{
			PrivateKey: m.Keys[m.selectedIndex].Private,
			Comment:    m.Keys[m.selectedIndex].Comment,
		}); err != nil {
			log.Fatal("Failed to load key to ssh-agent: ", err)
		}
		// Add the key to Model.
		m.AgentKeys = append(m.AgentKeys, m.Keys[m.selectedIndex].Public.Marshal())

		return m
	}
}

// unloadKeyFromAgent removes selected key from SSH agent and updates Model.
func unloadKeyFromAgent(m *Model) tea.Cmd {
	return func() tea.Msg {
		// Unload the key from agent.
		err := m.AgentClient.Remove(m.Keys[m.selectedIndex].Public)
		if err != nil {
			log.Fatal("Failed to unload key from ssh-agent: ", err)
		}

		m.Keys[m.selectedIndex].LoadedToAgent = false

		// Remove the key from Model.
		for i, ak := range m.AgentKeys {
			if bytes.Equal(ak, m.Keys[m.selectedIndex].Public.Marshal()) {
				m.AgentKeys = append(m.AgentKeys[:i], m.AgentKeys[i+1:]...)
			}
		}

		return m
	}
}
