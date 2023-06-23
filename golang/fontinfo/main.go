package main

import (
	"fmt"

	"github.com/liamg/fontinfo"
)

func main() {

	fonts, err := fontinfo.List()
	// fonts, err := fontinfo.Match(fontinfo.MatchFamily("Iosevka"))
	if err != nil {
		panic(err)
	}

	for _, font := range fonts {
		fmt.Printf("Family=%s Style=%s Path=%s\n", font.Family, font.Style, font.Path)
	}
}
