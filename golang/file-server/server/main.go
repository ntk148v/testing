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
