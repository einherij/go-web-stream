package streamvideo

import (
	"bytes"
	"log"
	"time"

	"gocv.io/x/gocv"
)

const (
	// default frame update time in millisecond
	defaultFrameUpdateTime = 50
)

var (
	// Frames contains channel with video frames
	Frames = make(chan []byte)
)

// StartReadingVideo is main function that gets frames from camera and sends it to the channel
func StartReadingVideo(frameUpdateTime time.Duration) {
	if frameUpdateTime == 0 {
		frameUpdateTime = defaultFrameUpdateTime
	}
	// Select camera device
	camera, err := gocv.VideoCaptureDevice(0)
	defer camera.Close()
	if err != nil {
		log.Fatal("!Can't find camera device. Error:", err)
	}

	// Create basic image container gocv
	img := gocv.NewMat()

	// Camera reading loop
	for {
		// read new frame and put in 'img'
		camera.Read(&img)
		if img.Empty() {
			continue
		}

		// convert Mat container to jpg bytes
		imgBytes, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			log.Fatal("!Can't encode image. Got error:", err)
		}

		// add header and footer for image (to use as body for response)
		bimg := [][]byte{[]byte("--frame\r\nContent-Type: image/jpeg\r\n\r\n"), imgBytes.GetBytes(), []byte("\r\n\r\n")}
		imgBytes.Close()

		// join header, image, footer and send to channel
		Frames <- bytes.Join(bimg, []byte(""))

		// wait between frames
		time.Sleep(frameUpdateTime)
	}
}
