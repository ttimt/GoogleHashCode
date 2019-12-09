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

func ServeHTTP() {
	// File server
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", HandleConnections)

	// Start http web server
	log.Println("http server started on :8081")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe:" + err.Error())
	}
}

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

	// Start using web socket
	m := Message{}

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
			if m.Action == actionSend && m.Data.(bool) && !isAlgorithmRunning {
				fmt.Println("starting algorithm")
				isAlgorithmRunning = true
				StartAlgorithm()

				// Switch off the flag
				isAlgorithmRunning = false
				maxScore = 0
			}
		}
	}
}

func WriteMessage() {
	for {
		value := <-broadcast

		err := client.WriteJSON(value)
		if err != nil {
			panic(err)
		}
	}
}

func ReadFile() (photos []Photo, nrOfPhotos int) {
	// Define file location
	fmt.Println("File used:", filePath)

	// Initialize ID
	var id int

	// Read line
	file, err := os.Open(filePath)

	if err != nil {
		panic("Cant open file for reading")
	}
	defer file.Close()
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

		// Create the instance
		photo := Photo{
			orientation: byte(lines[0][0]),
			nrOfTag:     nrOfTag,
			tags:        map[string]struct{}{},
			id:          id,
		}

		// Assign tags to photo
		for _, t := range lines[2:] {
			photo.tags[t] = struct{}{}
		}

		photos = append(photos, photo)

		id++
	}

	return
}

func AppendVerticalPhoto(photo, photo1 *Photo) {
	for t := range photo1.tags {
		if _, ok := photo.tags[t]; !ok {
			photo.nrOfTag++
			photo.tags[t] = struct{}{}
		}
	}
}

func EnsureUniqueNumber(x, y *int, len int) {
	if x == y {
		if *y+1 >= len {
			*y--
		} else {
			*y++
		}
	}
}

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

	return min(p1TagCount, p2TagCount, overlap)
}

func min(values ...int) int {
	lowest := values[0]

	for _, i := range values[1:] {
		if i < lowest {
			lowest = i
		}
	}

	return lowest
}
