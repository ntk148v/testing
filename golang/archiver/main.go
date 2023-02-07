package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/mholt/archiver/v4"
)

var (
	BufferSize = 1024
	TargetPath = "/tmp/test"
)

func mkdir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

// NOT WORKED YET
func main() {
	resp, err := http.Get("https://github.com/ryanoasis/nerd-fonts/releases/download/v2.3.3/HeavyData.zip")
	if err != nil || resp.StatusCode != http.StatusOK {
		panic(err)
	}

	defer resp.Body.Close()

	format, stream, err := archiver.Identify("", resp.Body)
	if err != nil {
		panic(err)
	}

	switch archive := format.(type) {
	case archiver.Extractor:
		if err := archive.Extract(context.Background(), stream, nil, func(ctx context.Context, f archiver.File) error {
			fmt.Println(f.FileInfo.Name())

			x, err := f.Open()
			if err != nil {
				return err
			}

			if f.IsDir() {
				mkdir(path.Join(TargetPath, f.NameInArchive))
			} else {
				t, err := os.Create(path.Join(TargetPath, f.NameInArchive))
				if err != nil {
					return err
				}
				defer t.Close()
				buf := make([]byte, BufferSize)
				for {
					n, err := x.Read(buf)
					if err != nil && err != io.EOF {
						return err
					}
					if n == 0 {
						break
					}

					if _, err := t.Write(buf[:n]); err != nil {
						return err
					}
				}
				t.Sync()
			}
			return nil
		}); err != nil {
			panic(err)
		}
	case archiver.Decompressor:
		decom, _ := format.(archiver.Decompressor)
		rc, err := decom.OpenReader(stream)
		if err != nil {
			panic(err)
		}

		defer rc.Close()
	}
}
