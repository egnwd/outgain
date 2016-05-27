package controller

import (
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func SpriteHandler(staticDir string) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the image exists, otherwise create it
		outputPath := staticDir + r.URL.String()
		if _, err := os.Stat(outputPath); err != nil {
			// Open a reader for the specified file
			path, err := filepath.Abs(staticDir + "/images/creature.png")
			if err != nil {
				log.Println(err)
			}
			reader, err := os.Open(path)
			if err != nil {
				log.Println(err)
			}
			defer reader.Close()

			// Open a writer for the specified output
			writer, err := os.Create(outputPath)
			if err != nil {
				log.Println(err)
			}
			defer writer.Close()

			err = createSprite(reader, writer, r.URL.String())
			if err != nil {
				log.Println(err)
			}
		}
		http.ServeFile(w, r, outputPath)
	})
}

func createSprite(reader io.Reader, writer io.Writer, url string) error {
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	// Find the bounds of the image and create a new one of the same size
	bounds := img.Bounds()
	minX, minY := bounds.Min.X, bounds.Min.Y
	maxX, maxY := bounds.Max.X, bounds.Max.Y
	colouredSprite := image.NewNRGBA(bounds)
	// Used for manipulating the new image's pixel RGBA values
	stride := colouredSprite.Stride
	pixels := colouredSprite.Pix

	// Extract hex code of colour from the URL and convert to an integer
	col := strings.TrimSuffix(strings.SplitAfter(url, "-")[1], ".png")
	colVal, err := strconv.ParseUint(col, 16, 32)
	if err != nil {
		return err
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

	// Write new image to file and serve it
	png.Encode(writer, colouredSprite)
	return nil
}
