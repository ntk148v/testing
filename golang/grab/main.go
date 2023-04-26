package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

func main() {
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", "https://github.com/ryanoasis/nerd-fonts/releases/download/v2.3.3/VictorMono.zip")

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf(" %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("\r  transferred %v / %v bytes (%.2f%%)",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress())
		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)
}
