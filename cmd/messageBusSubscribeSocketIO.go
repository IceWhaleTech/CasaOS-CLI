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
	"io"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/googollee/go-socket.io/parser"
	"github.com/spf13/cobra"
)

// messageBusSubscribeSocketIOCmd represents the messageBusSubscribeSocketIO command
var messageBusSubscribeSocketIOCmd = &cobra.Command{
	Use:   "socketio",
	Short: "subscribe to all entities in message bus via socketio",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		subscribeSIO(rootURL)
	},
}

func init() {
	messageBusSubscribeCmd.AddCommand(messageBusSubscribeSocketIOCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// messageBusSubscribeSocketIOCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// messageBusSubscribeSocketIOCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func subscribeSIO(rootURL string) {
	dialer := engineio.Dialer{
		Transports: []transport.Transport{
			websocket.Default,
			polling.Default,
		},
	}

	sioURL := fmt.Sprintf("http://%s/%s/socket.io", strings.TrimRight(rootURL, "/"), BasePathMessageBus)
	conn, err := dialer.Dial(sioURL, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()

	log.Printf("subscribed to %s via socketio", sioURL)

	decoder := parser.NewDecoder(conn)

	for {
		header := parser.Header{}
		name := ""
		if err := decoder.DecodeHeader(&header, &name); err != nil {
			if err == io.EOF {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			log.Fatalln(err.Error())
		}
		fmt.Printf("header: %+v, name: '%s'\n", header, name)

		values, err := decoder.DecodeArgs([]reflect.Type{
			reflect.TypeOf(map[string]interface{}{}),
		})
		if err != nil {
			if err == io.EOF {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			log.Println(err.Error())
		}
		decoder.Close()

		for _, value := range values {

			event := value.Interface().(map[string]interface{})

			output, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				log.Println(err.Error())
			}

			fmt.Println(string(output))
		}
	}
}
