package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	hugefile := flag.String("f", "/tmp/hugofile", "the path of static file to host")
	flag.Parse()

	http.Handle("/fileserver", http.FileServer(http.Dir(*directory)))
	http.Handle("/servefile", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Not handle error here (the file may not exist)
		http.ServeFile(rw, req, *hugefile)
	}))

	log.Printf("Serving on HTTP port: %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
