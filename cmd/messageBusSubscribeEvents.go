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
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// messageBusSubscribeEventsCmd represents the messageBusSubscribeEvents command
var messageBusSubscribeEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "subscribe to events in message bus",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		subscribeType, err := cmd.Flags().GetString(FlagMessageBusSubscribeType)
		if err != nil {
			log.Fatalln(err.Error())
		}

		switch subscribeType {
		case FlagMessageBusSubscribeTypeWS:

			sourceID, err := cmd.Flags().GetString(FlagMessageBusSourceID)
			if err != nil {
				log.Fatalln(err.Error())
			}

			eventNames, err := cmd.Flags().GetString(FlagMessageBusEventNames)
			if err != nil {
				log.Fatalln(err.Error())
			}

			bufferSize, err := cmd.Flags().GetUint(FlagMessageBusMessageBufferSize)
			if err != nil {
				log.Fatalln(err.Error())
			}

			subscribeWS(rootURL, "event", sourceID, eventNames, bufferSize)

		case FlagMessageBusSubscribeTypeSIO:

			subscribeSIO(rootURL, "event")

		default:
			log.Fatalf("invalid subscribe type - should be either '%s' or '%s'\n", FlagMessageBusSubscribeTypeWS, FlagMessageBusSubscribeTypeSIO)
		}
	},
}

func init() {
	messageBusSubscribeCmd.AddCommand(messageBusSubscribeEventsCmd)

	messageBusSubscribeEventsCmd.Flags().StringP(FlagMessageBusSubscribeType, "t", FlagMessageBusSubscribeTypeWS, fmt.Sprintf("subscribe type, either '%s' or '%s'", FlagMessageBusSubscribeTypeWS, FlagMessageBusSubscribeTypeSIO))

	messageBusSubscribeEventsCmd.Flags().StringP(FlagMessageBusSourceID, "s", "", "['websocket' only] source id")
	messageBusSubscribeEventsCmd.Flags().StringP(FlagMessageBusEventNames, "n", "", "['websocket' only] event names (separated by comma)")
	messageBusSubscribeEventsCmd.Flags().UintP(FlagMessageBusMessageBufferSize, "m", 1024, "['websocket' only] message buffer size")

	if err := messageBusSubscribeEventsCmd.MarkFlagRequired(FlagMessageBusSourceID); err != nil {
		log.Fatalln(err.Error())
	}
}
