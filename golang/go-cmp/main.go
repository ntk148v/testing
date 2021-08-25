// Copyright 2021 Kien Nguyen-Tuan <kiennt2609@gmail.com>
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

	"github.com/google/go-cmp/cmp"
)

type Person struct {
	Name    string
	Age     int
	Complex Complex
}

type Complex struct {
	Field1 map[string]int
	Field2 []string
}

func main() {
	m := make(map[string]int)
	m["k1"] = 7
	m["k2"] = 8
	a := Person{
		Name: "Tui",
		Age:  26,
		Complex: Complex{
			Field1: m,
			Field2: []string{"test"},
		},
	}
	b := Person{
		Name: "Nguoikhac",
		Age:  26,
	}
	if diff := cmp.Diff(a, b); diff != "" {
		fmt.Printf("%+v\n", diff)
	}
}
