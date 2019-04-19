package main

import (
	"io"
	"time"
)

// ImageData contains the data for a single Image
type ImageData struct {
	Hexcodes   Hexcodes
	Image      io.ReadCloser
	Url        string
	Duration   time.Duration
	Downloaded chan bool
	Processed  chan bool
}

// Hexcode contains the hexidecimal code for the color and the amount in the image
type Hexcode struct {
	Key   string
	Count int
}

// Hexcodes is a slice of Hexcode used for sorting
type Hexcodes []Hexcode

func (p Hexcodes) Len() int           { return len(p) }
func (p Hexcodes) Less(i, j int) bool { return p[i].Count < p[j].Count }
func (p Hexcodes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
