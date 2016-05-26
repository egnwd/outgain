package controller

import (
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func SpriteHandler(w http.ResponseWriter, r *http.Request, staticDir string) {
	// Load the image from the file path, converted for OS
	path, err := filepath.Abs("client/images/sprite.png")
	// Generate output image with corresponding id
	if err != nil {
		log.Println(err)
	}
	reader, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer reader.Close()
	// Convert the IO reader to an image
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Println(err)
	}
	// Find the bounds of the image and create a new one of the same size
	bounds := img.Bounds()
	minX, minY := bounds.Min.X, bounds.Min.Y
	maxX, maxY := bounds.Max.X, bounds.Max.Y
	colouredSprite := image.NewNRGBA(image.Rect(minX, minY, maxX, maxY))
	// Used for manipulating the new image's pixel RGBA values
	stride := colouredSprite.Stride
	pixels := colouredSprite.Pix

	// Extract hex code of colour from the URL and convert to an integer
	col := strings.TrimSuffix(strings.SplitAfter(r.URL.String(), "-")[1], ".png")
	colVal, err := strconv.ParseUint(col, 16, 32)
	if err != nil {
		log.Println(err)
	}
	// Get individual RGB values from hex colour code
	r1 := uint8((colVal & 0xff0000) >> 16)
	g1 := uint8((colVal & 0x00ff00) >> 8)
	b1 := uint8(colVal & 0x0000ff)

	// Iterate through pixels and colourise them
	for i := minX; i < maxX; i++ {
		for j := minY; j < maxY; j++ {
			// Get RGBA value from original image
			// Colour values are 8 bit, extended to 32, so conversion is safe
			r32, g32, b32, a32 := img.At(i, j).RGBA()
			r0 := uint8(r32)
			g0 := uint8(g32)
			b0 := uint8(b32)
			a0 := uint8(a32)
			current := (j-minY)*stride + (i-minX)*4
			// Set new pixel values, retaining white border
			if r0 == g0 && r0 == b0 && r0 == 0xff {
				pixels[current] = r0
				pixels[current+1] = g0
				pixels[current+2] = b0
				pixels[current+3] = a0
			} else {
				// Uncomment to make colours more saturated
				/*
				   if (r1 >= g1 && r1 >= b1) {
				     r1 = 0xff
				   }
				   if (g1 >= r1 && g1 >= b1) {
				     g1 = 0xff
				   }
				   if (b1 >= r1 && b1 >= g1) {
				     b1 = 0xff
				   }
				*/
				pixels[current] = r0/2 + r1/2
				pixels[current+1] = g0/2 + g1/2
				pixels[current+2] = b0/2 + b1/2
				pixels[current+3] = a0
			}
		}
	}
	// Write image to file
	outputPath := staticDir + r.URL.String()
	log.Printf(outputPath)
	writer, err := os.Create(outputPath)
	if err != nil {
		log.Println(err)
	}
	defer writer.Close()
	png.Encode(writer, colouredSprite)
	http.ServeFile(w, r, outputPath)
}
