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
	"os"

	"github.com/spf13/cobra"
)

// appManagementConvertAppFileCmd represents the appManagementConvertAppFile command
var appManagementConvertAppFileCmd = &cobra.Command{
	Use:   "appfile",
	Short: "convert to `docker-compose.yml` from an `appfile.json` exported by CasaOS v0.4.3 or earlier",
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := cmd.Flag(FlagAppManagementFile).Value.String()

		file, err := os.Open(filepath)
		if err != nil {
			return err
		}

		decoder := json.NewDecoder(file)

		return nil
	},
}

func init() {
	appManagementConvertCmd.AddCommand(appManagementConvertAppFileCmd)

	appManagementConvertAppFileCmd.Flags().StringP(FlagAppManagementFile, "f", "", "path to the `appfile.json` file")
	if err := appManagementConvertAppFileCmd.MarkFlagRequired(FlagAppManagementFile); err != nil {
		log.Fatalln(err.Error())
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementConvertAppFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementConvertAppFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
