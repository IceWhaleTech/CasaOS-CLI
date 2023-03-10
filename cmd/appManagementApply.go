/*
Copyright © 2023 IceWhaleTech

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
	"log"

	"github.com/spf13/cobra"
)

// appManagementApplyCmd represents the appManagementApply command
var appManagementApplyCmd = &cobra.Command{
	Use:   "apply <appid>",
	Short: "apply changes to an installed compose app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	appManagementCmd.AddCommand(appManagementApplyCmd)

	appManagementInstallCmd.Flags().StringP(FlagAppManagementFile, "f", "", "path to a compose file")
	if err := appManagementInstallCmd.MarkFlagRequired(FlagAppManagementFile); err != nil {
		log.Fatalln(err.Error())
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementApplyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementApplyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
