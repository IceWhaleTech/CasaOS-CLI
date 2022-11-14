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

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/local_storage"
	"github.com/IceWhaleTech/CasaOS-Common/utils"
	"github.com/spf13/cobra"
)

// localStorageSetMergeCmd represents the localStorageSetMerge command
var localStorageSetMergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "set a merge in local storage",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathLocalStorage)

		fsType, err := cmd.Flags().GetString(FlagLocalStorageFSType)
		if err != nil {
			log.Fatalln(err.Error())
		}

		mountPoint, err := cmd.Flags().GetString(FlagLocalStorageMountPoint)
		if err != nil {
			log.Fatalln(err.Error())
		}

		sourceBasePath, err := cmd.Flags().GetString(FlagLocalStorageSourceBasePath)
		if err != nil {
			log.Fatalln(err.Error())
		}

		sourceVolumeUUIDs, err := cmd.Flags().GetString(FlagLocalStorageSourceVolumeUUIDs)
		if err != nil {
			log.Fatalln(err.Error())
		}

		client, err := local_storage.NewClientWithResponses(url)
		if err != nil {
			log.Fatalln(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		request := local_storage.Merge{
			Fstype:     &fsType,
			MountPoint: mountPoint,
		}

		if sourceBasePath != "" {
			request.SourceBasePath = &sourceBasePath
		}

		if sourceVolumeUUIDs != "" {
			request.SourceVolumeUuids = utils.Ptr(strings.Split(sourceVolumeUUIDs, ","))
		}

		response, err := client.SetMergeWithResponse(ctx, request)
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
	localStorageSetCmd.AddCommand(localStorageSetMergeCmd)

	localStorageSetMergeCmd.Flags().StringP(FlagLocalStorageFSType, "t", "fuse.mergerfs", "merge type")
	localStorageSetMergeCmd.Flags().StringP(FlagLocalStorageMountPoint, "m", "", "mount point")
	localStorageSetMergeCmd.Flags().String(FlagLocalStorageSourceBasePath, "", "source base path")
	localStorageSetMergeCmd.Flags().String(FlagLocalStorageSourceVolumeUUIDs, "", "source volume uuids (separated by comma)")

	if err := localStorageSetMergeCmd.MarkFlagRequired(FlagLocalStorageMountPoint); err != nil {
		log.Fatalln(err.Error())
	}
}
