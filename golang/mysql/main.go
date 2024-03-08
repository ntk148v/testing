package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:secret@tcp(localhost:3306)/")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("SELECT product_id, vendor_id, instance_uuid FROM nova.pci_devices")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()

	var (
		product_id    []byte
		vendor_id     []byte
		instance_uuid []byte
	)

	for rows.Next() {
		err = rows.Scan(&product_id, &vendor_id, &instance_uuid)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(product_id), string(vendor_id), string(instance_uuid))
	}
}
