package main

import (
	"context"
	"fmt"
	"os"

	openfga "github.com/openfga/go-sdk"
	. "github.com/openfga/go-sdk/client"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Using no authentication - too lazy -_-
	fgaClient, err := NewSdkClient(&ClientConfiguration{
		ApiScheme: getEnv("FGA_API_SCHEME", "http"),         // optional, defaults to "https"
		ApiHost:   getEnv("FGA_API_HOST", "localhost:8080"), // required, define without the scheme (e.g. api.openfga.example instead of https://api.openfga.example),
	})

	if err != nil {
		panic(err)
	}

	// Create Store
	createResp, err := fgaClient.CreateStore(ctx).Body(ClientCreateStoreRequest{Name: "FGA Demo"}).Execute()
	if err != nil {
		panic(err)
	}

	fmt.Println("Created store", *createResp.Name)
	// store store.Id in database
	// update the storeId of the current instance
	fgaClient.SetStoreId(*createResp.Id)

	// Configure authorization model
	// model
	// 	schema 1.1
	// type user
	// type document
	// 	relations
	// 		define reader: [user]
	// 		define writer: [user]
	// 		define owner: [user]

	// Option 1:
	// var writeAuthorizationModelRequestString = "{\"schema_version\":\"1.1\",\"type_definitions\":[{\"type\":\"user\"},{\"type\":\"document\",\"relations\":{\"reader\":{\"this\":{}},\"writer\":{\"this\":{}},\"owner\":{\"this\":{}}},\"metadata\":{\"relations\":{\"reader\":{\"directly_related_user_types\":[{\"type\":\"user\"}]},\"writer\":{\"directly_related_user_types\":[{\"type\":\"user\"}]},\"owner\":{\"directly_related_user_types\":[{\"type\":\"user\"}]}}}}]}"
	// var body openfga.WriteAuthorizationModelRequest
	// if err := json.Unmarshal([]byte(writeAuthorizationModelRequestString), &body); err != nil {
	// 	panic(err)
	// }

	// writeAuthorModelResp, _, err := fgaClient.OpenFgaApi.WriteAuthorizationModel(context.Background()).Body(body).Execute()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Create authorization model with id ", *writeAuthorModelResp.AuthorizationModelId)

	// Option 2:
	writeAuthorReqBody := ClientWriteAuthorizationModelRequest{
		SchemaVersion: "1.1",
		TypeDefinitions: []openfga.TypeDefinition{
			{Type: "user", Relations: &map[string]openfga.Userset{}},
			{
				Type: "document",
				Relations: &map[string]openfga.Userset{
					"reader": {
						This: &map[string]interface{}{},
					},
					"writer": {
						This: &map[string]interface{}{},
					},
					"owner": {
						This: &map[string]interface{}{},
					},
				},
				Metadata: &openfga.Metadata{
					Relations: &map[string]openfga.RelationMetadata{
						"writer": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
						"reader": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
						"owner": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
					},
				},
			},
		},
	}

	writeAuthorModelResp, err := fgaClient.WriteAuthorizationModel(ctx).Body(writeAuthorReqBody).Execute()
	if err != nil {
		panic(err)
	}

	fmt.Println("Create authorization model with id ", *writeAuthorModelResp.AuthorizationModelId)

	// Add relationship tuples
	// {
	// 	user: 'user:anne',
	// 	relation: 'reader',
	// 	object: 'document:Z',
	// }
	writeOptions := ClientWriteOptions{
		AuthorizationModelId: writeAuthorModelResp.AuthorizationModelId,
	}

	writeReqBody := ClientWriteRequest{
		Writes: &[]ClientTupleKey{
			{
				User:     *openfga.PtrString("user:kiennt"),
				Relation: *openfga.PtrString("reader"),
				Object:   *openfga.PtrString("document:secret"),
			},
		},
	}

	writeResp, err := fgaClient.Write(ctx).Body(writeReqBody).Options(writeOptions).Execute()
	if err != nil {
		panic(err)
	}

	fmt.Println("Created relation response", writeResp.Writes[0].TupleKey)

	// Perform check
	checkOptions := ClientCheckOptions{
		AuthorizationModelId: writeAuthorModelResp.AuthorizationModelId,
	}

	checkReqBody := ClientCheckRequest{
		User:     "user:kiennt",
		Relation: "reader",
		Object:   "document:secret",
	}

	checkResp, err := fgaClient.Check(ctx).Body(checkReqBody).Options(checkOptions).Execute()
	if err != nil {
		panic(err)
	}

	fmt.Println("Check result", *checkResp.Allowed)
}
