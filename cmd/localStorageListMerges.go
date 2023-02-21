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

	"github.com/spf13/cobra"

	"github.com/IceWhaleTech/CasaOS-CLI/codegen/local_storage"
)

// localStorageListMergesCmd represents the localStorageListMerges command
var localStorageListMergesCmd = &cobra.Command{
	Use:   "merges",
	Short: "list merges in local storage",
	Run: func(cmd *cobra.Command, args []string) {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			log.Fatalln(err.Error())
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathLocalStorage)

		client, err := local_storage.NewClientWithResponses(url)
		if err != nil {
			log.Fatalln(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		response, err := client.GetMergesWithResponse(ctx, &local_storage.GetMergesParams{})
		if err != nil {
			log.Fatalln(err.Error())
		}

		if response.StatusCode() != http.StatusOK {
			log.Fatalln("unexpected status code", response.Status())
		}

		if response.JSON200.Data == nil || len(*response.JSON200.Data) == 0 {
			return
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "FSTYPE\tMOUNT_POINT\tSOURCE_BASE_PATH\tSOURCE_VOLUME_UUIDS\tCREATED_AT\tUPDATED_AT")
		fmt.Fprintln(w, "------\t-----------\t----------------\t-------------------\t----------\t----------")

		for _, merge := range *response.JSON200.Data {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				*merge.Fstype,
				merge.MountPoint,
				*merge.SourceBasePath,
				strings.Join(*merge.SourceVolumeUuids, ","),
				merge.CreatedAt,
				merge.UpdatedAt,
			)
		}
	},
}

func init() {
	localStorageListCmd.AddCommand(localStorageListMergesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localStorageListMergesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localStorageListMergesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
