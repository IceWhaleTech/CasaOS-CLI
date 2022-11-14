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
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
)

// messageBusSubscribeActionsCmd represents the messageBusSubscribeActions command
var messageBusSubscribeActionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "subscribe to actions in message bus",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		sourceID, err := cmd.Flags().GetString(FlagMessageBusSourceID)
		if err != nil {
			log.Fatalln(err.Error())
		}

		actionNames, err := cmd.Flags().GetString(FlagMessageBusActionNames)
		if err != nil {
			log.Fatalln(err.Error())
		}

		var wsURL string

		if actionNames == "" {
			wsURL = fmt.Sprintf("ws://%s/%s/action/%s", strings.TrimRight(rootURL, "/"), BasePathMessageBus, sourceID)
		} else {
			wsURL = fmt.Sprintf("ws://%s/%s/action/%s?names=%s", strings.TrimRight(rootURL, "/"), BasePathMessageBus, sourceID, actionNames)
		}

		bufferSize, err := cmd.Flags().GetUint(FlagMessageBusMessageBufferSize)
		if err != nil {
			log.Fatalln(err.Error())
		}

		ws, err := websocket.Dial(wsURL, "", "http://localhost")
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer ws.Close()

		log.Println("subscribed to", wsURL)

		for {
			msg := make([]byte, bufferSize)
			n, err := ws.Read(msg)
			if err != nil {
				log.Fatalln(err.Error())
			}
			log.Println(string(msg[:n]))
		}
	},
}

func init() {
	messageBusSubscribeCmd.AddCommand(messageBusSubscribeActionsCmd)

	messageBusSubscribeActionsCmd.Flags().UintP(FlagMessageBusMessageBufferSize, "m", 1024, "message buffer size")
	messageBusSubscribeActionsCmd.Flags().StringP(FlagMessageBusSourceID, "s", "", "source id")
	messageBusSubscribeActionsCmd.Flags().StringP(FlagMessageBusActionNames, "n", "", "action names (separated by comma)")

	if err := messageBusSubscribeActionsCmd.MarkFlagRequired(FlagMessageBusSourceID); err != nil {
		log.Fatalln(err.Error())
	}
}
