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
	// filePath = "qualification_round_2018.in/a_example.in"
	// filePath = "qualification_round_2018.in/b_should_be_easy.in"
	// filePath = "qualification_round_2018.in/c_no_hurry.in"
	// filePath = "qualification_round_2018.in/d_metropolis.in"
	filePath = "qualification_round_2018.in/e_high_bonus.in"
)

type problem struct {
	nrRows             int
	nrColumns          int
	nrVehicles         int
	nrRides            int
	nrSteps            int
	perRideOnTimeBonus int
}

type ride struct {
	startRow      int
	startColumn   int
	endRow        int
	endColumn     int
	earliestStart int
	latestEnd     int
	id            int
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

		// Add newly created ride to the rides set
		rs = append(rs, r)

		// Increment ID for the next new ride
		id++
	}

	return
}
