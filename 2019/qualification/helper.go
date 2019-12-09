package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// ServeHTTP will serve files to HTTP requests
func ServeHTTP() {
	// File server
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", HandleConnections)

	// Start http web server
	log.Println("HTTP server started on :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe:" + err.Error())
	}
}

// HandleConnections handles the web socket request
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade GET request to a web socket request
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err.Error())
	}

	// Close web socket at the end
	defer ws.Close()

	// Retrieve the client
	client = ws

	// Store message content
	m := Message{}

	// Start using web socket
	for {
		// Read action from user
		err := ws.ReadJSON(&m)

		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived) {
				log.Println("Read connection:", err)
			} else {
				log.Println("Connection closed")
			}
			break
		} else {
			// Message read success
			// Run algorithm if its not already running
			if m.Action == actionSend && m.Data.(bool) && !isAlgorithmRunning {
				fmt.Println("Starting algorithm ......")
				isAlgorithmRunning = true
				StartAlgorithm()

				// Switch off the flag
				maxScore = 0
				isAlgorithmRunning = false
			}
		}
	}
}

// WriteMessage writes to the client from a receiving channel
func WriteMessage() {
	var value Message

	for {
		value = <-broadcast

		err := client.WriteJSON(value)
		if err != nil {
			panic(err)
		}
	}
}

// ReadFile read the dataset of the problem
func ReadFile() (photos []Photo, nrOfPhotos int) {
	// Define file location
	fmt.Println("File used:", filePath)

	// Initialize ID
	var id int

	// Open file for reading
	file, err := os.Open(filePath)

	if err != nil {
		panic("Cant open file for reading" + err.Error())
	}
	defer file.Close()

	// Create a reader
	ioReader := bufio.NewReader(file)

	// Read first line: number of photos
	line, _ := ioReader.ReadString('\n')
	nrOfPhotos, _ = strconv.Atoi(strings.TrimSpace(line))

	// Read remaining lines: other photos
	for {
		line, err := ioReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err.Error())
		}

		// Process the line
		lines := strings.Split(strings.TrimSpace(line), " ")
		nrOfTag, _ := strconv.Atoi(lines[1])

		// Create the photo instance
		photo := Photo{
			orientation: lines[0][0],
			nrOfTag:     nrOfTag,
			tags:        map[string]struct{}{},
			id:          id,
		}

		// Assign tags to photo
		for _, t := range lines[2:] {
			photo.tags[t] = struct{}{}
		}

		// Add newly created photo to the photos set
		photos = append(photos, photo)

		// Increment ID for the next new photo
		id++
	}

	return
}

// AppendVerticalPhoto will merge 2 input photos into the first photo
func AppendVerticalPhoto(photo, photo1 *Photo) {
	for t := range photo1.tags {
		if _, ok := photo.tags[t]; !ok {
			photo.nrOfTag++
			photo.tags[t] = struct{}{}
		}
	}

	// Set the flag to indicate the vertical photo is used
	photo.isUsedAsVertical = true
	photo1.isUsedAsVertical = true
}

// CalcNumberOfOverlapTags sum the number of tags that overlap in the 2 input photos
func CalcNumberOfOverlapTags(photo, photo1 Photo) (total int) {
	for p1t := range photo.tags {
		if _, ok := photo1.tags[p1t]; ok {
			total++
		}
	}

	return
}

// CalcScoreBetweenTwo return the slide show transition score between the 2 input photos
func CalcScoreBetweenTwo(p1, p2 Photo) int {
	p1TagCount := p1.nrOfTag
	p2TagCount := p2.nrOfTag
	overlap := 0

	for t := range p1.tags {
		if _, ok := p2.tags[t]; ok {
			overlap++
			p1TagCount--
			p2TagCount--
		}
	}

	return Min(p1TagCount, p2TagCount, overlap)
}

// EnsureUniqueNumber will ensure the 2 input parameters will not have the same value
func EnsureUniqueNumber(x, y *int, len int) {
	if x == y {
		if *y+1 >= len {
			*y--
		} else {
			*y++
		}
	}
}

// Min returns the smallest value from the input parameter
func Min(values ...int) int {
	lowest := values[0]

	for _, i := range values[1:] {
		if i < lowest {
			lowest = i
		}
	}

	return lowest
}
