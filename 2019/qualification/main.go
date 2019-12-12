package main

import (
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

const (
	populationSize = 25
	repetition     = 100

	// Actions
	actionMaxScore = "maxScore"
	actionSend     = "send"
	actionEnd      = "end"
	actionData     = "data"

	// File paths
	// filePath = "qualification_round_2019/a_example.txt"
	// filePath = "qualification_round_2019/b_lovely_landscapes.txt"
	filePath = "qualification_round_2019/c_memorable_moments.txt"
	// filePath = "qualification_round_2019/d_pet_pictures.txt"
	// filePath = "qualification_round_2019/e_shiny_selfies.txt"
)

// Photo store imported photo information
type Photo struct {
	orientation      byte // H or V
	nrOfTag          int
	tags             map[string]struct{}
	isUsedAsVertical bool
	id               int
	used             bool
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

// TODO also score the best slide show to match with maxScore
var maxScore = 0
var r = rand.New(rand.NewSource(time.Now().Unix()))
var slideShowLength int

// TODO allow configuration of mutation rate, crossover rate etc
var mutationRate = 0.1

func main() {
	// Serve HTTP
	go ServeHTTP()
	go WriteMessage()

	select {}
}
