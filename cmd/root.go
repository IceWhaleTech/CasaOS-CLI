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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-ini/ini"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

const (
	BasePathCasaOS = "v2/casaos"

	FlagDir     = "dir"
	FlagDryRun  = "dry-run"
	FlagFile    = "file"
	FlagForce   = "force"
	FlagRootURL = "root-url"

	GatewayPath = "/etc/casaos/gateway.ini"

	DefaultTimeout = 10 * time.Second
	RootGroupID    = "casaos-cli"
)

var (
	Version string
	Commit  string
	Date    string

	json = jsoniter.ConfigCompatibleWithStandardLibrary
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
	url := ""

	rootCmd.PersistentFlags().StringP(FlagRootURL, "u", "", "root url of CasaOS API")

	if rootCmd.PersistentFlags().Changed(FlagRootURL) {
		url = rootCmd.PersistentFlags().Lookup(FlagRootURL).Value.String()
	} else {
		if _, err := os.Stat(GatewayPath); err == nil {
			cfgs, err := ini.Load(GatewayPath)
			if err != nil {
				log.Println("No gateway config found, use default root url")
			}

			port := cfgs.Section("gateway").Key("port").Value()
			if port != "" {
				url = fmt.Sprintf("localhost:%s", port)
			}
		}
	}

	if url == "" {
		url = "localhost:80"
	}

	rootCmd.PersistentFlags().Set(FlagRootURL, url)
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
