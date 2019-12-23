package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Result: 496,783
// 2019 Qualification Position: #1524
// A - 2
// B - 12
// C - 1,157
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
	populationSize = 25
	repetition     = 1000

	// Actions
	actionMaxScore = "maxScore"
	actionSend     = "send"
	actionEnd      = "end"
	actionData     = "data"

	// File paths
	filePathA = "qualification_round_2019/a_example.txt"
	filePathB = "qualification_round_2019/b_lovely_landscapes.txt"
	filePathC = "qualification_round_2019/c_memorable_moments.txt"
	filePathD = "qualification_round_2019/d_pet_pictures.txt"
	filePathE = "qualification_round_2019/e_shiny_selfies.txt"
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

	wg.Add(5)
	go startCategoryAlgorithm(filePathA)
	go startCategoryAlgorithm(filePathB)
	go startCategoryAlgorithm(filePathC)
	go startCategoryAlgorithm(filePathD)
	go startCategoryAlgorithm(filePathE)
	wg.Wait()
}
