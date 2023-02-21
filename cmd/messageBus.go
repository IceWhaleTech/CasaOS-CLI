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
	"github.com/spf13/cobra"
)

// messageBusCmd represents the messageBus command
var messageBusCmd = &cobra.Command{
	Use:     "message-bus",
	Short:   "All message bus related commands",
	GroupID: RootGroupID,
}

const (
	BasePathMessageBus = "v2/message_bus"

	FlagMessageBusActionName        = "action-name"
	FlagMessageBusActionNames       = "action-names"
	FlagMessageBusEventNames        = "event-names"
	FlagMessageBusMessageBufferSize = "message-buffer-size"
	FlagMessageBusProperties        = "properties"
	FlagMessageBusSourceID          = "source-id"
)

func init() {
	rootCmd.AddCommand(messageBusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// messageBusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// messageBusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
