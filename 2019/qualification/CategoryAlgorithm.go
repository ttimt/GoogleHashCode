package main

import (
	"fmt"
)

func StartCategoryAlgorithm() {
	// Read file
	fmt.Println("Importing ......")
	photos, nrOfPhotos := ReadFile()
	fmt.Println("Number of photos:", nrOfPhotos)

	// Assign vertical
	fmt.Println("Assigning vertical photos ......")
	photos = AssignVertical(photos)

	// Update photos length
	slideShowLength = len(photos)

	// Algorithm
	fmt.Println("Running algorithm ......")
	maxScore = CalcScore(photos)
	photos = CategoryAlgorithm(photos)

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

func CategoryAlgorithm(slideShow []Photo) []Photo {
	answer := make([]Photo, slideShowLength)
	answer = answer[0:1]

	photo := &slideShow[0]
	(*photo).used = true
	answer = append(answer, *photo)

	for i := 0; i < slideShowLength; i++ {
		done := false
		for j, pp := range slideShow {
			if !pp.used && CalcScoreBetweenTwo(*photo, pp) >= 3 {
				answer = append(answer, slideShow[j])
				slideShow[j].used = true
				*photo = slideShow[j]

				done = true
				fmt.Println("Iteration:", i)
				break
			}
		}

		if !done {
			for j, pp := range slideShow {
				if !pp.used {
					answer = append(answer, *photo)
					slideShow[j].used = true
					*photo = slideShow[j]

					fmt.Println("Iteration:", i)
					break
				}
			}
		}
	}

	return answer
}
