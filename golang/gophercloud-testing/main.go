package main

import (
	"fmt"
	"os"

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
	// stack := stacks.Find(client, "vsmart_kiennt_2")
	// stackBody, _ := stack.Extract()
	// // Convert output value (JSON string) to Map
	// outputValueMap := make(map[string]interface{})
	// outputValueRaw := stackBody.Outputs[0]["output_value"].(string)
	// err = json.Unmarshal([]byte(outputValueRaw), &outputValueMap)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(outputValueMap["service_name"])
	listopts := stacks.ListOpts{
		TenantID: "c5a8b5960ac04cc68f18a541a7a9c51e",
	}

	pager := stacks.List(client, listopts)
	fmt.Println(pager)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		stackList, err := stacks.ExtractStacks(page)
		if err != nil {
			return false, err
		}
		for _, s := range stackList {
			fmt.Println(s)
		}
		return true, nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
