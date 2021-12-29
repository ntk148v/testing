package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func main() {
	// Neo4j 4.0, defaults to no TLS therefore use bolt:// or neo4j://
	// Neo4j 3.5, defaults to self-signed certificates, TLS on, therefore use bolt+ssc:// or neo4j+ssc://
	driver, err := neo4j.NewDriver(os.Getenv("NEO4J_URI"), neo4j.BasicAuth(os.Getenv("NEO4J_USERNAME"), os.Getenv("NEO4J_PASSWORD"), ""))
	if err != nil {
		panic(err)
	}
	// Handle driver lifetime based on your application lifetime requirements  driver's lifetime is usually
	// bound by the application lifetime, which usually implies one driver instance per application
	defer driver.Close()

	// Load cypher commands from files
	buf, err := ioutil.ReadFile(os.Getenv("NEO4J_CYPHER_FILE"))
	if err != nil {
		panic(err)
	}
	cypher := string(buf)

	// Create session
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	for _, query := range strings.SplitAfter(cypher, ";") {
		if query == "\n" || query == "" {
			continue
		}
		fmt.Println("Execute query", query)
		if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return tx.Run(query, nil)
		}); err != nil {
			panic(err)
		}
	}
}
