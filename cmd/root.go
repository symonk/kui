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
	"log/slog"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/symonk/kui/internal/env"
	"github.com/symonk/kui/internal/gui"
	"github.com/symonk/kui/internal/kafka"
	"github.com/symonk/kui/internal/log"
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
		logger, path, closer := setupLogger()
		defer closer()
		logger.Info("starting kui...", path)
		cfg := viper.GetViper().ConfigFileUsed()
		cfgMap, err := kafka.FileToKafkaMap(cfg)
		logCh := make(chan confluentKafka.LogEvent)
		p := makeDummyProducer(&cfgMap, logCh)
		defer p.Close()
		cobra.CheckErr(err)
		// TODO: This is all horrible, but buggy in confluent kafka go right now.
		// see the issue in the docstring of makeDummyProducer
		client, err := kafka.New(p)
		defer client.Close()
		cobra.CheckErr(err)
		done := make(chan struct{})
		go func() {
			redirectLogs(logger, logCh, done)
		}()
		defer close(done)
		program := tea.NewProgram(gui.New(client, logger, path))
		if _, err := program.Run(); err != nil {
			logger.Error("critical failure", slog.Any("reason", err))
		}
	},
}

type Closer func() error

// setupLogger configures an application logger, primarily used for debug purposes
// but also outputs the kafka stderr stream.
func setupLogger() (*slog.Logger, string, Closer) {
	tmpF := os.TempDir()
	logFile := path.Join(tmpF, "kui.log")
	writer, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := log.New(handler)
	return l, logFile, writer.Close
}

// redirectLogs asynchronously redirects kafka logs to the application logger.
func redirectLogs(logger *slog.Logger, ch chan confluentKafka.LogEvent, done chan struct{}) {
	for {
		select {
		case e := <-ch:
			logger.Debug("kafka_log",
				slog.String("name", e.Name),
				slog.String("tag", e.Tag),
				slog.String("message", e.Message),
				slog.Int("level", int(e.Level)),
				slog.Time("timestamp", e.Timestamp))
		case <-done:
			logger.Debug("exit signal, exiting")
			return
		}
	}
}

// makeDummyProducer returns a dummy kafka producer as a work around for an issue with
// providing a client supplied logger channel to the admin client.
//
// https://github.com/confluentinc/confluent-kafka-go/issues/1119
func makeDummyProducer(cfgMap *confluentKafka.ConfigMap, logsChan chan confluentKafka.LogEvent) *confluentKafka.Producer {
	cfgMap.SetKey("go.logs.channel.enable", true)
	cfgMap.SetKey("go.logs.channel", logsChan)
	p, err := confluentKafka.NewProducer(cfgMap)
	if err != nil {
		panic(err)
	}
	return p
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
			p = path.Join(home, ".config", "kui.conf")
			viper.SetConfigFile(p)
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
