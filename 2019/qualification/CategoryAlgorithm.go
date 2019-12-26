package main

import (
	"fmt"
)

func startCategoryAlgorithm(filePath string) {
	// Read file
	fmt.Println("Importing ......")
	photos, _, maxNrOfTags := ReadFile(filePath)
	// fmt.Println("Number of photos:", nrOfPhotos)

	// Assign vertical
	// fmt.Println("Assigning vertical photos ......")
	// photos = AssignVertical(photos)
	photos = assignEasyVertical(&photos)
	//
	// Photos Length
	// fmt.Println("Slide show length:", len(photos))

	// Update current score
	updateAllCurrentScore(photos)
	score := CalcScore(photos)

	// Algorithm
	fmt.Println(filePath, "- Running category algorithm ......")
	// fmt.Println("Initial score:", score)
	var ok bool

	ite := maxNrOfTags / 2
	// fmt.Println("Max iteration to make:", ite)
	for i := 0; i < ite; i++ {
		photos, ok = CategoryAlgorithm(photos, i, filePath)

		if !ok {
			break
		}
	}

	// Final score
	score = CalcScore(photos)
	fmt.Println(filePath, "- Final score:", score)
	wg.Done()
}

// CategoryAlgorithm greedy
func CategoryAlgorithm(photos []Photo, maxNrOfTags int, filePath string) ([]Photo, bool) {
	maxScoreNew := CalcScore(photos)
	fmt.Println(filePath, "- New score:", maxScoreNew)
	if maxScore == maxScoreNew {
		return photos, false
	} else {
		maxScore = maxScoreNew
	}

	for k := range photos {
		// fmt.Println("1", k)
		currentTotal := 0

		for j := range photos {
			// fmt.Println("2", k, j)
			currentTotal = photos[k].currentScore

			if k != j {
				currentTotal += photos[j].currentScore

				// fmt.Println("3", " ")
				// Get swap score
				newTotal := 0
				if k-1 >= 0 {
					newTotal += CalcScoreBetweenTwo(photos[k-1], photos[j])
				}
				if k+1 < len(photos) {
					newTotal += CalcScoreBetweenTwo(photos[j], photos[k+1])
				}
				if j-1 >= 0 {
					newTotal += CalcScoreBetweenTwo(photos[j-1], photos[k])
				}
				if j+1 < len(photos) {
					newTotal += CalcScoreBetweenTwo(photos[k], photos[j+1])
				}

				// fmt.Println("New", newTotal, "Current", currentTotal)
				if newTotal > currentTotal && newTotal != 0 {
					// fmt.Println("4", " ")
					// Swap
					temp := photos[k]
					photos[k] = photos[j]
					photos[j] = temp

					photos[k].currentScore = updateCurrentScore(photos, k)
					photos[j].currentScore = updateCurrentScore(photos, j)
					break
				}
				// fmt.Println("5", " ")
			}
		}
	}

	return photos, true
}

func updateAllCurrentScore(photos []Photo) {
	for k := range photos {
		score := updateCurrentScore(photos, k)
		photos[k].currentScore = score
	}
}

func updateCurrentScore(photos []Photo, pos int) int {
	score := 0

	if pos-1 >= 0 {
		score += CalcScoreBetweenTwo(photos[pos-1], photos[pos])
	}

	if pos+1 < len(photos) {
		score += CalcScoreBetweenTwo(photos[pos], photos[pos+1])
	}

	return score
}
