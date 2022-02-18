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
