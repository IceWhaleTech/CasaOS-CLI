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

// appManagementSetGlobalCmd represents the appManagementSetGlobal command
var appManagementSetGlobalCmd = &cobra.Command{
	Use:   "global",
	Short: "set a global variable",
	Args: cobra.MatchAll(cobra.ExactArgs(2), func(cmd *cobra.Command, args []string) error {
		// the func sould to check the args.
		// to force the Key  is not a number is hard to do this.
		// So I didn't do this. this is a null check. May be you can do this.
		return nil
	}),
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

		key := args[0]
		body := app_management.GlobalSetting{
			Value: args[1],
		}
		response, err := client.UpdateGlobalSettingWithResponse(ctx, key, body)
		if err != nil {
			return err
		}

		var baseResponse app_management.BaseResponse

		if response.StatusCode() != http.StatusOK {
			if err := json.Unmarshal(response.Body, &baseResponse); err != nil {
				return fmt.Errorf("%s - %s", response.Status(), response.Body)
			}

			return fmt.Errorf("%s - %s", response.Status(), *baseResponse.Message)
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "Global Key\tGlobal Value")
		fmt.Fprintln(w, "--------------\t------------")

		fmt.Fprintf(w, "%s\t%s\t\n",
			*response.JSON200.Data.Key,
			response.JSON200.Data.Value,
		)

		return nil
	},
}

func init() {
	appManagementSetCmd.AddCommand(appManagementSetGlobalCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementSetGlobalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementSetGlobalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
