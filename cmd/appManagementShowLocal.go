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
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/tabwriter"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/app_management"
	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"
)

// appManagementShowLocalCmd represents the appManagementShowLocal command
var appManagementShowLocalCmd = &cobra.Command{
	Use:   "local <appid>",
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

		appID := cmd.Flags().Arg(0)

		client, err := app_management.NewClientWithResponses(url)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		if useYAML {
			return showYAML(ctx, cmd.OutOrStdout(), client, appID, useColor)
		}

		if err := showAppList(ctx, cmd.OutOrStdout(), client, appID, useColor); err != nil {
			return err
		}

		if err := showContainers(ctx, cmd.OutOrStdout(), client, appID); err != nil {
			return err
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

func showYAML(ctx context.Context, writer io.Writer, client *app_management.ClientWithResponses, appID string, useColor bool) error {
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
		if err := quick.Highlight(writer, string(response.Body), "yaml", "terminal8", "native"); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(writer, string(response.Body)); err != nil {
			return err
		}
	}

	return nil
}

func showContainers(ctx context.Context, writer io.Writer, client *app_management.ClientWithResponses, appID string) error {
	response, err := client.ComposeAppContainersWithResponse(ctx, appID)
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

	w := tabwriter.NewWriter(writer, 0, 0, 3, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "CONTAINER NAME\tCONTAINER ID\tIMAGE\tSTATE")
	fmt.Fprintln(w, "--------------\t------------\t-----\t-----")

	mainApp := *response.JSON200.Data.Main
	for id, container := range *response.JSON200.Data.Containers {

		name := container.Name
		if id == mainApp {
			name = fmt.Sprintf("%s (main)", name)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			name,
			container.ID,
			container.Image,
			container.State,
		)
	}

	return nil
}

func showAppList(ctx context.Context, writer io.Writer, client *app_management.ClientWithResponses, appID string, useColor bool) error {
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
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			message := string(body)
			if message == "" {
				message = "is the casaos-app-management service running?"
			}

			return fmt.Errorf("%s - %s", response.Status, message)
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

	storeInfo, err := composeAppStoreInfo(data)
	if err != nil {
		return err
	}

	mainApp := "unknown"
	if storeInfo.Main != nil {
		mainApp = *storeInfo.Main
	}

	for name, app := range *storeInfo.Apps {
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
		fmt.Fprintln(writer, line)
		fmt.Fprintln(writer, strings.Repeat("-", len(line)))

		if useColor {
			if err := quick.Highlight(writer, buf.String(), "json", "terminal8", "native"); err != nil {
				return err
			}
		} else {
			if _, err := io.WriteString(writer, buf.String()); err != nil {
				return err
			}
		}

		fmt.Fprintln(writer)
	}

	return nil
}
