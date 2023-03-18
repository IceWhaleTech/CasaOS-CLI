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
	"io"
	"net"
	"net/http"
	"strings"
	"text/tabwriter"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/app_management"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// appManagementListAppsCmd represents the appManagementListApps command
var appManagementListAppsCmd = &cobra.Command{
	Use:     "apps",
	Short:   "list locally installed apps",
	Aliases: []string{"app"},
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

		response, err := client.MyComposeAppList(ctx)
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

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "APPID\tSTATUS\tWEB UI\tDESCRIPTION")
		fmt.Fprintln(w, "-----\t------\t------\t-----------")

		for id, app := range data {
			status, err := status(app)
			if err != nil {
				status = "unknown"
			}

			mainApp, appList, err := appList(app)
			if err != nil {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					id,
					status,
					"n/a",
					"(not a CasaOS compose app)",
				)
				return nil
			}

			mainAppStoreInfo, ok := appList[mainApp]
			if !ok {
				return fmt.Errorf("main app not found in app list")
			}

			scheme := "http"
			if mainAppStoreInfo.Container.Scheme != nil {
				scheme = string(*mainAppStoreInfo.Container.Scheme)
			}

			hostname, err := hostname()
			if err != nil {
				return err
			}

			if mainAppStoreInfo.Container.Hostname != nil {
				hostname = *mainAppStoreInfo.Container.Hostname
			}

			webUI := fmt.Sprintf("%s://%s:%s/%s",
				scheme,
				hostname,
				mainAppStoreInfo.Container.PortMap,
				strings.TrimLeft(mainAppStoreInfo.Container.Index, "/"),
			)

			description := map[string]string{
				"en_US": "No description available",
			}

			if mainAppStoreInfo.Description != nil {
				description = mainAppStoreInfo.Description
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				id,
				status,
				webUI,
				trim(lo.Values(description)[0], 78),
			)
		}

		return nil
	},
}

func init() {
	appManagementListCmd.AddCommand(appManagementListAppsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementListAppsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementListAppsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func status(composeApp interface{}) (string, error) {
	composeAppMapStruct, ok := composeApp.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("app is not a map[string]interface{}")
	}

	_, ok = composeAppMapStruct["status"]
	if !ok {
		return "", fmt.Errorf("app does not have \"status\"")
	}

	status, ok := composeAppMapStruct["status"].(string)
	if !ok {
		return "", fmt.Errorf("app[\"status\"] is not a string")
	}

	return status, nil
}

func appList(composeApp interface{}) (string, map[string]app_management.AppStoreInfo, error) {
	composeAppMapStruct, ok := composeApp.(map[string]interface{})
	if !ok {
		return "", nil, fmt.Errorf("app is not a map[string]interface{}")
	}

	_, ok = composeAppMapStruct["store_info"]
	if !ok {
		return "", nil, fmt.Errorf("app does not have \"store_info\"")
	}

	composeAppStoreInfo, ok := composeAppMapStruct["store_info"].(map[string]interface{})
	if !ok {
		return "", nil, fmt.Errorf("app[\"store_info\"] is not a map[string]interface{}")
	}

	_, ok = composeAppStoreInfo["main_app"]
	if !ok {
		return "", nil, fmt.Errorf("app[\"store_info\"] does not have \"main_app\"")
	}

	mainApp, ok := composeAppStoreInfo["main_app"].(string)
	if !ok {
		return "", nil, fmt.Errorf("app[\"store_info\"][\"main_app\"] is not a string")
	}

	_, ok = composeAppStoreInfo["apps"]
	if !ok {
		return "", nil, fmt.Errorf("app[\"store_info\"] does not have \"apps\"")
	}

	appListMapStruct, ok := composeAppStoreInfo["apps"].(map[string]interface{})
	if !ok {
		return "", nil, fmt.Errorf("app[\"store_info\"][\"apps\"] is not a map[string]interface{}")
	}

	var appList map[string]app_management.AppStoreInfo

	if err := mapstructure.Decode(appListMapStruct, &appList); err != nil {
		return "", nil, err
	}

	return mainApp, appList, nil
}

func hostname() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("could not find hostname")
}
