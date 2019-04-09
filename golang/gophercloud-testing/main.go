package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/gophercloud/gophercloud/pagination"
)

func main() {
	// Option 1: Pass in the values yourself
	// opts := gophercloud.AuthOptions{
	// 	IdentityEndpoint: "https://openstack.example.com:5000/v2.0",
	// 	Username:         "admin",
	// 	Password:         "{password}",
	// }
	// Option 2: Use a utility function to retrieve all your environment variables
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println(err)
		return
	}
	client, err := openstack.NewOrchestrationV1(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	listopts := stacks.ListOpts{
		SortKey: "stack_name",
		Tags:    "scale",
	}

	results := make(map[string]interface{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i <= 10; i++ {
			pager := stacks.List(client, listopts)
			fmt.Println(pager)
			err = pager.EachPage(func(page pagination.Page) (bool, error) {
				stackList, err := stacks.ExtractStacks(page)
				if err != nil {
					return false, err
				}
				var stack stacks.GetResult
				for _, s := range stackList {
					stack = stacks.Get(client, s.Name, s.ID)
					stackBody, _ := stack.Extract()
					outputValues := make(map[string]interface{})
					if len(stackBody.Outputs) == 0 {
						continue
					}
					// Convert output value (JSON string) to Map
					for _, v := range stackBody.Outputs {
						outputValueMap := make(map[string]interface{})
						outputValueRaw := v["output_value"].(string)
						if err := json.Unmarshal([]byte(outputValueRaw), &outputValueMap); err != nil {
							continue
						}
						outputValues[v["output_key"].(string)] = outputValueMap
					}
					results[s.ID] = outputValues
				}
				return true, nil
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, value := range results {
				fmt.Println(value)
			}
			time.Sleep(time.Second * 2)
		}
	}()

	wg.Wait()
	// pager := stacks.List(client, listopts)
	// fmt.Println(pager)
	// err = pager.EachPage(func(page pagination.Page) (bool, error) {
	// 	stackList, err := stacks.ExtractStacks(page)
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	var stack stacks.GetResult
	// 	for _, s := range stackList {
	// 		stack = stacks.Get(client, s.Name, s.ID)
	// 		stackBody, _ := stack.Extract()
	// 		outputValues := make(map[string]interface{})
	// 		if len(stackBody.Outputs) == 0 {
	// 			continue
	// 		}
	// 		// Convert output value (JSON string) to Map
	// 		for _, v := range stackBody.Outputs {
	// 			outputValueMap := make(map[string]interface{})
	// 			outputValueRaw := v["output_value"].(string)
	// 			if err := json.Unmarshal([]byte(outputValueRaw), &outputValueMap); err != nil {
	// 				continue
	// 			}
	// 			outputValues[v["output_key"].(string)] = outputValueMap
	// 		}
	// 		results[s.ID] = outputValues
	// 	}
	// 	return true, nil
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// for _, value := range results {
	// 	fmt.Println(value)
	// }
}
