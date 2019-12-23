package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Result: 702,864‬
// 2019 Qualification Position: #873
// A - 2
// B - 205,623
// C - 1,627
// D - 292,762
// E - 202,850
//
// Best score:
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
	filePathA = "C:\\Users\\Timothy\\go\\src\\github.com\\ttimt\\GoogleHashCode\\2019\\qualification\\qualification_round_2019/a_example.txt"
	filePathB = "C:\\Users\\Timothy\\go\\src\\github.com\\ttimt\\GoogleHashCode\\2019\\qualification\\qualification_round_2019/b_lovely_landscapes.txt"
	filePathC = "C:\\Users\\Timothy\\go\\src\\github.com\\ttimt\\GoogleHashCode\\2019\\qualification\\qualification_round_2019/c_memorable_moments.txt"
	filePathD = "C:\\Users\\Timothy\\go\\src\\github.com\\ttimt\\GoogleHashCode\\2019\\qualification\\qualification_round_2019/d_pet_pictures.txt"
	filePathE = "C:\\Users\\Timothy\\go\\src\\github.com\\ttimt\\GoogleHashCode\\2019\\qualification\\qualification_round_2019/e_shiny_selfies.txt"
	filePathF = "C:\\Users\\Timothy\\go\\src\\github.com\\ttimt\\GoogleHashCode\\2019\\qualification\\qualification_round_2019/f.txt"
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

// TODO also score the best slide show to match with maxScore
var maxScore = 0
var r = rand.New(rand.NewSource(time.Now().Unix()))

// TODO allow configuration of mutation rate, crossover rate etc
var mutationRate = 0.2

func main() {
	// Serve HTTP
	// go ServeHTTP()
	// go WriteMessage()
	//
	// select {}

	wg.Add(1)
	// go startTagAlgorithm(filePathA)
	// go startTagAlgorithm(filePathB)
	// go startTagAlgorithm(filePathC)
	go StartAlgorithm(filePathD)
	// go startCategoryAlgorithm(filePathE)

	// go startTagAlgorithm(filePathF)
	wg.Wait()
}
