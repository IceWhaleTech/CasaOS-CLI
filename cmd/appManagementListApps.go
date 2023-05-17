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

		ctx, cancel := context.WithCancel(context.Background())
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

			message := "(empty response body)"
			if baseResponse.Message != nil {
				message = *baseResponse.Message
			}

			return fmt.Errorf("%s - %s", response.Status, message)
		}

		data := json.Get(buf, "data")

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "APPID\tSTATUS\tWEB UI\tIMAGES\tDESCRIPTION")
		fmt.Fprintln(w, "-----\t------\t------\t------\t-----------")

		for _, id := range data.Keys() {
			app := data.Get(id)

			status := app.Get("status").ToString()

			images := []string{}
			compose := app.Get("compose")
			if compose.LastError() == nil {
				services := compose.Get("services")
				if services.LastError() == nil {
					for _, id := range services.Keys() {
						service := services.Get(id)
						if service.LastError() == nil {
							image := service.Get("image").ToString()
							if image != "" {
								images = append(images, image)
							}
						}
					}
				}
			}

			storeInfo := app.Get("store_info")
			if storeInfo.LastError() != nil {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					id,
					status,
					"n/a",
					strings.Join(images, ","),
					"(not a CasaOS compose app)",
				)
				continue
			}

			scheme := "http"
			schemeAny := storeInfo.Get("scheme")
			if schemeAny.LastError() == nil && schemeAny.ToString() != "" {
				scheme = schemeAny.ToString()
			}

			hostname, err := hostname()
			if err != nil {
				return err
			}

			hostnameAny := storeInfo.Get("hostname")
			if hostnameAny.LastError() == nil && hostnameAny.ToString() != "" {
				hostname = hostnameAny.ToString()
			}

			portMap := "unknown"
			portMapAny := storeInfo.Get("port_map")
			if portMapAny.LastError() == nil && portMapAny.ToString() != "" {
				portMap = portMapAny.ToString()
			}

			index := ""
			indexAny := storeInfo.Get("index")
			if indexAny.LastError() == nil && indexAny.ToString() != "" {
				index = indexAny.ToString()
			}

			webUI := fmt.Sprintf("%s://%s:%s/%s",
				scheme,
				hostname,
				portMap,
				strings.TrimLeft(index, "/"),
			)

			description := map[string]string{
				DefaultLanguage: "No description available",
			}

			descriptionAny := storeInfo.Get("description")
			if descriptionAny.LastError() == nil {
				for _, key := range descriptionAny.Keys() {
					description[key] = descriptionAny.Get(key).ToString()
				}
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				id,
				status,
				webUI,
				strings.Join(images, ","),
				trim(
					lo.If(
						description[DefaultLanguage] != "", description[DefaultLanguage],
					).Else(
						lo.Values(description)[0],
					),
					78,
				),
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

func composeAppStoreInfo(composeApp interface{}) (*app_management.ComposeAppStoreInfo, error) {
	composeAppMapStruct, ok := composeApp.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("app is not a map[string]interface{}")
	}

	_, ok = composeAppMapStruct["store_info"]
	if !ok {
		return nil, fmt.Errorf("app does not have \"store_info\"")
	}

	composeAppStoreInfoMapStruct, ok := composeAppMapStruct["store_info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("app[\"store_info\"] is not a map[string]interface{}")
	}

	composeAppStoreInfo := &app_management.ComposeAppStoreInfo{}

	if err := mapstructure.Decode(composeAppStoreInfoMapStruct, composeAppStoreInfo); err != nil {
		return nil, err
	}

	return composeAppStoreInfo, nil
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
