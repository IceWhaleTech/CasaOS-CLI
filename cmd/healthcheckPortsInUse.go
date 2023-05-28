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
	"sort"
	"text/tabwriter"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/casaos"
	"github.com/spf13/cobra"
)

// healthcheckPortsInUseCmd represents the healthcheckPortsInUse command
var healthcheckPortsInUseCmd = &cobra.Command{
	Use:     "ports-in-use",
	Short:   "get ports in use",
	Aliases: []string{"ports", "port-in-use"},
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

		response, err := client.GetHealthPortsWithResponse(ctx)
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

		fmt.Fprintln(w, "PORT\tTYPE\t")
		fmt.Fprintln(w, "----\t----\t")

		if response.JSON200.Data.TCP != nil {
			sort.Ints(*response.JSON200.Data.TCP)
			for _, port := range *response.JSON200.Data.TCP {
				fmt.Fprintf(w, "%d\t%s\n", port, "TCP")
			}
		}

		if response.JSON200.Data.UDP != nil {
			sort.Ints(*response.JSON200.Data.UDP)
			for _, port := range *response.JSON200.Data.UDP {
				fmt.Fprintf(w, "%d\t%s\n", port, "UDP")
			}
		}

		return nil
	},
}

func init() {
	healthcheckCmd.AddCommand(healthcheckPortsInUseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// healthcheckPortsInUseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// healthcheckPortsInUseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
