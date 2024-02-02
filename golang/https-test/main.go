package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Max-Age", "3600")

		fmt.Fprint(w, "Hello world")
	})

	// if err := http.ListenAndServe("127.0.0.1:8084", nil); err != nil {
	// 	panic(err)
	// }

	if err := http.ListenAndServeTLS("127.0.0.1:8085", "server.crt", "server.key", nil); err != nil {
		panic(err)
	}
}
