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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/message_bus"
	engineio "github.com/googollee/go-engine.io"
	"github.com/googollee/go-engine.io/transport"
	transportPolling "github.com/googollee/go-engine.io/transport/polling"
	transportWS "github.com/googollee/go-engine.io/transport/websocket"
	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
)

// messageBusSubscribeCmd represents the messageBusSubscribe command
var messageBusSubscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "subscribe to events/actions in message bus",
}

func init() {
	messageBusCmd.AddCommand(messageBusSubscribeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// messageBusSubscribeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// messageBusSubscribeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func subscribeWS(rootURL, entityType, sourceID, names string, bufferSize uint) {
	var wsURL string

	if names == "" {
		wsURL = fmt.Sprintf("ws://%s/%s/%s/%s", strings.TrimRight(rootURL, "/"), BasePathMessageBus, entityType, sourceID)
	} else {
		wsURL = fmt.Sprintf("ws://%s/%s/%s/%s?names=%s", strings.TrimRight(rootURL, "/"), BasePathMessageBus, entityType, sourceID, names)
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
		log.Println(string(output))
	}
}

func subscribeSIO(rootURL, entityType string) {
	dialer := engineio.Dialer{
		Transports: []transport.Transport{
			transportWS.Default,
			transportPolling.Default,
		},
	}

	sioURL := fmt.Sprintf("http://%s/%s/%s", strings.TrimRight(rootURL, "/"), BasePathMessageBus, entityType)
	conn, err := dialer.Dial(sioURL, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()

	log.Printf("subscribed to %s via socketio", sioURL)

	for {
		_, r, err := conn.NextReader()
		if err != nil {
			log.Println(err.Error())
			break
		}
		b, err := io.ReadAll(r)
		if err != nil {
			r.Close()
			log.Println(err.Error())
			break
		}
		if err := r.Close(); err != nil {
			log.Println(err.Error())
			break
		}

		fmt.Println(string(b))
	}
}
