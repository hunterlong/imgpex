package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"sort"
)

// getPixels accepts a local filename and will return a map[string]int as map[hexcode]count
func getPixels(d *ImageData) (map[string]int, error) {
	defer d.Image.Close()
	img, _, err := image.Decode(d.Image)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	uniqueMap := make(map[string]int)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			hex := toHex(img.At(x, y).RGBA())
			uniqueMap[hex]++
		}
	}
	return uniqueMap, nil
}

// returnTopColors takes a map[hex]count of colors in the image and returns the top 3 colors
func returnTopColors(colors map[string]int) Hexcodes {
	unique := make(Hexcodes, len(colors))
	var index int
	for k, v := range colors {
		unique[index] = Hexcode{k, v}
		index++
	}
	sort.Sort(unique)
	return unique[len(unique)-3:]
}

// toHex converts a pixel RGB to hexcode
func toHex(r, g, b, a uint32) string {
	return fmt.Sprintf("%.2X%.2X%.2X", int(r>>8), int(g>>8), int(b>>8))
}
