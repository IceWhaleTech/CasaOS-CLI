//go:generate bash -c "mkdir -p codegen/message_bus && go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.2 -package message_bus https://raw.githubusercontent.com/IceWhaleTech/CasaOS-MessageBus/add_action_type/api/message_bus/openapi.yaml > codegen/message_bus/api.go"
// TODO update OpenAPI URL above when release ^^^

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

import "github.com/IceWhaleTech/CasaOS-CLI/cmd"

func main() {
	cmd.Execute()
}
