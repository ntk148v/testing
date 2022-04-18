// Copyright 2022 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"go.i3wm.org/i3/v4"
)

func main() {
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		panic(err)
	}

	for _, workspace := range workspaces {
		fmt.Printf("%+v\n", workspace)
	}

	// Get bar
	barIDs, err := i3.GetBarIDs()
	if err != nil {
		panic(err)
	}
	for _, barID := range barIDs {
		barCfg, err := i3.GetBarConfig(barID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%v\n", barCfg)
	}

	// Get event
	subscription := i3.Subscribe(i3.WorkspaceEventType)
	for subscription.Next() {
		event := subscription.Event()
		fmt.Printf("%+v\n", event)
	}
}
