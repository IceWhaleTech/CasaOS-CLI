//go:generate bash -c "mkdir -p codegen/app_management && go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -generate types,client -package app_management https://raw.githubusercontent.com/IceWhaleTech/CasaOS-AppManagement/main/api/app_management/openapi.yaml > codegen/app_management/api.go"
//go:generate bash -c "mkdir -p codegen/casaos && go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -generate types,client -package casaos https://raw.githubusercontent.com/IceWhaleTech/CasaOS/main/api/casaos/openapi.yaml > codegen/casaos/api.go"
//go:generate bash -c "mkdir -p codegen/local_storage && go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -generate types,client -package local_storage https://raw.githubusercontent.com/IceWhaleTech/CasaOS-LocalStorage/main/api/local_storage/openapi.yaml > codegen/local_storage/api.go"
//go:generate bash -c "mkdir -p codegen/message_bus && go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -generate types,client -package message_bus https://raw.githubusercontent.com/IceWhaleTech/CasaOS-MessageBus/main/api/message_bus/openapi.yaml > codegen/message_bus/api.go"
//go:generate bash -c "mkdir -p codegen/user_service && go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -generate types,client -package user_service https://raw.githubusercontent.com/IceWhaleTech/CasaOS-UserService/main/api/user-service/openapi.yaml > codegen/user_service/api.go"

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
package main

import (
	"github.com/IceWhaleTech/CasaOS-CLI/cmd"
)

var (
	// see https://goreleaser.com/cookbooks/using-main.version
	version = "0.4.4"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date

	cmd.Execute()
}
