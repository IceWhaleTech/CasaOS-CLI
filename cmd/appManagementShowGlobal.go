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
	"text/tabwriter"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/app_management"
	"github.com/spf13/cobra"
)

// appManagementShowGlobalCmd represents the appManagementShowGlobal command
var appManagementShowGlobalCmd = &cobra.Command{
	Use:  "global <Key>",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathAppManagement)

		client, err := app_management.NewClientWithResponses(url)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		response, err := client.GetGlobalSettingsWithResponse(ctx)
		if err != nil {
			return err
		}

		var baseResponse app_management.GetGlobalSettingsResponse

		if response.StatusCode() != http.StatusOK {
			if err := json.Unmarshal(response.Body, &baseResponse); err != nil {
				return fmt.Errorf("%s - %s", response.Status(), response.Body)
			}

			return fmt.Errorf("%s - %s", response.Status(), baseResponse.Body)
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "Global Key\tGlobal Value")
		fmt.Fprintln(w, "--------------\t------------")

		for _, value := range *response.JSON200.Data {

			fmt.Fprintf(w, "%s\t%s\t\n",
				*value.Key,
				value.Value,
			)
		}

		return nil
	},
}

func init() {
	appManagementShowCmd.AddCommand(appManagementShowGlobalCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementShowGlobalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementShowGlobalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
