package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Province struct {
	Name      string     `json:"name"`
	Code      int        `json:"code"`
	Districts []District `json:"districts"`
}

type District struct {
	Name  string `json:"name"`
	Code  int    `json:"code"`
	Wards []Ward `json:"wards"`
}

type Ward struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}

func readCSVFile(csvFile string) ([][]string, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	return r.ReadAll()

}

func main() {
	// Download Excel from https://danhmuchanhchinh.gso.gov.vn/
	// Convert to csv (cause I want to work with csv)
	records, err := readCSVFile("vn-provinces.csv")
	if err != nil {
		panic(err)
	}

	for _, r := range records[1:] {
		fmt.Println(r)
	}
}
