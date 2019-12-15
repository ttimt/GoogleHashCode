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

type vehicle struct {
	rs []*ride
}

type ride struct {
	startRow      int // the row of the start intersection (0 ≤ a < R)
	startColumn   int // the column of the start intersection (0 ≤ b < C)
	endRow        int // the row of the finish intersection (0 ≤ x < R)
	endColumn     int // the column of the finish intersection (0 ≤ y < C)
	earliestStart int // the earliest start (0 ≤ s < T)
	latestEnd     int // the latest finish (0 ≤ f ≤ T) , (f ≥ s + |x − a| + |y − b|)
	id            int // ID to easily identify each ride

	distance     int  // Distance between end and start point
	latestStart  int  // Latest end - distance
	earliestStep int  // max of (0,0), earliest start, and distance from start to previous end plus previous steps
	startStep    int  // Real start
	endStep      int  // Real end
	isAssigned   bool // Is ride assigned to a vehicle
}

var p problem
var vs []vehicle
var rs []ride

func init() {
	p = problem{}

	readFile()
}

func main() {
	debugProblem()
	// debugVehicles()
	// debugRides()
	runAlgorithm()
	vs[0].assignRide(&rs[0])
	vs[1].assignRide(&rs[1])
	fmt.Println("Ride assigned!")
	debugVehicles()
	debugRides()
}

func runAlgorithm() {

}

func calcScore() {

}

func debugProblem() {
	fmt.Println("Problem:")
	fmt.Println("Number of rows:", p.nrRows)
	fmt.Println("Number of columns:", p.nrColumns)
	fmt.Println("Number of vehicles:", p.nrVehicles)
	fmt.Println("Number of steps:", p.nrSteps)
	fmt.Println("On time bonus for a single ride:", p.perRideOnTimeBonus)
	fmt.Println()
}

func debugVehicles() {
	fmt.Println("Vehicles:")
	for k, v := range vs {
		fmt.Println("Vehicle index:", k)
		fmt.Println("Number of rides assigned:", len(v.rs))
		fmt.Println("Last grid:", v.getLastRow(), v.getLastColumn())
		fmt.Println("Last step:", v.getLastStep())
		fmt.Println()
	}
}

func debugRides() {
	fmt.Println("Rides:")
	for _, r := range rs {
		fmt.Println("ID:", r.id)
		fmt.Println("Start grid:", r.startRow, r.startColumn)
		fmt.Println("End grid:", r.endRow, r.endColumn)
		fmt.Println("Earliest start & Latest start:", r.earliestStart, r.latestStart)
		fmt.Println("Latest end:", r.latestEnd)
		fmt.Println("Earliest step:", r.earliestStep)
		fmt.Println("Start and end step:", r.startStep, r.endStep)
		fmt.Println("Distance:", r.distance)
		fmt.Println("Is assigned:", r.isAssigned)
		fmt.Println()
	}
}

func readFile() {
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

	// Create vehicles
	createVehicles()

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
		r.updateDistance()
		r.updateLatestStart()
		r.declarativeUpdateEarliestStep()

		// Add newly created ride to the rides set
		rs = append(rs, r)

		// Increment ID for the next new ride
		id++
	}

	return
}

func createVehicles() {
	for i := 0; i < p.nrVehicles; i++ {
		newVehicle := vehicle{}
		vs = append(vs, newVehicle)
	}
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

func max(values ...int) int {
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

func getEarliestAvailableVehicle() *vehicle {
	earliestVehicle := &vs[0]

	for k, v := range vs[1:] {
		if v.getLastStep() < earliestVehicle.getLastStep() {
			earliestVehicle = &vs[k+1]
		}

		if earliestVehicle.getLastStep() == 0 {
			break
		}
	}

	return earliestVehicle
}

func declarativeUpdateAllEarliestStep() {
	for k, r := range rs {
		if !r.isAssigned {
			rs[k].declarativeUpdateEarliestStep()
		}
	}
}

func (v *vehicle) getLastStep() int {
	if len(v.rs) == 0 {
		return 0
	}

	return v.getLastRide().endStep
}

func (v *vehicle) getLastRow() int {
	if len(v.rs) == 0 {
		return 0
	}

	return v.getLastRide().endRow
}

func (v *vehicle) getLastColumn() int {
	if len(v.rs) == 0 {
		return 0
	}

	return v.getLastRide().endColumn
}

func (v *vehicle) getLastRide() *ride {
	if len(v.rs) == 0 {
		return nil
	}

	return v.rs[len(v.rs)-1]
}

func (v *vehicle) assignRide(r *ride) {
	r.startStep = max(r.earliestStart, r.startRow+r.startColumn-v.getLastRow()-v.getLastColumn()+v.getLastStep())
	r.isAssigned = true

	v.rs = append(v.rs, r)

	r.declarativeUpdateEndStep()
	declarativeUpdateAllEarliestStep()
}

func (r *ride) updateDistance() {
	r.distance = abs(r.endColumn-r.startColumn) + abs(r.endRow-r.startRow)
}

func (r *ride) updateLatestStart() {
	r.latestStart = r.latestEnd - r.distance
}

func (r *ride) declarativeUpdateEarliestStep() {
	lastVehicle := getEarliestAvailableVehicle()
	previousGrid := lastVehicle.getLastRow() + lastVehicle.getLastColumn()

	r.earliestStep = max(r.startRow+r.startColumn-previousGrid+lastVehicle.getLastStep(), r.earliestStart)
}

func (r *ride) declarativeUpdateEndStep() {
	r.endStep = r.startStep + r.distance
}

// TODO fix row and column calculation
