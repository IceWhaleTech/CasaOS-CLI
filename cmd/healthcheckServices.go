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
	"strings"
	"text/tabwriter"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/casaos"
	"github.com/spf13/cobra"
)

// healthcheckServicesCmd represents the healthcheckServices command
var healthcheckServicesCmd = &cobra.Command{
	Use:     "services",
	Short:   "get running status of each `casaos-*` service",
	Aliases: []string{"svc", "service"},
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathCasaOS)

		client, err := casaos.NewClientWithResponses(url)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		response, err := client.GetHealthServicesWithResponse(ctx)
		if err != nil {
			return err
		}

		if response.StatusCode() != http.StatusOK {
			var baseResponse casaos.BaseResponse
			if err := json.Unmarshal(response.Body, &baseResponse); err != nil {
				return fmt.Errorf("%s - %s", response.Status(), response.Body)
			}

			return fmt.Errorf("%s - %s", response.Status(), *baseResponse.Message)
		}

		if response.JSON200 == nil || response.JSON200.Data == nil {
			return fmt.Errorf("response body is empty")
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "NAME\tSTATUS\t")
		fmt.Fprintln(w, "----\t------\t")

		if response.JSON200.Data.Running != nil {
			for _, service := range *response.JSON200.Data.Running {
				fmt.Fprintf(w, "%s\t%s\n", strings.TrimSuffix(service, ".service"), "running")
			}
		}

		if response.JSON200.Data.NotRunning != nil {
			for _, service := range *response.JSON200.Data.NotRunning {
				fmt.Fprintf(w, "%s\t%s\n", strings.TrimSuffix(service, ".service"), "not running")
			}
		}

		return nil
	},
}

func init() {
	healthcheckCmd.AddCommand(healthcheckServicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// healthcheckServicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// healthcheckServicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
