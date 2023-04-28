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
	"fmt"
	"log"
	"strings"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/message_bus"
	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
)

// messageBusSubscribeWebSocketCmd represents the messageBusSubscribeWebSocket command
var messageBusSubscribeWebSocketCmd = &cobra.Command{
	Use:   "websocket",
	Short: "subscribe to all entities in message bus via websocket",
}

func init() {
	messageBusSubscribeCmd.AddCommand(messageBusSubscribeWebSocketCmd)

	messageBusSubscribeWebSocketCmd.PersistentFlags().StringP(FlagMessageBusSourceID, "s", "", "source id")
	messageBusSubscribeWebSocketCmd.PersistentFlags().UintP(FlagMessageBusMessageBufferSize, "m", 1024, "message buffer size")

	if err := messageBusSubscribeWebSocketCmd.MarkPersistentFlagRequired(FlagMessageBusSourceID); err != nil {
		log.Fatalln(err.Error())
	}
}

func subscribeWS(rootURL, messageType, sourceID, names string, bufferSize uint) {
	var wsURL string

	if names == "" {
		wsURL = fmt.Sprintf("ws://%s/%s/%s/%s", strings.TrimRight(rootURL, "/"), BasePathMessageBus, messageType, sourceID)
	} else {
		wsURL = fmt.Sprintf("ws://%s/%s/%s/%s?names=%s", strings.TrimRight(rootURL, "/"), BasePathMessageBus, messageType, sourceID, names)
	}

	ws, err := websocket.Dial(wsURL, "", "http://localhost")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer ws.Close()

	log.Printf("subscribed to %s via websocket", wsURL)

	for {
		msg := make([]byte, bufferSize)
		n, err := ws.Read(msg)
		if err != nil {
			log.Fatalln(err.Error())
		}

		var event message_bus.Event

		if err := json.Unmarshal(msg[:n], &event); err != nil {
			log.Println(err.Error())
		}

		output, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			log.Println(err.Error())
		}

		fmt.Println(string(output))
	}
}
