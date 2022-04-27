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
	"fmt"
	"log"
)

func main() {
	sizeN := 5
	lookupSizeM := uint64(13) //(must be prime number)

	var names []string
	for i := 0; i < sizeN; i++ {
		names = append(names, fmt.Sprintf("backend-%d", i))
	}
	//backend-0 ~ backend-4

	mm, err := NewMaglev(names, lookupSizeM)
	if err != nil {
		log.Fatal("Init maglev", err)
	}
	v, _ := mm.Get("IP1")
	log.Println("node1:", v)
	//node1: backend-2
	v, _ = mm.Get("IP2")
	log.Println("node2:", v)
	//node2: backend-1
	v, _ = mm.Get("IPasdasdwni2")
	log.Println("node3:", v)
	//node3: backend-0

	if err := mm.Remove("backend-0"); err != nil {
		log.Fatal("Remove failed", err)
	}
	v, _ = mm.Get("IPasdasdwni2")
	log.Println("node3-D:", v)
	//node3-D: Change from "backend-0" to "backend-1"
}
