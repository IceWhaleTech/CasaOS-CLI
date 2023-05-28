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
	"log"
	"net/http"
	"os"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/app_management"
	"github.com/alecthomas/chroma/quick"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// appManagementConvertAppFileCmd represents the appManagementConvertAppFile command
var appManagementConvertAppFileCmd = &cobra.Command{
	Use:   "appfile",
	Short: "convert `appfile.json` to Docker Compose YAML (for local conversion, use `appfile2compose` command)",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			return err
		}

		filepath := cmd.Flag(FlagFile).Value.String()

		useColor, err := cmd.Flags().GetBool(FlagAppManagementUseColor)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathAppManagement)

		client, err := app_management.NewClientWithResponses(url)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		file, err := os.Open(filepath)
		if err != nil {
			return err
		}

		params := app_management.ConvertParams{Type: lo.ToPtr(app_management.Appfile)}

		response, err := client.ConvertWithBodyWithResponse(ctx, &params, MINEApplicationJSON, file)
		if err != nil {
			fmt.Println("Error: Unable to reach CasaOS API. Try convert locally using `appfile2compose` command.")
			return err
		}

		if response.StatusCode() != http.StatusOK {
			var baseResponse app_management.BaseResponse
			if err := yaml.Unmarshal(response.Body, &baseResponse); err != nil {
				return fmt.Errorf("%s - %s", response.Status(), response.Body)
			}

			return fmt.Errorf("%s - %s", response.Status(), *baseResponse.Message)
		}

		if useColor {
			return quick.Highlight(cmd.OutOrStdout(), string(response.Body), "yaml", "terminal8", "native")
		}

		fmt.Println(string(response.Body))

		return nil
	},
}

func init() {
	appManagementConvertCmd.AddCommand(appManagementConvertAppFileCmd)

	appManagementConvertAppFileCmd.Flags().StringP(FlagFile, "f", "", "path to the `appfile.json` file")
	if err := appManagementConvertAppFileCmd.MarkFlagRequired(FlagFile); err != nil {
		log.Fatalln(err.Error())
	}

	appManagementConvertAppFileCmd.Flags().BoolP(FlagAppManagementUseColor, "c", false, "colorize output")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementConvertAppFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementConvertAppFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
