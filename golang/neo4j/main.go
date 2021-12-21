package main

import (
	"fmt"
	"os"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func main() {
	// Neo4j 4.0, defaults to no TLS therefore use bolt:// or neo4j://
	// Neo4j 3.5, defaults to self-signed certificates, TLS on, therefore use bolt+ssc:// or neo4j+ssc://
	driver, err := neo4j.NewDriver(os.Getenv("NEO4J_URL"), neo4j.BasicAuth(os.Getenv("NEO4J_USERNAME"), os.Getenv("NEO4J_PASSWORD"), ""))
	if err != nil {
		panic(err)
	}
	// Handle driver lifetime based on your application lifetime requirements  driver's lifetime is usually
	// bound by the application lifetime, which usually implies one driver instance per application
	defer driver.Close()
	item, err := insertItem(driver)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", item)
}

func insertItem(driver neo4j.Driver) (*Item, error) {
	// Sessions are short-lived, cheap to create and NOT thread safe. Typically create one or more sessions
	// per request in your web application. Make sure to call Close on the session when done.
	// For multi-database support, set sessionConfig.DatabaseName to requested database
	// Session config will default to write mode, if only reads are to be used configure session for
	// read mode.
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	result, err := session.WriteTransaction(createItemFn)
	if err != nil {
		return nil, err
	}
	return result.(*Item), nil
}

func createItemFn(tx neo4j.Transaction) (interface{}, error) {
	records, err := tx.Run("CREATE (n:Item { id: $id, name: $name }) RETURN n.id, n.name", map[string]interface{}{
		"id":   1,
		"name": "Item 1",
	})
	// In face of driver native errors, make sure to return them directly.
	// Depending on the error, the driver may try to execute the function again.
	if err != nil {
		return nil, err
	}
	record, err := records.Single()
	if err != nil {
		return nil, err
	}
	// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
	return &Item{
		Id:   record.Values[0].(int64),
		Name: record.Values[1].(string),
	}, nil
}

type Item struct {
	Id   int64
	Name string
}
