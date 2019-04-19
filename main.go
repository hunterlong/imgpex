package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	imagesList     = "images.txt"
	resultFilename = "result.csv"
	queue          = 4
)

var (
	imgData    chan *ImageData
	csvWriter  *csv.Writer
	complete   chan bool
	concurrent chan struct{}
	wg         *sync.WaitGroup
)

func init() {
	imgData = make(chan *ImageData)
	concurrent = make(chan struct{}, queue)
	complete = make(chan bool)
	wg = new(sync.WaitGroup)
}

func handleError(err error) {
	log.Println("Error: ", err)
	os.Exit(1)
}

func main() {
	csvFile, err := os.Create(resultFilename)
	if err != nil {
		handleError(err)
	}
	csvWriter = csv.NewWriter(csvFile)
	defer csvFile.Close()

	go loadTextFile(imagesList, imgData)

	log.Println("Fetching top colors for each line.")

	// run a loop and break the loop once images.txt has reach the last line
Renderer:
	for {
		select {
		case d := <-imgData:
			concurrent <- struct{}{}
			go ProcessImage(d)
		case <-complete:
			break Renderer
		default:
			continue
		}
	}

}

// ProcessImage is the main routine to process an image in the follow steps:
// // 1. download image
// // 2. get all color for all pixels in hexcode format
// // 3. get unique colors, sort and return top 3
// // 4. format CSV data
// // 5. save CSV data
func ProcessImage(data *ImageData) {
	go data.Download()

	<-data.Downloaded

	go data.Process()

	<-data.Processed
	wg.Done()
	<-concurrent
}

func (data *ImageData) Download() {
	t1 := time.Now()
	file, err := downloadImage(data)
	if err != nil {
		handleError(err)
	}
	t2 := time.Now().Sub(t1)
	data.Duration = t2
	data.Image = file
	data.Downloaded <- true
	log.Printf("Downloaded %v in %.2f seconds\n", data.Url, data.Duration.Seconds())
}

func (data *ImageData) Process() {
	t1 := time.Now()
	// returns a map[string]int as map[hexcode]count of all pixels in image
	pixels, err := getPixels(data)
	if err != nil {
		handleError(err)
	}

	// returns the top 3 colors in hexcode
	uniqueColors := returnTopColors(pixels)
	data.Hexcodes = uniqueColors

	// format data into CSV
	csvData := formatToCsv(data.Url, uniqueColors)

	// add new row with data into CSV file
	csvWriter.Write(csvData[0:])
	csvWriter.Flush()

	speed := time.Now().Sub(t1)
	log.Printf("Processed %v in %.2f seconds with colors: [%v %v %v]\n", data.Url, speed.Seconds(), data.Hexcodes[0].Key, data.Hexcodes[1].Key, data.Hexcodes[2].Key)

	data.Image = nil
	data.Processed <- true
}

// formatToCsv will take the URL of file and top 3 hexcodes and convert into a CSV [4]string
func formatToCsv(file string, colors []Hexcode) [4]string {
	var data [4]string
	data[0] = file
	for k, v := range colors {
		data[k+1] = v.Key
	}
	return data
}

// loadTextFile loads the image list file into memory
func loadTextFile(filename string, data chan *ImageData) {
	file, err := os.Open(filename)
	if err != nil {
		handleError(err)
	}
	defer func() {
		file.Close()
		wg.Wait()
		complete <- true
	}()
	bf := bufio.NewReader(file)
	for {
		line, _, err := bf.ReadLine()
		if err == io.EOF {
			break
		}
		wg.Add(1)
		data <- &ImageData{
			Url:        string(line),
			Downloaded: make(chan bool),
			Processed:  make(chan bool),
		}
	}
}

// downloadImage will download an image from URL and return the response body
func downloadImage(data *ImageData) (io.ReadCloser, error) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Get(data.Url)
	if err != nil {
		return nil, err
	}
	defer client.CloseIdleConnections()
	return response.Body, err
}
