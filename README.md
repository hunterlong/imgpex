# Imgpex
This golang application accepts a list of images from a text file and creates a CSV file with URL,color1,color2,color3 scheme. 

## Focus
- Handling unlimited image URL's by reading each line and passing URL into a channel.
- Increase CPU time by not copying the image data to local machine, image is store **in memory** and released when image processing is complete.
- Save CSV data as soon as color information is fetched, writing a new line.
- Application **does not** store unique data based on image URL, `images.txt` has 40 unique images out of 1,000. But I am thinking this file could be "billions" of URLs. I am considering all URL's to be unique and/or require processing.

## Theory
Based on 1 CPU the image processing should be the most expensive task. This application will process 4 (`queue` const) images as a queue. 
It will wait for 4 images to be processed before continuing to download the others images.

## Process Queue Example
- 4 images are downloaded (queue 4)
  - 1 image completes processing (queue 3)
    - 1 image downloads (queue 4)
      - 3 images complete processing (queue 1)
        - 3 images download (queue 4)
          - 4 images complete processing
            - done
            
## Process
1. Open `images.txt` and read URL's line by line, the URL is send to the `imgData` channel.
2. Load 4 URL's and download all 4 images (`queue`)
3. Processing each image =>
    1. Get all colors in all pixels `map[hexcode]count`
    2. Get the unique pixels, sort the Count, and return the top 3 hexcodes with highest Count.
    3. Format CSV data as `url, color1, color2, color3`
    4. Append data as new row to `result.csv` CSV file.
    5. 1 image from queue is complete
4. After 1 image is processed, 1 image will be downloaded (if queue is at 3)
5. The process will end when it reaches the bottom of the text file.

## Docker Container
You can run the commands below to build this golang project and force the Docker container to use `512mb` of RAM, and 1 CPU. 
This Docker image will build from the golang source, and then insert the binary into Alpine linux.

1. `docker build -t imgpex .`
2. `docker run -it --memory="512m" --cpus=1 imgpex`

## Logs
```
Downloaded http://i.imgur.com/TKLs9lo.jpg in 0.12 seconds
Downloaded http://i.imgur.com/FApqk3D.jpg in 0.18 seconds
Downloaded https://i.redd.it/4m5yk8gjrtzy.jpg in 0.20 seconds
Downloaded https://i.redd.it/d8021b5i2moy.jpg in 0.48 seconds
Processed http://i.imgur.com/TKLs9lo.jpg in 0.65 seconds with colors: [F7F7F7 FEFEFE FFFFFF]
Downloaded https://i.redd.it/xae65ypfqycy.jpg in 0.07 seconds
Processed https://i.redd.it/4m5yk8gjrtzy.jpg in 2.87 seconds with colors: [030001 020001 010101]
Downloaded http://i.imgur.com/lcEUZHv.jpg in 0.06 seconds
Processed http://i.imgur.com/FApqk3D.jpg in 4.69 seconds with colors: [F3C300 000000 FFFFFF]
Downloaded https://i.redd.it/1nlgrn49x7ry.jpg in 0.28 seconds
Processed http://i.imgur.com/lcEUZHv.jpg in 2.15 seconds with colors: [444A5A C6BCBF FFFFFF]
```
