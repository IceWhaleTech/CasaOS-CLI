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

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/app_management"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const (
	FlagAppManagementCategory   = "category"
	FlagAppManagementAuthorType = "type"
	FlagAppManagementRecommend  = "recommend"
)

var authorTypes = []string{
	string(app_management.Official), string(app_management.ByCasaos), string(app_management.Community),
}

// appManagementSearchCmd represents the appManagementSearch command
var appManagementSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "search for apps in app store",
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

		params := &app_management.ComposeAppStoreInfoListParams{}

		category, err := cmd.Flags().GetString(FlagAppManagementCategory)
		if err != nil {
			return err
		}

		if category != "" {
			params.Category = &category
		}

		authorType, err := cmd.Flags().GetString(FlagAppManagementAuthorType)
		if err != nil {
			return err
		}

		if authorType != "" {
			if !lo.Contains(authorTypes, authorType) {
				return fmt.Errorf("invalid author type %s, should be one of %s", authorType, strings.Join(authorTypes, ", "))
			}

			params.AuthorType = (*app_management.StoreAppAuthorType)(&authorType)
		}

		recommend, err := cmd.Flags().GetBool(FlagAppManagementRecommend)
		if err != nil {
			return err
		}

		if recommend {
			params.Recommend = &recommend
		}

		response, err := client.ComposeAppStoreInfoListWithResponse(ctx, params)
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

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		defer w.Flush()

		if response.JSON200.Data.List == nil || len(*response.JSON200.Data.List) == 0 {
			return fmt.Errorf("no compose app found from this store")
		}

		installedList := []string{}
		if response.JSON200.Data.Installed != nil && len(*response.JSON200.Data.Installed) != 0 {
			installedList = *response.JSON200.Data.Installed
		}

		fmt.Fprintln(w, "Name\tCategory\tAuthor\tDeveloper\tDescription")
		fmt.Fprintln(w, "----\t--------\t------\t---------\t-----------")

		for storeAppID, composeApp := range *response.JSON200.Data.List {
			if lo.Contains(installedList, storeAppID) {
				storeAppID = fmt.Sprintf("%s [installed]", storeAppID)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", storeAppID, composeApp.Category, composeApp.Author, composeApp.Developer, trim(composeApp.Description[DefaultLanguage], 78))
		}

		return nil
	},
}

func init() {
	appManagementCmd.AddCommand(appManagementSearchCmd)

	appManagementSearchCmd.Flags().StringP(FlagAppManagementAuthorType, "t", "", fmt.Sprintf("author type of the app (%s)", strings.Join(authorTypes, ", ")))

	appManagementSearchCmd.Flags().StringP(FlagAppManagementCategory, "c", "", "category of the app")

	appManagementSearchCmd.Flags().BoolP(FlagAppManagementRecommend, "r", false, "recommend apps")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementSearchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementSearchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
