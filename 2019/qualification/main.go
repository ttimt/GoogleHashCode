package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"sync"
	"time"
)

// My result: 1,111,595
// 2019 Qualification Position: #19
// A - 2        - all under 1 minute
// B - 205,689  - all under 1 minute
// C - 1,627    - all under 1 minute
// D - 443,825  - 2 minutes assigning verticals (Score: 229,448) - 14 minutes creating slide show
// E - 460,452  - 19 minutes assigning verticals (Score: 178,982) - 2 hrs 50 minutes on creating slide show
//
// Best known score: 1,243,566
// A - 2
// B - 239,997
// C - 1,775
// D - 441,242
// E - 560,550

const (
	populationSize = 15
	repetition     = 100

	// Actions
	actionMaxScore = "maxScore"
	actionSend     = "send"
	actionEnd      = "end"
	actionData     = "data"

	// File paths
	filePathPrepend = "C:\\Users\\Timothy\\go\\src\\github.com\\ttimt\\GoogleHashCode\\2019\\qualification\\qualification_round_2019\\"
	// filePathPrepend = "C:\\Users\\CLM6\\go\\src\\github.com\\ttimt\\GoogleHashCode\\2019\\qualification\\qualification_round_2019\\"
	filePathA = "a_example.txt"
	filePathB = "b_lovely_landscapes.txt"
	filePathC = "c_memorable_moments.txt"
	filePathD = "d_pet_pictures.txt"
	filePathE = "e_shiny_selfies.txt"
)

// Photo store imported photo information
type Photo struct {
	orientation      byte // H or V
	nrOfTag          int
	tags             map[string]struct{}
	isUsedAsVertical bool
	id               int
	used             bool
	currentScore     int
}

// Result store the result used to send to the UI graph
type Result struct {
	X string `json:"x"` // HH::mm:ss
	Y int    `json:"y"` // the score
}

// Message is a wrapper data struct for communication through web socket
type Message struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

var client *websocket.Conn
var broadcast = make(chan Message)
var isAlgorithmRunning = false
var wg sync.WaitGroup

var maxScore = 0
var r = rand.New(rand.NewSource(time.Now().Unix()))
var mutationRate = 0.2

func main() {
	// Serve HTTP
	// go ServeHTTP()
	// go WriteMessage()
	//
	// select {}

	wg.Add(2)
	fmt.Println(time.Now().Format(time.Kitchen))
	go startTagAlgorithm(filePathA)
	// go startTagAlgorithm(filePathB)
	go startTagAlgorithm(filePathC)
	// go startTagAlgorithm(filePathD)
	// go startTagAlgorithm(filePathE)
	wg.Wait()
}
