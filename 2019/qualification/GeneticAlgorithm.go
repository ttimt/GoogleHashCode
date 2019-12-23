package main

import (
	"fmt"
)

// StartAlgorithm will execute the whole algorithm
func StartAlgorithm(filePath string) {
	// Read file
	fmt.Println("Importing ......")
	photos, nrOfPhotos, maxNrOfTags := ReadFile(filePath)
	fmt.Println("Number of photos:", nrOfPhotos)
	fmt.Println("Max number of tags:", maxNrOfTags)

	// Assign vertical
	fmt.Println("Assigning vertical photos ......")
	photos = AssignVertical(photos)

	// Update photos length
	slideShowLength := len(photos)

	// Genetic algorithm
	fmt.Println("Running algorithm ......")
	// maxScore = CalcScore(photos)
	photos = GeneticAlgorithm(photos, repetition, slideShowLength)

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

// GeneratePopulation generates random population
func GeneratePopulation(slideShow *[]Photo, slideShowLength int) *[][]Photo {
	// Use to store the population and store the original slide show in the set
	var set [][]Photo
	set = append(set, *slideShow)

	// Store the random number of mutation to occur per slide show in a set
	var numberOfMutation, firstPhotoPosition, secondPhotoPosition int

	// Store the temporary photo when swapping
	var swap Photo

	// TODO verify random swap should have at least 1 score in between slide shows
	//  or else attempt a few more times then only swap with photo of lowest number of tag
	for i := 0; i < populationSize; i++ {
		// Store the new instance of slide show
		newSlideShow := make([]Photo, slideShowLength)
		copy(newSlideShow, *slideShow)

		// Generate a number for number of mutation from the original slide show
		numberOfMutation = r.Intn(slideShowLength-1) + 1

		// Ensure there's at least one mutation to be different from the first slide show
		if numberOfMutation == 0 {
			numberOfMutation++
		}

		numberOfRetry := 0

		// Randomly select 2 photo to swap for numberOfMutation iteration
		for j := 0; j < numberOfMutation; j++ {
			// Get 2 random photo in the slide show to swap
			firstPhotoPosition = r.Intn(slideShowLength)
			secondPhotoPosition = r.Intn(slideShowLength)

			// fmt.Println("1.4 Repetition:", repetition)

			// Ensure the 2 positions are unique
			EnsureUniqueNumber(&firstPhotoPosition, &secondPhotoPosition, slideShowLength)

			// fmt.Println("1.5 Repetition:", repetition)

			initialScore := 5
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

			if initialScore >= newScore {
				j--
				numberOfRetry++

				if numberOfRetry > 5 {
					j++
					numberOfRetry = 0
				} else {
					continue
				}
			}
			// if newScore <= 0 {
			// 	j--
			// 	continue
			// }

			// Swap the photo
			swap = newSlideShow[firstPhotoPosition]
			newSlideShow[firstPhotoPosition] = newSlideShow[secondPhotoPosition]
			newSlideShow[secondPhotoPosition] = swap
		}

		// Store the new slide show into the population
		set = append(set, newSlideShow)
	}

	return &set
}

// SelectFittest get best slide show
func SelectFittest(set *[][]Photo) int {
	// Store the fittest genetic
	fittestSlideShowPosition := 0
	highestScore := 0
	currentKScore := 0

	// Traverse to all slide shows of population and get the fittest slide show in set
	for k := range *set {
		currentKScore = CalcScore((*set)[k])
		if currentKScore > highestScore {
			fittestSlideShowPosition = k
			highestScore = currentKScore
		}
	}

	return fittestSlideShowPosition
}

// CreateOffspring create offspring from 2 parents
func CreateOffspring(set *[][]Photo, fittestSlideShowPosition int, slideShowLength int) []Photo {
	// The random slide show selected could be the fittest slide show as well,
	// which will cause the new offspring to have
	// the same gene as the fittest slide show prior to mutation

	// Get the second best parent
	randomParent := (*set)[r.Intn(len(*set))]

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
	startPositionFirstParent := r.Intn(slideShowLength / 2)
	lengthGeneOfFirstParent := r.Intn(slideShowLength - startPositionFirstParent)
	endPositionFirstParent := startPositionFirstParent + lengthGeneOfFirstParent

	startPositionFirstParent = 0
	lengthGeneOfFirstParent = slideShowLength * 3 / 4
	endPositionFirstParent = startPositionFirstParent + lengthGeneOfFirstParent

	// Insert the selected first parent gene into the offspring
	for _, p := range (*set)[fittestSlideShowPosition][startPositionFirstParent:endPositionFirstParent] {
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

	return offspring
}

// OffspringMutation mutates the offspring
func OffspringMutation(offspring *[]Photo, slideShowLength int) {
	hasMutation := r.Float64() < mutationRate

	if hasMutation {
		// Get 2 random photo in the slide show to swap
		firstPhotoPosition := r.Intn(slideShowLength)
		secondPhotoPosition := r.Intn(slideShowLength)

		// Ensure the 2 positions are unique
		EnsureUniqueNumber(&firstPhotoPosition, &secondPhotoPosition, slideShowLength)

		swap := (*offspring)[firstPhotoPosition]
		(*offspring)[firstPhotoPosition] = (*offspring)[secondPhotoPosition]
		(*offspring)[secondPhotoPosition] = swap
	}
}

// GeneticAlgorithm is the optimization algorithm
func GeneticAlgorithm(slideShow []Photo, repetition int, slideShowLength int) []Photo {
	// 1. Generate population / a set of slide shows
	// 2. Pick the fittest
	// 3. Create an offspring from the fittest and a random slide show
	// 4. Mutation to the offspring
	// 5. Repeat by adding the mutated offspring to the population set in (1)

	for i := 0; i < repetition; i++ {
		// 1. Generate a population / a set of slide shows
		fmt.Println("1.0 Repetition:", i)
		set := GeneratePopulation(&slideShow, slideShowLength)

		// 2. Calculate and pick the fittest slide show
		fmt.Println("2.0 Repetition:", i)
		fittestSlideShowPosition := SelectFittest(set)

		// 3. Create an offspring from the fittest slide show and a random slide show
		fmt.Println("3.0 Repetition:", i)
		offspring := CreateOffspring(set, fittestSlideShowPosition, slideShowLength)

		// 4. Mutate the offspring
		fmt.Println("4.0 Repetition:", i)
		OffspringMutation(&offspring, slideShowLength)

		// 4.5 Set fittest slide show as offspring if the offspring has lower score
		// if CalcScore((*set)[fittestSlideShowPosition]) > CalcScore(offspring) {
		// 	offspring = (*set)[fittestSlideShowPosition]
		// }

		// 5. Repeat by adding mutated offspring to the population set
		fmt.Println("5.0 Repetition:", i)
		slideShow = offspring

		// Set current offspring score to the UI
		// broadcast <- Message{
		// 	Action: actionData,
		// 	Data: Result{
		// 		X: time.Now().Format("15:04:05"),
		// 		Y: CalcScore(offspring),
		// 	},
		// }

		// Update max score
		if CalcScore(offspring) > maxScore {
			maxScore = CalcScore(offspring)
		}

		// Send the max score to the UI
		// broadcast <- Message{
		// 	Action: actionMaxScore,
		// 	Data:   maxScore,
		// }
	}

	return slideShow
}
