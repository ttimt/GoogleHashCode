package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

const (
	populationSize = 100
	repetition     = 1000
	// TODO allow configuration of mutation rate, crossover rate etc
	mutationRate = 0.01

	// Actions
	actionMaxScore = "maxScore"
	actionSend     = "send"
	actionEnd      = "end"
	actionData     = "data"

	// File paths
	// filePath = "qualification_round_2019/a_example.txt"
	// filePath       = "qualification_round_2019/b_lovely_landscapes.txt"
	// filePath = "qualification_round_2019/c_memorable_moments.txt"
	// filePath = "qualification_round_2019/d_pet_pictures.txt"
	filePath = "qualification_round_2019/e_shiny_selfies.txt"
)

// Photo store imported photo information
type Photo struct {
	orientation      byte // H or V
	nrOfTag          int
	tags             map[string]struct{}
	isUsedAsVertical bool
	id               int
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
var maxScore = 0

func main() {
	// Serve HTTP
	go ServeHTTP()
	go WriteMessage()

	select {}
}

// StartAlgorithm will execute the whole algorithm
func StartAlgorithm() {
	// Read file
	fmt.Println("Importing ......")
	photos, nrOfPhotos := ReadFile()
	fmt.Println("Number of photos:", nrOfPhotos)

	// Assign vertical
	fmt.Println("Assigning vertical photos ......")
	photos = AssignVertical(photos)

	// Genetic algorithm
	fmt.Println("Running algorithm ......")
	photos = GeneticAlgorithm(photos, rand.New(rand.NewSource(time.Now().Unix())), repetition)

	// Final score
	fmt.Println("Final score:")
	fmt.Println(CalcScore(photos))

	// Notify the UI that algorithm has ended
	m := Message{
		Action: actionEnd,
		Data:   true,
	}
	broadcast <- m
}

// AssignVertical assign vertical photos
func AssignVertical(photos []Photo) (answer []Photo) {
	var singleVertical []Photo

	// Filter out all horizontal images
	for _, p := range photos {
		if p.orientation != 'V' {
			answer = append(answer, p)
		} else {
			singleVertical = append(singleVertical, p)
		}
	}

	// Return if no vertical image or only one
	if len(singleVertical) <= 1 {
		return answer
	}

	// Process all vertical images
	for k, v := range singleVertical {
		// Process current image
		// If no vertical image left, discard this image from use
		if !v.isUsedAsVertical && k+1 < len(singleVertical) {
			// Pick another vertical photo to form a slide
			smallestOverlapPhotoPosition := k + 1
			smallestOverlap := CalcNumberOfOverlapTags(v, singleVertical[k+1])

			// Find the one with least overlap
			for j, v1 := range singleVertical[k+1:] {
				if !v1.isUsedAsVertical && CalcNumberOfOverlapTags(v, v1) < smallestOverlap {
					smallestOverlapPhotoPosition = k + 1 + j
					smallestOverlap = CalcNumberOfOverlapTags(v, v1)
				}
			}

			// Append smallestOverlapPhoto into v
			AppendVerticalPhoto(&singleVertical[k], &singleVertical[smallestOverlapPhotoPosition])
			answer = append(answer, singleVertical[k])
		}
	}

	return
}

// GeneticAlgorithm is the optimization algorithm
func GeneticAlgorithm(slideShow []Photo, r *rand.Rand, repetition int) []Photo {
	// 1. Generate population / a set of slide shows
	// 2. Pick the fittest
	// 3. Create an offspring from the fittest and a random slide show
	// 4. Mutation to the offspring
	// 5. Repeat by adding the mutated offspring to the population set in (1)

	// 1. Generate a population / a set of slide shows
	fmt.Println("1.0 Repetition:", repetition)

	// Use to store the population and store the original slide show in the set
	var set [][]Photo
	set = append(set, slideShow)

	// Store the random number of mutation to occur per slide show in a set
	var numberOfMutation, firstPhotoPosition, secondPhotoPosition int

	// Length of slide show
	lenSlideShow := len(slideShow)

	// Store the temporary photo when swapping
	var swap Photo

	// fmt.Println("1.1 Repetition:", repetition)
	// TODO verify random swap should have at least 1 score in between slide shows
	//  or else attempt a few more times then only swap with photo of lowest number of tag
	for i := 0; i < populationSize; i++ {
		// Store the new instance of slide show
		newSlideShow := make([]Photo, lenSlideShow)
		copy(newSlideShow, slideShow)

		// fmt.Println("1.2 Repetition:", repetition)

		// Generate a number for number of mutation from the original slide show
		numberOfMutation = r.Intn(lenSlideShow / 2)

		// Ensure there's at least one mutation to be different from the first slide show
		if numberOfMutation == 0 {
			numberOfMutation++
		}

		// fmt.Println("1.3 Repetition:", repetition)

		// Randomly select 2 photo to swap for numberOfMutation iteration
		for j := 0; j < numberOfMutation; j++ {
			// Get 2 random photo in the slide show to swap
			firstPhotoPosition = r.Intn(lenSlideShow)
			secondPhotoPosition = r.Intn(lenSlideShow)

			// fmt.Println("1.4 Repetition:", repetition)

			// Ensure the 2 positions are unique
			EnsureUniqueNumber(&firstPhotoPosition, &secondPhotoPosition, lenSlideShow)

			// fmt.Println("1.5 Repetition:", repetition)

			initialScore := 0
			if len(newSlideShow) > firstPhotoPosition+1 {
				initialScore += CalcScoreBetweenTwo(newSlideShow[firstPhotoPosition], newSlideShow[firstPhotoPosition+1])
			}

			if 0 <= firstPhotoPosition-1 {
				initialScore += CalcScoreBetweenTwo(newSlideShow[firstPhotoPosition], newSlideShow[firstPhotoPosition-1])
			}

			if len(newSlideShow) > secondPhotoPosition+1 {
				initialScore += CalcScoreBetweenTwo(newSlideShow[secondPhotoPosition], newSlideShow[secondPhotoPosition+1])
			}

			if 0 <= secondPhotoPosition-1 {
				initialScore += CalcScoreBetweenTwo(newSlideShow[secondPhotoPosition], newSlideShow[secondPhotoPosition-1])
			}

			newScore := 0
			if len(newSlideShow) > secondPhotoPosition+1 {
				newScore += CalcScoreBetweenTwo(newSlideShow[firstPhotoPosition], newSlideShow[secondPhotoPosition+1])
			}

			if 0 <= secondPhotoPosition-1 {
				newScore += CalcScoreBetweenTwo(newSlideShow[firstPhotoPosition], newSlideShow[secondPhotoPosition-1])
			}

			if len(newSlideShow) > firstPhotoPosition+1 {
				newScore += CalcScoreBetweenTwo(newSlideShow[secondPhotoPosition], newSlideShow[firstPhotoPosition+1])
			}

			if 0 <= firstPhotoPosition-1 {
				newScore += CalcScoreBetweenTwo(newSlideShow[secondPhotoPosition], newSlideShow[firstPhotoPosition-1])
			}

			if initialScore > newScore {
				// j--
				continue
			}

			// Swap the photo
			swap = newSlideShow[firstPhotoPosition]
			newSlideShow[firstPhotoPosition] = newSlideShow[secondPhotoPosition]
			newSlideShow[secondPhotoPosition] = swap
		}

		// fmt.Println("1.6 Repetition:", repetition)

		// Store the new slide show into the population
		set = append(set, newSlideShow)
	}

	// 2. Calculate and pick the fittest slide show
	// fmt.Println("2.0 Repetition:", repetition)

	// Store the fittest genetic
	fittestSlideShow := 0
	secondFittestSlideShow := 0
	highestScore := 0
	secondHighestScore := 0

	// Traverse to all slide shows of population and get the fittest slide show in set
	for k := range set {
		if CalcScore(set[k]) > highestScore {
			secondFittestSlideShow = fittestSlideShow
			secondHighestScore = highestScore

			fittestSlideShow = k
			highestScore = CalcScore(set[k])
		} else if CalcScore(set[k]) > secondHighestScore {
			secondFittestSlideShow = k
			highestScore = CalcScore(set[k])
		}
	}

	// 3. Create an offspring from the fittest slide show and a random slide show
	// fmt.Println("3.0 Repetition:", repetition)

	// The random slide show selected could be the fittest slide show as well,
	// which will cause the new offspring to have
	// the same gene as the fittest slide show prior to mutation

	// Get the second best parent
	randomParent := set[secondFittestSlideShow]

	// Mate the two parents:
	// Select a random point and length in the first parent
	// Put the genes into the new offspring
	// Traverse through the second parent starting
	// at the end of  position of the selected gene of first parent
	// Insert the gene into the offspring if the gene does not exist in the offspring

	// Create the new offspring
	offspring := make([]Photo, 0)

	// Create a map to store the photos id
	offSpringIDList := make(map[int]struct{})

	// Select start and length of gene from the first parent
	startPositionFirstParent := r.Intn(lenSlideShow)
	lengthGeneOfFirstParent := r.Intn(lenSlideShow - startPositionFirstParent)
	endPositionFirstParent := startPositionFirstParent + lengthGeneOfFirstParent

	// Insert the selected first parent gene into the offspring
	for _, p := range set[fittestSlideShow][startPositionFirstParent:endPositionFirstParent] {
		offspring = append(offspring, p)
		offSpringIDList[p.id] = struct{}{}
	}

	// Iterate second parent from end gene position of first parent till end
	for _, p := range randomParent[endPositionFirstParent:] {
		// Go to next gene if current gene already exist in offspring
		if _, ok := offSpringIDList[p.id]; ok {
			continue
		}

		// Add gene to offspring if this iteration is not skipped
		offspring = append(offspring, p)
		offSpringIDList[p.id] = struct{}{}
	}

	// Iterate second parent from start to end gene of first parent
	for _, p := range randomParent[:endPositionFirstParent] {
		// Go to next gene if current gene already exist in offspring
		if _, ok := offSpringIDList[p.id]; ok {
			continue
		}

		offspring = append(offspring, p)
		offSpringIDList[p.id] = struct{}{}
	}

	// 4. Mutate the offspring
	// fmt.Println("4.0 Repetition:", repetition)

	numberOfMutation = r.Intn(lenSlideShow / 2)
	// numberOfMutation = r.Intn(7)
	// if numberOfMutation > 0 {
	// 	numberOfMutation = 1
	// }
	for i := 0; i < numberOfMutation; i++ {
		// Get 2 random photo in the slide show to swap
		firstPhotoPosition = r.Intn(lenSlideShow)
		secondPhotoPosition = r.Intn(lenSlideShow)

		// Ensure the 2 positions are unique
		EnsureUniqueNumber(&firstPhotoPosition, &secondPhotoPosition, lenSlideShow)

		swap = offspring[firstPhotoPosition]
		offspring[firstPhotoPosition] = offspring[secondPhotoPosition]
		offspring[secondPhotoPosition] = swap
	}

	// 5. Repeat by adding mutated offspring to the population set
	// fmt.Println("5.0 Repetition:", repetition)

	// Send max score to the UI and record highest slide show
	highestSlideShow := slideShow

	if CalcScore(slideShow) > maxScore {
		maxScore = CalcScore(slideShow)
		highestSlideShow = slideShow
	}

	if CalcScore(offspring) > maxScore {
		maxScore = CalcScore(offspring)
		highestSlideShow = offspring
	}

	if CalcScore(set[fittestSlideShow]) > maxScore {
		maxScore = CalcScore(set[fittestSlideShow])
		highestSlideShow = set[fittestSlideShow]
	}

	broadcast <- Message{
		Action: actionMaxScore,
		Data:   maxScore,
	}

	broadcast <- Message{
		Action: actionData,
		Data: Result{
			X: time.Now().Format("15:04:05"),
			Y: maxScore,
		},
	}

	// If there are any more repetition, call it again
	if repetition != 0 {
		repetition--

		// Recursive
		// TODO pass offspring to next generation even if parent has better score to avoid local optimum
		highestSlideShow = GeneticAlgorithm(highestSlideShow, r, repetition)
	}

	return highestSlideShow
}

// CalcScore is the fitness score calculator
func CalcScore(slideShow []Photo) (score int) {
	if len(slideShow) <= 1 {
		return 0
	}

	for k, p := range slideShow[1:] {
		currentScore := CalcScoreBetweenTwo(p, slideShow[k])

		score += currentScore
	}

	return
}
