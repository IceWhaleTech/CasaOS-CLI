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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/app_management"
	"github.com/spf13/cobra"
)

const (
	FlagAppManagementInstallFile = "file"
)

// appManagementInstallCmd represents the appManagementInstall command
var appManagementInstallCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"add", "create", "up"},
	Short:   "install a compose app",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathAppManagement)

		filepath := cmd.Flag(FlagAppManagementInstallFile).Value.String()

		file, err := os.Open(filepath)
		if err != nil {
			return err
		}

		client, err := app_management.NewClientWithResponses(url)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		response, err := client.InstallComposeAppWithBodyWithResponse(ctx, "application/yaml", file)
		if err != nil {
			return err
		}

		if response.StatusCode() != http.StatusOK {
			var baseResponse app_management.BaseResponse
			if err := json.Unmarshal(response.Body, &baseResponse); err != nil {
				return fmt.Errorf("%s - %s", response.Status(), response.Body)
			}

			return fmt.Errorf("%s - %s", response.Status(), *baseResponse.Message)
		}

		log.Println(*response.JSON200.Message)

		return nil
	},
}

func init() {
	appManagementCmd.AddCommand(appManagementInstallCmd)

	appManagementInstallCmd.Flags().StringP(FlagAppManagementInstallFile, "f", "", "path to a compose file")
	if err := appManagementInstallCmd.MarkFlagRequired(FlagAppManagementInstallFile); err != nil {
		log.Fatalln(err.Error())
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementInstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementInstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}