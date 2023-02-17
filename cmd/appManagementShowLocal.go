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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/app_management"
	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"
)

// appManagementShowLocalCmd represents the appManagementShowLocal command
var appManagementShowLocalCmd = &cobra.Command{
	Use:   "local [appid]",
	Short: "show information of a locally installed app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			return err
		}

		useYAML, err := cmd.Flags().GetBool(FlagAppManagementYAML)
		if err != nil {
			return err
		}

		useColor, err := cmd.Flags().GetBool(FlagAppManagementUseColor)
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

		appID := cmd.Flags().Arg(0)

		if useYAML {
			response, err := client.MyComposeAppWithResponse(ctx, appID, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("Accept", "application/yaml")
				return nil
			})
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

			if useColor {
				if err := quick.Highlight(os.Stdout, string(response.Body), "yaml", "terminal8", "native"); err != nil {
					return err
				}
			} else {
				fmt.Println(string(response.Body))
			}

			return nil
		}

		response, err := client.MyComposeApp(ctx, appID)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		buf, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			var baseResponse app_management.BaseResponse
			if err := json.Unmarshal(buf, &baseResponse); err != nil {
				return fmt.Errorf("%s - %s", response.Status, response.Body)
			}

			return fmt.Errorf("%s - %s", response.Status, *baseResponse.Message)
		}

		// get mapstruct of response body - can't unmarshal directly due to https://github.com/compose-spec/compose-go/issues/353
		var body map[string]interface{}
		if err := json.Unmarshal(buf, &body); err != nil {
			return err
		}

		_, ok := body["data"]
		if !ok {
			return fmt.Errorf("body does not contain `data`")
		}

		data, ok := body["data"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("data is not a map[string]interface")
		}

		mainApp, appList, err := appList(data)
		if err != nil {
			return err
		}

		for name, app := range appList {
			var buf bytes.Buffer

			enc := json.NewEncoder(&buf)
			enc.SetIndent("", "  ")

			if err := enc.Encode(app); err != nil {
				return err
			}

			line := name
			if name == mainApp {
				line += " (main)"
			}
			fmt.Println(line)
			fmt.Println(strings.Repeat("-", len(line)))

			if useColor {
				if err := quick.Highlight(os.Stdout, buf.String(), "json", "terminal8", "native"); err != nil {
					return err
				}
			} else {
				fmt.Println(buf.String())
			}

			fmt.Println()
		}

		return nil
	},
}

func init() {
	appManagementShowCmd.AddCommand(appManagementShowLocalCmd)

	appManagementShowLocalCmd.Flags().BoolP(FlagAppManagementYAML, "", false, "output in raw YAML format")
	appManagementShowLocalCmd.Flags().BoolP(FlagAppManagementUseColor, "c", false, "colorize output")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementShowLocalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementShowLocalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
