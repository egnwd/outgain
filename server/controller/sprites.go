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

	"github.com/ajstarks/svgo"
	"github.com/gorilla/mux"
)

const creaturePath = `M59.5,28.2c-1.4,0-2.7,1-3.4,2h-2.8c-0.4-4.9-2.2-8.4-4.7-11.6l2.2-2.2c1.2,0.2,2.5-0.2,3.4-1.1
	c1.6-1.6,1.6-4,0-5.5c-1.6-1.6-4-1.6-5.5,0c-1,1-1.4,2.4-1,3.7l-2.2,2c-3.3-2.7-6.7-4.3-11.6-4.8v-3c1-0.7,1.9-1.9,1.9-3.3
	c0-2.2-1.7-3.9-3.8-3.9s-3.7,1.8-3.7,3.9c0,1.4,0.9,2.6,1.9,3.3v3c-4.9,0.4-8.4,2.2-11.6,4.7l-2.1-2.1c0.3-1.3-0.1-2.7-1.1-3.7
	c-1.6-1.6-4-1.6-5.5,0c-1.6,1.6-1.6,4,0,5.5c0.9,0.9,2.2,1.3,3.4,1.1l2.2,2.3c-2.7,3.3-4.3,6.7-4.8,11.6H7.6c-0.7-1-1.9-2-3.4-2
	c-2.2,0-3.9,1.7-3.9,3.8s1.8,3.7,3.9,3.7c1.3,0,2.4-0.7,3.2-1.7h3c0.4,4.9,2.2,8.5,4.7,11.6l-1.9,1.9c-1.3-0.3-2.7,0.1-3.7,1.1
	c-1.6,1.6-1.6,4,0,5.5c1.6,1.6,4,1.6,5.5,0c0.9-0.9,1.3-2.2,1.1-3.4l2.1-2c3.3,2.7,6.7,4.3,11.6,4.8v2.7c-1,0.7-1.9,1.9-1.9,3.3
	c0,2.2,1.7,3.9,3.8,3.9s3.7-1.8,3.7-3.9c0-1.4-0.9-2.6-1.9-3.3v-2.7c4.9-0.4,8.4-2.2,11.6-4.7l2.1,2c-0.2,1.2,0.2,2.5,1.1,3.4
	c1.6,1.6,4,1.6,5.5,0c1.6-1.6,1.6-4,0-5.5c-1-1-2.4-1.3-3.7-1l-1.9-2c2.7-3.3,4.3-6.7,4.8-11.6h3c0.7,1,1.9,1.7,3.2,1.7
	c2.2,0,3.9-1.7,3.9-3.8C63.5,29.8,61.7,28.2,59.5,28.2z`

func SVGSpriteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	const size = 64
	vars := mux.Vars(r)
	colour := vars["colour"]

	creature := svg.New(w)
	creature.Start(size, size)
	style := "fill:#" + colour
	creature.Path(creaturePath, style)
	creature.End()
}

func SpriteHandler(staticDir string) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the image exists, otherwise create it
		outputPath := staticDir + r.URL.String()
		if _, err := os.Stat(outputPath); err != nil {
			// Open a reader for the specified file and read the image
			path, err := filepath.Abs(staticDir + "/images/creature.png")
			if err != nil {
				log.Println(err)
				return
			}
			reader, err := os.Open(path)
			if err != nil {
				log.Println(err)
				return
			}
			defer reader.Close()
			img, _, err := image.Decode(reader)
			if err != nil {
				log.Println(err)
				return
			}

			// Generate the new image
			newImg, err := createSprite(img, r.URL.String())
			if err != nil {
				log.Println(err)
				return
			}

			// Open a writer for the specified output and write the new image
			writer, err := os.Create(outputPath)
			if err != nil {
				log.Println(err)
				return
			}
			defer writer.Close()
			png.Encode(writer, newImg)

		}
		http.ServeFile(w, r, outputPath)
	})
}

func createSprite(img image.Image, url string) (image.Image, error) {

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
		return nil, err
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

	return colouredSprite, nil
}
