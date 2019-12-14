package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	filePath = "qualification_round_2018.in/a_example.in"
	// filePath = "qualification_round_2018.in/b_should_be_easy.in"
	// filePath = "qualification_round_2018.in/c_no_hurry.in"
	// filePath = "qualification_round_2018.in/d_metropolis.in"
	// filePath = "qualification_round_2018.in/e_high_bonus.in"

	maxNr = 10000
)

type problem struct {
	nrRows             int // number of rows of the grid (1<=R<=10000)
	nrColumns          int // number of columns of the grid (1<=C<=10000)
	nrVehicles         int // number of vehicles in the fleet (1<=F<=1000)
	nrRides            int // number of rides (1<=N<=10000)
	nrSteps            int // number of steps in the simulation (1<=T<=10^9)
	perRideOnTimeBonus int // per-ride bonus for starting the ride on time (1<=B<=10000)
}

type ride struct {
	startRow      int // the row of the start intersection (0 ≤ a < R)
	startColumn   int // the column of the start intersection (0 ≤ b < C)
	endRow        int // the row of the finish intersection (0 ≤ x < R)
	endColumn     int // the column of the finish intersection (0 ≤ y < C)
	earliestStart int // the earliest start (0 ≤ s < T)
	latestEnd     int // the latest finish (0 ≤ f ≤ T) , (f ≥ s + |x − a| + |y − b|)
	id            int // ID to easily identify each ride

	latestStart int // Latest end - distance
}

var p problem
var rs []ride

func init() {
	p = problem{}

	ReadFile()
}

func main() {
	fmt.Println("Problem:", p)
	fmt.Println("Rides:", rs)

	runAlgorithm()
}

func runAlgorithm() {

}

func calcScore() {

}

// ReadFile read the dataset of the problem
func ReadFile() {
	// Define file location
	fmt.Println("File used:", filePath)

	// Open file for reading
	file, err := os.Open(filePath)

	if err != nil {
		panic("Cant open file for reading" + err.Error())
	}
	defer file.Close()

	// Create a reader
	ioReader := bufio.NewReader(file)

	// Read first line: problem configuration
	line, _ := ioReader.ReadString('\n')
	lines := strings.Split(line, " ")
	p.nrRows, _ = strconv.Atoi(strings.TrimSpace(lines[0]))
	p.nrColumns, _ = strconv.Atoi(strings.TrimSpace(lines[1]))
	p.nrVehicles, _ = strconv.Atoi(strings.TrimSpace(lines[2]))
	p.nrRides, _ = strconv.Atoi(strings.TrimSpace(lines[3]))
	p.perRideOnTimeBonus, _ = strconv.Atoi(strings.TrimSpace(lines[4]))
	p.nrSteps, _ = strconv.Atoi(strings.TrimSpace(lines[5]))

	// Store ID
	id := 0

	// Read remaining lines: rides
	for {
		line, err := ioReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err.Error())
		}

		// Process the line
		lines = strings.Split(strings.TrimSpace(line), " ")

		startRow, _ := strconv.Atoi(lines[0])
		startColumn, _ := strconv.Atoi(lines[1])
		endRow, _ := strconv.Atoi(lines[2])
		endColumn, _ := strconv.Atoi(lines[3])
		earliestStart, _ := strconv.Atoi(lines[4])
		latestEnd, _ := strconv.Atoi(lines[5])

		// Create the ride instance
		r := ride{
			startRow:      startRow,
			startColumn:   startColumn,
			endRow:        endRow,
			endColumn:     endColumn,
			earliestStart: earliestStart,
			latestEnd:     latestEnd,
			id:            id,
		}

		// Update declarative logic
		r.updateLatestStart()

		// Add newly created ride to the rides set
		rs = append(rs, r)

		// Increment ID for the next new ride
		id++
	}

	return
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

// Max returns the highest value from the input parameter
func Max(values ...int) int {
	highest := values[0]

	for _, i := range values[1:] {
		if i > highest {
			highest = i
		}
	}

	return highest
}

func abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}

func (r ride) updateLatestStart() {
	r.latestStart = r.latestEnd - abs(r.endRow-r.startRow) - abs(r.endColumn-r.startColumn)
}
