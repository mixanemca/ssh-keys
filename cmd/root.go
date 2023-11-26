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

package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mixanemca/ssh-keys/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version string = "unknown"
	build   string = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "ssh-keys",
	Short:         "Work with SSH keys easily!",
	SilenceErrors: false,
	SilenceUsage:  false,
	Version:       version,
	Run:           run,
}

func init() {
	vt := rootCmd.VersionTemplate()
	rootCmd.SetVersionTemplate(vt[:len(vt)-1] + " (" + build + ")\n")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Create a new TUI model which will be rendered in Bubbletea.
	state, err := ui.NewModel()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error starting init command: %s\n", err))
		os.Exit(1)
	}
	// tea.NewProgram starts the Bubbletea framework which will render our
	// application using our state.
	if err := tea.NewProgram(state).Start(); err != nil {
		log.Fatal(err)
	}
}
