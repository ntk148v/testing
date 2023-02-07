package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractZip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dst, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dst, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractTar(src, dst string) error {
	reader, err := os.Open(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	var tarReader *tar.Reader

	// Handle special case

	if strings.HasSuffix(src, "gz") {
		zipReader, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		tarReader = tar.NewReader(zipReader)
	} else {
		tarReader = tar.NewReader(reader)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(dst, header.Name)
		info := header.FileInfo()

		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		case tar.TypeReg:
			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}
		case tar.TypeLink:
			link := filepath.Join(dst, header.Name)
			linkTarget := filepath.Join(dst, header.Linkname)
			// lazy link creation. just to make sure all files are available
			defer os.Link(link, linkTarget)
		case tar.TypeSymlink:
			linkPath := filepath.Join(dst, header.Name)
			if err := os.Symlink(header.Linkname, linkPath); err != nil {
				if !os.IsExist(err) {
					return err
				}
			}
		}
	}
	return nil
}

func extract(src, dst string) error {
	if strings.HasSuffix(src, ".zip") {
		return extractZip(src, dst)
	} else if strings.HasSuffix(src, ".gz") || strings.HasSuffix(src, ".tar") {
		return extractTar(src, dst)
	}

	return errors.New("unsupported archive type")
}

func main() {
	if err := extract("/tmp/test.zip", "/tmp/test"); err != nil {
		panic(err)
	}
}
