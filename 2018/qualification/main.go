package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Score : 45,643,519 - Position #706
// A - 10 - Perfect
// B - 176,877 - Perfect
// C - 15,770,067 - 53,562 more
// D - 8,230,620 - 4,061,925‬ more
// E - 21,465,945 - Perfect
//
// Top #2 score - 49759006
// C - 15,823,629
// D - 12,292,545

// My latest step end should be earlier than the other rides latest start
const (
	// filePath = "qualification_round_2018.in/a_example.in"
	// filePath = "qualification_round_2018.in/b_should_be_easy.in"
	// filePath = "qualification_round_2018.in/c_no_hurry.in"
	filePath = "qualification_round_2018.in/d_metropolis.in"
	// filePath = "qualification_round_2018.in/e_high_bonus.in"
)

const (
	NW = iota
	NE
	SE
	SW
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
	id int
}

type ride struct {
	startRow      int // the row of the start intersection (0 ≤ a < R)
	startColumn   int // the column of the start intersection (0 ≤ b < C)
	endRow        int // the row of the finish intersection (0 ≤ x < R)
	endColumn     int // the column of the finish intersection (0 ≤ y < C)
	earliestStart int // the earliest start (0 ≤ s < T)
	latestEnd     int // the latest finish (0 ≤ f ≤ T) , (f ≥ s + |x − a| + |y − b|)
	id            int // ID to easily identify each ride

	distance        int      // Distance between end and start point
	latestStart     int      // Latest end - distance
	earliestEnd     int      // Earliest start + distance
	earliestStep    int      // max of (0,0), earliest start, and distance from start to previous end plus previous steps
	startStep       int      // Real start
	endStep         int      // Real end
	isAssigned      bool     // Is ride assigned to a vehicle
	earliestVehicle *vehicle // Store the vehicle used in earliest step calculation
	assignedVehicle *vehicle // The vehicle if the ride is assigned
}

var p problem
var vs []*vehicle
var rs []*ride
var earliestStepRide *ride
var score int

func init() {
	p = problem{}

	readFile()
}

func main() {
	debugProblem()
	// debugVehicles()
	// debugRides()
	runAlgorithm()
	fmt.Println("Ride assigned!")
	// printResult()
	calcScore()
	fmt.Println("Total score:", score)
	// printUnassignedRides()
	// processUnassignedRides()
	frequencyAnalysis()
}

func runAlgorithm() {
	for i := 0; i < len(rs); i++ {
		if earliestStepRide.isAssigned || earliestStepRide.earliestStep+earliestStepRide.distance > p.nrSteps {
			break
		}

		earliestStepRide.earliestVehicle.assignRide(earliestStepRide)
	}
}

func calcScore() {
	for kv := range vs {
		for kr := range vs[kv].rs {
			score += vs[kv].rs[kr].distance

			// fmt.Print(strconv.Itoa(vs[kv].rs[kr].distance) + " ")
			if vs[kv].rs[kr].startStep == vs[kv].rs[kr].earliestStart {
				score += p.perRideOnTimeBonus
			}
		}
	}
}

func printResult() {
	for kv := range vs {
		fmt.Print("Vehicle: ", vs[kv].id, " - Rides: ")
		for kr := range vs[kv].rs {
			fmt.Print(" ", vs[kv].rs[kr].id)
		}
		fmt.Println()
	}
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
	for _, v := range vs {
		fmt.Println("Vehicle ID:", v.id)
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
		fmt.Println("Earliest end & Latest end:", r.earliestEnd, r.latestEnd)
		fmt.Println("Earliest step & Latest step:", r.earliestStep, r.earliestStep+r.distance)
		fmt.Println("Start and end step:", r.startStep, r.endStep)
		fmt.Println("Distance:", r.distance)
		fmt.Println("Is assigned:", r.isAssigned)
		if r.isAssigned {
			fmt.Println("Assigned vehicle:", r.assignedVehicle.id)
		} else {
			fmt.Println("Earliest vehicle:", r.earliestVehicle.id)
		}
		fmt.Println()
	}
}

func printUnassignedRides() {
	fmt.Println("Unassigned rides:")

	for k := range rs {
		if !rs[k].isAssigned {
			fmt.Println("Ride ID:", rs[k].id, " - distance: ", rs[k].distance, coordinateLocation(rs[k].endRow, rs[k].endColumn))
		}
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
		r := &ride{
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
		r.updateEarliestEnd()
		r.declarativeUpdateEarliestStep()

		// Add newly created ride to the rides set
		rs = append(rs, r)

		// Increment ID for the next new ride
		id++
	}

	return
}

func createVehicles() {
	for i := 1; i <= p.nrVehicles; i++ {
		newVehicle := &vehicle{
			id: i,
		}
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

func calculateDistance(startRow, startColumn, endRow, endColumn int) int {
	return abs(endColumn-startColumn) + abs(endRow-startRow)
}

func getEarliestAvailableVehicle() *vehicle {
	earliestVehicle := vs[0]

	for k, v := range vs[1:] {
		if v.getLastStep() < earliestVehicle.getLastStep() {
			earliestVehicle = vs[k+1]
		}

		if earliestVehicle.getLastStep() == 0 {
			break
		}
	}

	return earliestVehicle
}

func declarativeUpdateAllEarliestStep() {
	for k := range rs {
		if !rs[k].isAssigned {
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
	r.isAssigned = true
	r.assignedVehicle = v
	r.updateStartStep(v)

	v.rs = append(v.rs, r)

	r.declarativeUpdateEndStep()
	declarativeUpdateAllEarliestStep()
}

func (r *ride) updateDistance() {
	r.distance = calculateDistance(r.startRow, r.startColumn, r.endRow, r.endColumn)
}

func (r *ride) updateLatestStart() {
	r.latestStart = r.latestEnd - r.distance
}

func (r *ride) updateEarliestEnd() {
	r.earliestEnd = r.earliestStart + r.distance
}

func (r *ride) updateStartStep(v *vehicle) {
	r.startStep = max(r.earliestStart, calculateDistance(r.startRow, r.startColumn, v.getLastRow(), v.getLastColumn())+v.getLastStep())
}

func (r *ride) declarativeUpdateEarliestStep() {
	lastVehicle := getEarliestAvailableVehicle()
	r.earliestVehicle = lastVehicle

	r.earliestStep = max(calculateDistance(r.startRow, r.startColumn, lastVehicle.getLastRow(), lastVehicle.getLastColumn())+lastVehicle.getLastStep(), r.earliestStart)

	// Skip
	if r.earliestStep+r.distance > r.latestEnd { // || coordinateLocation(r.endRow, r.endColumn) != SW { + 2 millions points
		return
	}

	if earliestStepRide == nil || earliestStepRide.isAssigned || r.earliestStep < earliestStepRide.earliestStep {
		earliestStepRide = r
	}
}

func (r *ride) declarativeUpdateEndStep() {
	r.endStep = r.startStep + r.distance
}

func processUnassignedRides() {
	for k := range rs {
		if !rs[k].isAssigned {
			// Try to add this ride
		}
	}
}

func frequencyAnalysis() {
	var northWest, northEast, southEast, southWest int

	for k := range rs {
		switch coordinateLocation(rs[k].startRow, rs[k].startColumn) {
		case NW:
			northWest++
		case NE:
			northEast++
		case SE:
			southEast++
		case SW:
			southWest++
		}
	}

	fmt.Println("Frequency analysis")
	fmt.Println("Start frequency")
	fmt.Println("North west:", float64(northWest)/float64(p.nrRides))
	fmt.Println("North east:", float64(northEast)/float64(p.nrRides))
	fmt.Println("South east:", float64(southEast)/float64(p.nrRides))
	fmt.Println("South west:", float64(southWest)/float64(p.nrRides))

	northEast = 0
	northWest = 0
	southWest = 0
	southEast = 0

	for k := range rs {
		switch coordinateLocation(rs[k].endRow, rs[k].endColumn) {
		case NW:
			northWest++
		case NE:
			northEast++
		case SE:
			southEast++
		case SW:
			southWest++
		}
	}

	fmt.Println("End frequency")
	fmt.Println("North west:", float64(northWest)/float64(p.nrRides))
	fmt.Println("North east:", float64(northEast)/float64(p.nrRides))
	fmt.Println("South east:", float64(southEast)/float64(p.nrRides))
	fmt.Println("South west:", float64(southWest)/float64(p.nrRides))
}

func coordinateLocation(row, column int) int {
	location := 0

	if row <= p.nrRows/2 {
		// North
		if column <= p.nrColumns/2 {
			location = SW
		} else {
			location = SE
		}
	} else {
		// South
		if column <= p.nrColumns/2 {
			location = NW
		} else {
			location = NE
		}
	}

	return location
}
