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
	orientation byte // H or V
	nrOfTag     int
	tags        []string
	takenV      bool
}

func main() {
	var photos []Photo
	var nrOfPhotos int
	
	// Read line
	// file, err := os.Open("qualification_round_2019/a_example.txt")
	// file, err := os.Open("qualification_round_2019/b_lovely_landscapes.txt")
	file, err := os.Open("qualification_round_2019/c_memorable_moments.txt")
	// file, err := os.Open("qualification_round_2019/d_pet_pictures.txt")
	if err != nil {
		panic("Cant open file for reading")
	}
	defer file.Close()
	x := bufio.NewReader(file)
	// Read number of photos
	line, _ := x.ReadString('\n')
	nrOfPhotos, _ = strconv.Atoi(strings.TrimSpace(line))
	
	// Read photos
	for {
		line, err := x.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Eod of file reached!")
			break
		}
		
		lines := strings.Split(strings.TrimSpace(line), " ")
		nrOfTag, _ := strconv.Atoi(lines[1])
		
		photo := Photo{
			orientation: byte(lines[0][0]),
			nrOfTag:     nrOfTag,
			tags:        lines[2:],
		}
		photos = append(photos, photo)
	}
	
	fmt.Println("Number of photos:", nrOfPhotos)
	
	var answer []Photo
	for k, p := range photos {
		if p.orientation != 'V' {
			answer = append(answer, p)
		} else if !p.takenV {
			p.takenV = true
			for j, pp := range photos[k+1:] {
				if pp.orientation == 'V' && !pp.takenV {
					photos[k+1+j].takenV = true
					for _, pptags := range pp.tags {
						exist := false
						for _, ptags := range p.tags {
							if ptags == pptags {
								exist = true
								break
							}
						}
						if !exist {
							p.tags = append(p.tags, pptags)
						}
					}
				}
			}
			answer = append(answer, p)
		}
	}
	
	fmt.Println("Answer:")
	for _, p := range answer {
		if true {
			fmt.Printf("%c, %d, %s\n", p.orientation, p.nrOfTag, p.tags)
		}
	}
	
	// Algorithm
	fmt.Println("*****Algorithm start")
	answer = GeneticAlgo(answer, rand.New(rand.NewSource(time.Now().Unix())), 500)
	fmt.Println("******Algorithm end")
	
	fmt.Println("Score:")
	fmt.Println(CalcScore(answer))
}

func GeneticAlgo(answer []Photo, r *rand.Rand, repetition int) []Photo {
	var set [][]Photo
	
	size := 15
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
		fmt.Println(allZero)
		fmt.Println(allZeroSlice)
		fmt.Println(CalcScore(answer))
		panic("")
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
		
		// fmt.Println("rand1",rand1, "rand2", rand2)
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
		for i := 0; i < 1; i++ {
			// Random mutation
			rand3 := r.Intn(maxLen)
			anotherRand := rand3
			
			if rand3+1 == maxLen {
				anotherRand--
			} else {
				anotherRand++
			}
			temp := set[maxPosition][rand3]
			set[maxPosition][rand3] = set[maxPosition][anotherRand]
			set[maxPosition][anotherRand] = temp
		}
		
		repetition--
		
		// Recursive
		answer = GeneticAlgo(set[maxPosition], r, repetition)
	}
	
	if CalcScore(set[maxPosition]) > CalcScore(answer) {
		answer = set[maxPosition]
	}
	
	return answer
}

func CalcScore(answer []Photo) int {
	if len(answer) == 0 {
		return 0
	}
	
	score := 0
	var scoreArr []int
	
	for k, p := range answer[1:] {
		leftCount := len(answer[k+1-1].tags)
		rightCount := len(answer[k+1].tags)
		overlapCount := 0
		
		for _, t1 := range p.tags {
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
		
		score += min(leftCount, rightCount, overlapCount)
		scoreArr = append(scoreArr, min(leftCount, rightCount, overlapCount))
	}
	
	fmt.Println(scoreArr)
	
	return score
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
