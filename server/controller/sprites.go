package controller

import (
  "os"
  "path/filepath"
  "image"
  "image/png"
  "log"
  "net/http"
)

func SpriteHandler(w http.ResponseWriter, r *http.Request) {
  log.Printf("got here")
  // Load the image from the file path, converted for OS
  path, err := filepath.Abs("client/images/sprite.png")
  // Generate output image with corresponding id
  if err != nil {
    log.Fatal(err)
  }
  reader, err := os.Open(path)
  if err != nil {
    log.Fatal(err)
  }
  defer reader.Close()
  // Convert the IO reader to an image
  img, _, err := image.Decode(reader)
  if err != nil {
    log.Fatal(err)
  }
  // Find the bounds of the image and create a new one of the same size
  bounds := img.Bounds()
  minX, minY := bounds.Min.X, bounds.Min.Y
  maxX, maxY := bounds.Max.X, bounds.Max.Y
  colouredSprite := image.NewNRGBA(image.Rect(minX, minY, maxX, maxY))
  // Used for manipulating the new image's pixel RGBA values
  stride := colouredSprite.Stride
  pixels := colouredSprite.Pix
  // Iterate through pixels and colourise them
  for i := minX; i < maxX; i++ {
    for j := minY; j < maxY; j++ {
    // Get RGBA value from original image
    // Colour values are 8 bit, extended to 32, so conversion is safe
    r32, g32, b32, a32 := img.At(i, j).RGBA()
    r := uint8(r32)
    g := uint8(g32)
    b := uint8(b32)
    a := uint8(a32)
    current := (j - minY) * stride + (i - minX) * 4
    pixels[current] = 255 - r
    pixels[current + 1] = 255 - g
    pixels[current + 2] = 255 - b
    pixels[current + 3] = a
    }
  }
  // Write image to file
  outputPath := r.URL.String()
  log.Printf(outputPath)
  writer, err := os.Create(outputPath)
  if err != nil {
    log.Fatal(err)
  }
  defer writer.Close()
  png.Encode(writer, colouredSprite)
  http.ServeFile(w, r, outputPath)
}
