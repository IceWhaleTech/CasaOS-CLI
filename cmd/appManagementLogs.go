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
	"fmt"
	"net/http"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/app_management"
	"github.com/spf13/cobra"
)

const (
	FlagAppManagementLogsLines = "lines"
)

// appManagementLogsCmd represents the appManagementLogs command
var appManagementLogsCmd = &cobra.Command{
	Use:   "logs <appid>",
	Short: "retrieve logs of a compose app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			return err
		}

		lines, err := cmd.Flags().GetInt(FlagAppManagementLogsLines)
		if err != nil {
			return err
		}

		if lines < 0 {
			return fmt.Errorf("lines must be greater than 0")
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathAppManagement)

		client, err := app_management.NewClientWithResponses(url)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		appID := cmd.Flags().Arg(0)
		response, err := client.ComposeAppLogsWithResponse(ctx, appID, &app_management.ComposeAppLogsParams{Lines: &lines})
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

		fmt.Printf("(showing last %d lines)\n", lines)
		fmt.Println("...")
		fmt.Println(*response.JSON200.Data)

		return nil
	},
}

func init() {
	appManagementCmd.AddCommand(appManagementLogsCmd)

	appManagementLogsCmd.Flags().IntP(FlagAppManagementLogsLines, "l", 1000, "Follow log output")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementLogsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementLogsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
