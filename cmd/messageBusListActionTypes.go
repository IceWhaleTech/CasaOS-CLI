/*
Copyright © 2022 IceWhaleTech

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
	"strings"
	"text/tabwriter"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/message_bus"
	"github.com/spf13/cobra"
)

// messageBusListActionTypesCmd represents the messageBusListActionTypes command
var messageBusListActionTypesCmd = &cobra.Command{
	Use:   "action-types",
	Short: "list action types registered in message bus",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathMessageBus)

		client, err := message_bus.NewClientWithResponses(url)
		if err != nil {
			log.Fatalln(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		response, err := client.GetActionTypesWithResponse(ctx)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if response.StatusCode() != http.StatusOK {
			log.Fatalln("unexpected status code", response.Status())
		}

		if response.JSON200 == nil || len(*response.JSON200) == 0 {
			return
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "SOURCE ID\tACTION NAME\tPROPERTY TYPES")
		fmt.Fprintln(w, "---------\t----------\t--------------")

		for _, actionType := range *response.JSON200 {
			propertyTypes := make([]string, 0)
			for _, propertyType := range actionType.PropertyTypeList {
				propertyTypes = append(propertyTypes, propertyType.Name)
			}

			fmt.Fprintf(w, "%s\t%s\t{%s}\n", actionType.SourceID, actionType.Name, strings.Join(propertyTypes, ", "))
		}
	},
}

func init() {
	messageBusListCmd.AddCommand(messageBusListActionTypesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// messageBusListActionTypesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// messageBusListActionTypesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
