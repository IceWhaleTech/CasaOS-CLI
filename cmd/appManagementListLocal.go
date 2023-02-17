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
	"log"
	"net/http"
	"text/tabwriter"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/app_management"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

// appManagementListLocalCmd represents the appManagementListLocal command
var appManagementListLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "list locally installed apps",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathAppManagement)

		client, err := app_management.NewClientWithResponses(url)
		if err != nil {
			log.Fatalln(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		response, err := client.MyComposeAppList(ctx)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Fatalln("unexpected status code", response.Status)
		}

		// get mapstruct of response body - can't unmarshal directly due to https://github.com/compose-spec/compose-go/issues/353
		dec := json.NewDecoder(response.Body)

		var body map[string]interface{}
		if err := dec.Decode(&body); err != nil {
			log.Fatalln(err.Error())
		}

		_, ok := body["data"]
		if !ok {
			log.Fatalln("unexpected response body")
		}

		data, ok := body["data"].(map[string]interface{})
		if !ok {
			log.Fatalln("unexpected response body")
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "ID\tDESCRIPTION")
		fmt.Fprintln(w, "--\t-----------")

		for id, app := range data {
			mainApp, appList, err := appList(app)
			if err != nil {
				log.Fatalln(err.Error())
			}

			fmt.Fprintf(w, "%s\t%s\n",
				id,
				appList[mainApp].Description["en_US"][0:60]+"...",
			)
		}
	},
}

func init() {
	appManagementListCmd.AddCommand(appManagementListLocalCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementListLocalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementListLocalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
