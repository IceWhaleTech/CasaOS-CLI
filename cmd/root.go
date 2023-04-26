/*
Copyright © 2022 IceWhaleTech

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
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	FlagRootURL = "root-url"
	FlagForce   = "force"
	FlagDryRun  = "dry-run"

	DefaultTimeout = 10 * time.Second
	RootGroupID    = "casaos-cli"
)

var (
	Version string
	Commit  string
	Date    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "casaos-cli",
	Short: "A command line interface for CasaOS",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// TODO - read from /etc/casaos/gateway.ini
	rootCmd.PersistentFlags().StringP(FlagRootURL, "u", "localhost:80", "root url of CasaOS API")
	rootCmd.AddGroup(&cobra.Group{
		ID:    RootGroupID,
		Title: "Services",
	})
}

func trim(s string, l uint) string {
	if len(s) > int(l) {
		return s[:l] + "..."
	}
	return s
}
