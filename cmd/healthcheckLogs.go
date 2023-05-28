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
	"os"
	"time"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/casaos"
	"github.com/spf13/cobra"
)

// healthcheckLogsCmd represents the healthcheckLogs command
var healthcheckLogsCmd = &cobra.Command{
	Use:     "logs",
	Short:   "get all `casaos-*` logs and save to a ZIP file",
	Aliases: []string{"log"},
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

		fmt.Println("getting logs...")

		response, err := client.GetHealthlogsWithResponse(ctx)
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

		outDir, err := cmd.Flags().GetString(FlagDir)
		if err != nil {
			return err
		}

		if outDir == "" {
			outDir, err = os.MkdirTemp("", "casaos-cli-*")
			if err != nil {
				return err
			}
		}

		zipFilePath := fmt.Sprintf("%s/casaos-%s-logs-%s.zip", outDir, Version, time.Now().Format("20060102150405"))

		if err := os.WriteFile(zipFilePath, response.Body, 0o600); err != nil {
			return err
		}

		fmt.Printf("logs saved to %s\n", zipFilePath)

		return nil
	},
}

func init() {
	healthcheckCmd.AddCommand(healthcheckLogsCmd)

	healthcheckLogsCmd.Flags().StringP(FlagDir, "d", "", "output directory")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// healthcheckLogsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// healthcheckLogsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
