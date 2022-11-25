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

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/message_bus"
	"github.com/spf13/cobra"
)

// messageBusTriggerActionCmd represents the messageBusTriggerAction command
var messageBusTriggerActionCmd = &cobra.Command{
	Use:   "trigger",
	Short: "trigger an action via message bus",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		sourceID, err := cmd.Flags().GetString(FlagMessageBusSourceID)
		if err != nil {
			log.Fatalln(err.Error())
		}

		actionName, err := cmd.Flags().GetString(FlagMessageBusActionName)
		if err != nil {
			log.Fatalln(err.Error())
		}

		properties, err := cmd.Flags().GetString(FlagMessageBusProperties)
		if err != nil {
			log.Fatalln(err.Error())
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathMessageBus)

		client, err := message_bus.NewClientWithResponses(url)
		if err != nil {
			log.Fatalln(err.Error())
		}

		request := map[string]string{}

		for _, property := range strings.Split(properties, ",") {
			kv := strings.Split(property, "=")
			if len(kv) != 2 {
				log.Fatalln("invalid property:", property)
			}

			request[kv[0]] = kv[1]
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		response, err := client.TriggerActionWithResponse(ctx, sourceID, actionName, request)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if response == nil {
			log.Fatalln("empty response")
		}

		if response.StatusCode() != http.StatusOK {
			log.Fatalln("unexpected status code", response.Status())
		}
	},
}

func init() {
	messageBusCmd.AddCommand(messageBusTriggerActionCmd)

	messageBusTriggerActionCmd.Flags().StringP(FlagMessageBusSourceID, "s", "", "source id")
	messageBusTriggerActionCmd.Flags().StringP(FlagMessageBusActionName, "n", "", "action name")
	messageBusTriggerActionCmd.Flags().StringP(FlagMessageBusProperties, "p", "", "action properties (in form of `K=V` and separated by comma)")

	if err := messageBusTriggerActionCmd.MarkFlagRequired(FlagMessageBusSourceID); err != nil {
		log.Fatalln(err.Error())
	}

	if err := messageBusTriggerActionCmd.MarkFlagRequired(FlagMessageBusActionName); err != nil {
		log.Fatalln(err.Error())
	}
}
