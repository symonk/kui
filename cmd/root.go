/*
Copyright © 2024 Simon Kerr

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
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/symonk/kui/internal/env"
	"github.com/symonk/kui/internal/gui"
	"github.com/symonk/kui/internal/kafka"
)

const (
	// cfgEnvironLookupKey is set to point at a config file and used
	// if no explicit -c value is provided.
	cfgEnvironLookupKey string = "KUI_CONFIG"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kui",
	Short: "A terminal ui for manging kafka",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := viper.GetViper().ConfigFileUsed()
		cfgMap, err := kafka.FileToKafkaMap(cfg)
		cobra.CheckErr(err)
		client, err := kafka.New(cfgMap)
		cobra.CheckErr(err)
		p := tea.NewProgram(gui.New(client))
		if _, err := p.Run(); err != nil {
			fmt.Println("critical failure", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "c", "", "config file (default is $HOME/.config/kui.conf)")
}

// initConfig attempts to resolve a configuration file of
// librdkafka properties.  The order is as follows:
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// check if the user has $KUI_CONFIG exported in their environment
		// if so, use it, otherwise look in the default config directory
		// for a file
		if p, ok := env.KeyIsInEnvironment(cfgEnvironLookupKey); ok {
			viper.SetConfigFile(p)
		} else {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			// Search config in home directory with name ".kui" (without extension).
			viper.AddConfigPath(path.Join(home, ".config"))
			viper.SetConfigType("conf")
			viper.SetConfigName("kui")
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
