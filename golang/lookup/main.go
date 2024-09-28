package main

import (
	"fmt"
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/deluan/lookup"
)

// Helper function to load an image from the filesystem
func loadImageFromFile(imgPath string) image.Image {
	imageFile, err := os.Open(imgPath)
	if err != nil {
		panic(err)
	}
	defer imageFile.Close()
	img, _, err := image.Decode(imageFile)
	if err != nil {
		panic(err)
	}
	return img
}

func main() {
	// Create an OCR object with an accuracy of 1.0
	ocr := lookup.NewOCR(1.0)

	// Load an image to recognize
	img := loadImageFromFile("captcha.jpg")

	// Recognize text in image
	text, err := ocr.Recognize(img)
	if err != nil {
		panic(err)
	}

	// Print the results
	fmt.Printf("Text found in image: %s\n", text)
}
