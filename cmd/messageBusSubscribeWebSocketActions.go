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
	"log"

	"github.com/spf13/cobra"
)

// messageBusSubscribeWebSocketActionsCmd represents the messageBusSubscribeActions command
var messageBusSubscribeWebSocketActionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "subscribe to actions in message bus via websocket",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		sourceID, err := messageBusSubscribeWebSocketCmd.PersistentFlags().GetString(FlagMessageBusSourceID)
		if err != nil {
			log.Fatalln(err.Error())
		}

		bufferSize, err := messageBusSubscribeWebSocketCmd.PersistentFlags().GetUint(FlagMessageBusMessageBufferSize)
		if err != nil {
			log.Fatalln(err.Error())
		}

		actionNames, err := cmd.Flags().GetString(FlagMessageBusActionNames)
		if err != nil {
			log.Fatalln(err.Error())
		}

		subscribeWS(rootURL, "action", sourceID, actionNames, bufferSize)
	},
}

func init() {
	messageBusSubscribeWebSocketCmd.AddCommand(messageBusSubscribeWebSocketActionsCmd)

	messageBusSubscribeWebSocketActionsCmd.Flags().StringP(FlagMessageBusActionNames, "n", "", "action names (separated by comma)")
}
