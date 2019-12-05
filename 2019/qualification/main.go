package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Photo
type Photo struct {
	orientation      byte // H or V
	nrOfTag          int
	tags             map[string]struct{}
	isUsedAsVertical bool
}

func main() {
	// Read file
	photos, nrOfPhotos := ReadFile()
	fmt.Println("Number of photos:", nrOfPhotos)

	// Assign vertical
	answer := AssignVertical(photos)

	// Genetic algorithm
	fmt.Println(time.Now().Local())
	answer = GeneticAlgorithm(answer, rand.New(rand.NewSource(time.Now().Unix())), 10000)
	fmt.Println(time.Now().Local())

	// Debug
	// fmt.Println("DEBUG:", answer)

	// Final score
	fmt.Println("Final score:")
	fmt.Println(CalcScore(answer))
}

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
		if !v.isUsedAsVertical {
			// If no vertical image left, discard this image from use
			if k+1 == len(singleVertical) {
				break
			}

			// Pick another vertical photo to form a slide
			for j, v1 := range singleVertical[k+1:] {
				if !v1.isUsedAsVertical && CalcScoreBetweenTwo(v, v1) > 0 {
					// Append v1 into v
					AppendVerticalPhoto(&singleVertical[k], &singleVertical[k+1+j])

					singleVertical[k].isUsedAsVertical = true
					singleVertical[k+1+j].isUsedAsVertical = true

					answer = append(answer, singleVertical[k])

					break
				}
			}
		}
	}

	// Process all remaining images (Like those left because their score with other vertical images are all zero)
	for k, v := range singleVertical {
		if !v.isUsedAsVertical && k+1 < len(singleVertical) {
			for j, v1 := range singleVertical[k+1:] {
				if !v1.isUsedAsVertical {
					// Append next unassigned vertical photo to v
					AppendVerticalPhoto(&singleVertical[k], &singleVertical[k+1+j])

					singleVertical[k].isUsedAsVertical = true
					singleVertical[k+1+j].isUsedAsVertical = true

					break
				}
			}
		}
	}

	return
}

func AppendVerticalPhoto(photo, photo1 *Photo) {
	for t, _ := range photo1.tags {
		if _, ok := photo.tags[t]; !ok {
			photo.nrOfTag++
			photo.tags[t] = struct{}{}
		}
	}
}

func ReadFile() (photos []Photo, nrOfPhotos int) {
	// Define file location
	// filePath := "qualification_round_2019/a_example.txt"
	// filePath := "qualification_round_2019/b_lovely_landscapes.txt"
	filePath := "qualification_round_2019/c_memorable_moments.txt"
	// filePath := "qualification_round_2019/d_pet_pictures.txt"
	fmt.Println("Importing ......")
	fmt.Println("File used:", filePath)
	fmt.Println()

	// Read line
	file, err := os.Open(filePath)

	if err != nil {
		panic("Cant open file for reading")
	}
	defer file.Close()
	ioReader := bufio.NewReader(file)

	// Read first line: number of photos
	line, _ := ioReader.ReadString('\n')
	nrOfPhotos, _ = strconv.Atoi(strings.TrimSpace(line))

	// Read remaining lines: other photos
	for {
		line, err := ioReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err.Error())
		}

		// Process the line
		lines := strings.Split(strings.TrimSpace(line), " ")
		nrOfTag, _ := strconv.Atoi(lines[1])

		// Create the instance
		photo := Photo{
			orientation: byte(lines[0][0]),
			nrOfTag:     nrOfTag,
			tags:        map[string]struct{}{},
		}

		// Assign tags to photo
		for _, t := range lines[2:] {
			photo.tags[t] = struct{}{}
		}

		photos = append(photos, photo)
	}

	return
}

func GeneticAlgorithm(answer []Photo, r *rand.Rand, repetition int) []Photo {
	var set [][]Photo

	size := 100
	maxLen := len(answer)
	x := make(map[string]struct{})

	for i := 0; i < size; i++ {
		newAns := make([]Photo, len(answer))
		copy(newAns, answer)

		// Create a map that store position and randomly select them
		allZero := make(map[int]struct{})
		allZeroSlice := make([]int, 1)

		for k, _ := range newAns[1:] {
			leftCount := len(newAns[k+1-1].tags)
			rightCount := len(newAns[k+1].tags)
			overlapCount := 0

			for _, t1 := range newAns[k+1].tags {
				overlap := false
				for _, t0 := range answer[k+1-1].tags {
					if t0 == t1 {
						overlap = true
						break
					}
				}
				if overlap {
					leftCount--
					rightCount--
					overlapCount++
				}
			}
			score := min(leftCount, rightCount, overlapCount)
			if score == 0 {
				allZero[k+1] = struct{}{}
				allZeroSlice = append(allZeroSlice, k+1)
			}
		}

		rand1 := allZeroSlice[r.Intn(len(allZeroSlice))]
		rand2 := r.Intn(maxLen)

		for _, ok := allZero[rand2]; ok; _, ok = allZero[rand2] {
			rand2 = r.Intn(maxLen)
		}

		if rand1 == rand2 {
			if rand2+1 != maxLen {
				rand2++
			} else {
				rand2--
			}
		}

		if _, ok := x[string(rand1)+string(rand2)]; ok {
			i--
			continue
		}
		x[string(rand1)+string(rand2)] = struct{}{}
		x[string(rand2)+string(rand1)] = struct{}{}

		tempPhoto := newAns[rand1]

		newAns[rand1] = newAns[rand2]
		newAns[rand2] = tempPhoto
		set = append(set, newAns)
	}

	// Get the best genetic
	maxPosition := 0
	maxScore := CalcScore(set[0])

	for k := range set[1:] {
		if CalcScore(set[k+1]) > maxScore {
			maxPosition = k + 1
		}
	}

	// Select a random instance and reproduce
	// randomInstance := r.Intn(len(set))
	// for randomInstance == maxPosition {
	// 	randomInstance = r.Intn(len(set))
	// }
	if repetition != 0 {
		// for i := 0; i < 1; i++ {
		// 	// Random mutation
		// 	rand3 := r.Intn(maxLen)
		// 	anotherRand := rand3
		//
		// 	if rand3+1 == maxLen {
		// 		anotherRand--
		// 	} else {
		// 		anotherRand++
		// 	}
		// 	temp := set[maxPosition][rand3]
		// 	set[maxPosition][rand3] = set[maxPosition][anotherRand]
		// 	set[maxPosition][anotherRand] = temp
		// }

		repetition--

		// Recursive
		answer = GeneticAlgorithm(set[maxPosition], r, repetition)
	}

	if CalcScore(set[maxPosition]) > CalcScore(answer) {
		answer = set[maxPosition]
	}

	return answer
}

func CalcScore(answer []Photo) int {
	if len(answer) <= 1 {
		return 0
	}

	score := 0
	var scoreArr []int

	for k, p := range answer[1:] {
		currentScore := CalcScoreBetweenTwo(p, answer[k])

		score += currentScore
		scoreArr = append(scoreArr, currentScore)
	}

	// fmt.Println(scoreArr)

	return score
}

func CalcScoreBetweenTwo(p1, p2 Photo) int {
	p1TagCount := p1.nrOfTag
	p2TagCount := p2.nrOfTag
	overlap := 0

	for t, _ := range p1.tags {
		if _, ok := p2.tags[t]; ok {
			overlap++
			p1TagCount--
			p2TagCount--
		}
	}

	return min(p1TagCount, p2TagCount, overlap)
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
