package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var (
	uniq    map[string]int
	testImg *ImageData
)

func init() {
	testImg = &ImageData{
		Url:        "https://graf1x.com/wp-content/uploads/2017/11/tiger-orange-color-paint-code-swatch-chart-rgb-html-hex.png",
		Downloaded: make(chan bool),
		Processed:  make(chan bool),
	}
}

func TestDownloadImage(t *testing.T) {
	file, err := downloadImage(testImg)
	testImg.Image = file
	assert.Nil(t, err)
	d, err := ioutil.ReadAll(file)
	assert.Nil(t, err)
	assert.Equal(t, 14589, len(d))
}

func TestGetPixels(t *testing.T) {
	file, err := downloadImage(testImg)
	testImg.Image = file
	assert.Nil(t, err)
	uniq, err = getPixels(testImg)
	assert.Nil(t, err)
	assert.Equal(t, 18, len(uniq))
}

func TestReturnTopPixels(t *testing.T) {
	top := returnTopColors(uniq)
	testImg.Hexcodes = top
	assert.Equal(t, 3, len(top))
	assert.Equal(t, "FE9042", top[0].Key)
	assert.Equal(t, 1075, top[0].Count)
}

func BenchmarkDownloadImage(b *testing.B) {
	for n := 0; n < b.N; n++ {
		downloadImage(testImg)
	}
}

func benchmarkPixels(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getPixels(testImg)
	}
}

func BenchmarkTopPixels(b *testing.B)      { benchmarkPixels(b) }
func BenchmarkTopPixels100(b *testing.B)   { benchmarkPixels(b) }
func BenchmarkTopPixels10000(b *testing.B) { benchmarkPixels(b) }
