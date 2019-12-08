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

const (
	populationSize = 15
	repetition     = 2000
	// filePath       = "qualification_round_2019/a_example.txt"
	// filePath       = "qualification_round_2019/b_lovely_landscapes.txt"
	filePath = "qualification_round_2019/c_memorable_moments.txt"
	// filePath       = "qualification_round_2019/d_pet_pictures.txt"
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
	fmt.Println("Importing ......")
	photos, nrOfPhotos := ReadFile()
	fmt.Println("Number of photos:", nrOfPhotos)

	// Assign vertical
	fmt.Println("Assigning vertical photos ......")
	answer := AssignVertical(photos)

	// Genetic algorithm
	fmt.Println("Running algorithm ......")
	answer = GeneticAlgorithm(answer, rand.New(rand.NewSource(time.Now().Unix())), repetition)

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
	fmt.Println("File used:", filePath)

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

func GeneticAlgorithm(slideShow []Photo, r *rand.Rand, repetition int) []Photo {
	// 1. Generate population / a set of slide shows
	// 2. Pick the fittest
	// 3. Create an offspring from the fittest and a random slide show
	// 4. Mutation to the offspring
	// 5. Repeat by adding the mutated offspring to the population set in (1)

	// 1. Generate a population / a set of slide shows
	// Use to store the population and store the original slide show in the set
	var set [][]Photo
	set = append(set, slideShow)

	// Store the random number of mutation to occur per slide show in a set
	var numberOfMutation, firstPhotoPosition, secondPhotoPosition int

	// Length of slide show
	lenSlideShow := len(slideShow)

	// Store the temporary photo when swapping
	var swap Photo

	for i := 0; i < populationSize; i++ {
		// Store the new instance of slide show
		newSlideShow := make([]Photo, lenSlideShow)
		copy(newSlideShow, slideShow)

		// Generate a number for number of mutation from the original slide show
		numberOfMutation = r.Intn(lenSlideShow / 2)

		// Ensure there's at least one mutation to be different from the first slide show
		if numberOfMutation == 0 {
			numberOfMutation++
		}

		// Randomly select 2 photo to swap for numberOfMutation iteration
		for j := 0; j < numberOfMutation; j++ {
			// Get 2 random photo in the slide show to swap
			firstPhotoPosition = r.Intn(lenSlideShow)
			secondPhotoPosition = r.Intn(lenSlideShow)

			// Ensure the 2 positions are unique
			EnsureUniqueNumber(&firstPhotoPosition, &secondPhotoPosition, lenSlideShow)

			// Swap the photo
			swap = newSlideShow[firstPhotoPosition]
			newSlideShow[firstPhotoPosition] = newSlideShow[secondPhotoPosition]
			newSlideShow[secondPhotoPosition] = swap
		}

		// Store the new slide show into the population
		set = append(set, newSlideShow)
	}

	// 2. Calculate and pick the fittest slide show
	// Store the fittest genetic
	fittestSlideShow := 0
	highestScore := 0

	// Traverse to all slide shows of population and get the fittest slide show in set
	for k := range set {
		if CalcScore(set[k]) > highestScore {
			fittestSlideShow = k
		}
	}

	// 3. Create an offspring from the fittest slide show and a random slide show
	// The random slide show selected could be the fittest slide show as well,
	// which will cause the new offspring to have
	// the same gene as the fittest slide show prior to mutation

	// Get the random parent
	randomParent := set[r.Intn(len(set))]

	// Mate the two parents:
	// Select a random point and length in the first parent
	// Put the genes into the new offspring
	// Traverse through the second parent starting
	// at the end of  position of the selected gene of first parent
	// Insert the gene into the offspring if the gene does not exist in the offspring

	// Create the new offspring
	offspring := make([]Photo, lenSlideShow)

	// Select start and length of gene from the first parent
	startPositionFirstParent := r.Intn(lenSlideShow)
	lengthGeneOfFirstParent := r.Intn(lenSlideShow - startPositionFirstParent)
	endPositionFirstParent := startPositionFirstParent + lengthGeneOfFirstParent

	// Insert the selected first parent gene into the offspring
	for _, p := range set[fittestSlideShow][startPositionFirstParent:endPositionFirstParent] {
		offspring = append(offspring, p)
	}

	// Iterate second parent from end gene position of first parent till end
	for _, p := range randomParent[endPositionFirstParent:] {
		// Go to next gene if current gene already exist in offspring
		for _, pOff := range offspring {
			if IsPhotoEqual(&p, &pOff) {
				continue
			}
		}

		// Add gene to offspring if this iteration is not skipped
		offspring = append(offspring, p)
	}

	// Iterate second parent from start to start gene of first parent
	for _, p := range randomParent[:startPositionFirstParent] {
		// Go to next gene if current gene already exist in offspring
		for _, pOff := range offspring {
			if IsPhotoEqual(&p, &pOff) {
				continue
			}
		}

		offspring = append(offspring, p)
	}

	// 4. Mutate the offspring
	numberOfMutation = r.Intn(lenSlideShow / 2)
	numberOfMutation = 0
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
	if repetition != 0 {
		repetition--

		// Recursive
		slideShow = GeneticAlgorithm(set[fittestSlideShow], r, repetition)
	}

	if CalcScore(set[fittestSlideShow]) > CalcScore(slideShow) {
		slideShow = set[fittestSlideShow]
	}

	return slideShow
}

func IsPhotoEqual(photo1, photo2 *Photo) bool {
	// Compare if two Photo struct are equal

	// Compare orientation
	if photo1.orientation != photo2.orientation {
		return false
	}

	// Compare number of tags
	if photo1.nrOfTag != photo2.nrOfTag {
		return false
	}

	// Compare length of maps
	if len(photo1.tags) != len(photo2.tags) {
		return false
	}

	// Compare tags
	for k := range photo1.tags {
		if _, ok := photo2.tags[k]; !ok {
			return false
		}
	}

	return true
}

func EnsureUniqueNumber(x, y *int, len int) {
	if x == y {
		if *y+1 >= len {
			*y--
		} else {
			*y++
		}
	}
}

func CalcScore(slideShow []Photo) int {
	if len(slideShow) <= 1 {
		return 0
	}

	score := 0
	// var scoreArr []int

	for k, p := range slideShow[1:] {
		currentScore := CalcScoreBetweenTwo(p, slideShow[k])

		score += currentScore
		// scoreArr = append(scoreArr, currentScore)
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
