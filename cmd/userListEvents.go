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

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/user_service"
	"github.com/spf13/cobra"
)

// userListEventsCmd represents the userListEvents command
var userListEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "list all events received by the user",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathUsers)

		client, err := user_service.NewClientWithResponses(url)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		response, err := client.GetEventsWithResponse(ctx, &user_service.GetEventsParams{})
		if err != nil {
			return err
		}

		if response.StatusCode() != http.StatusOK {
			var baseResponse user_service.BaseResponse
			if err := json.Unmarshal(response.Body, &baseResponse); err != nil {
				return fmt.Errorf("%s - %s", response.Status(), response.Body)
			}

			return fmt.Errorf("%s - %s", response.Status(), *baseResponse.Message)
		}

		if response.JSON200 == nil || len(*response.JSON200) == 0 {
			fmt.Println("No events received")
			return nil
		}

		for _, event := range *response.JSON200 {
			fmt.Printf("%s\t%s\t%s\t%s\t%#v\n", event.Timestamp, event.EventUuid, event.SourceID, event.Name, event.Properties)
		}

		return nil
	},
}

func init() {
	userListCmd.AddCommand(userListEventsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userListEventsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userListEventsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
