package main

import (
	"fmt"
)

type tag struct {
	photo    *Photo
	previous *tag
	next     *tag
}

type tagWrap struct {
	start *tag
	end   *tag
}

func startTagAlgorithm(filePath string) {
	// Read file
	fmt.Println("Importing ......")
	photos, _, _ := ReadFile(filePath)
	// fmt.Println("Number of photos:", nrOfPhotos)

	// Assign vertical
	// fmt.Println("Assigning vertical photos ......")
	photos = AssignVertical(photos)
	// photos = assignEasyVertical(photos)

	// Photos Length
	// fmt.Println("Slide show length:", len(photos))

	// Update current score
	updateAllCurrentScore(photos)
	score := CalcScore(photos)

	// Algorithm
	fmt.Println(filePath, "- Running tag algorithm ......")
	fmt.Println(filePath, "- Initial score:", score)
	fmt.Println(filePath, "- Initial Length:", len(photos))

	photos = tagAlgorithm(photos)

	// Final score
	score = CalcScore(photos)
	fmt.Println(filePath, "-- Final score:", score)
	fmt.Println(filePath, "-- Final Length:", len(photos))

	wg.Done()
}

func tagAlgorithm(photos []Photo) []Photo {
	// Create map of tags and node list
	tagMap := make(map[string]*tagWrap)

	for k := range photos {
		for j := range photos[k].tags {
			if _, ok := tagMap[j]; !ok {
				tagMap[j] = initializeNewTag(&photos[k])
			} else {
				// If tag exist before
				appendNewTag(&photos[k], tagMap[j])
			}
		}
	}
	// END

	// Get photos that has the tag
	samePhotosMap := make(map[int][]*Photo)

	for k := range photos {
		var samePhotos []*Photo
		for j := range photos[k].tags {
			samePhotos = append(samePhotos, getPhotosInTag(&photos[k], tagMap[j])...)
		}
		samePhotosMap[photos[k].id] = samePhotos
	}
	// END

	// Assign to slide show
	assigned := make(map[int]struct{})

	// Store answer
	var answer []Photo
	currentPhoto := &photos[0]

	// Store unassigned
	var storage []*Photo

	for i := 0; i < len(photos); i++ {
		if _, ok := assigned[photos[i].id]; !ok {
			currentPhoto = &photos[i]
			answer = append(answer, photos[i])
			assigned[photos[i].id] = struct{}{}

			solve(currentPhoto, &assigned, samePhotosMap, &answer, &storage)
		}
	}

	// fmt.Println("Storage length:", len(storage))
	return answer
}

func solve(currentPhoto *Photo, assigned *map[int]struct{}, samePhotosMap map[int][]*Photo, answer *[]Photo, storage *[]*Photo) {
	// Get max score
	var maxPhoto *Photo
	var maxScore int
	for j := range samePhotosMap[currentPhoto.id] {
		if _, ok := (*assigned)[samePhotosMap[currentPhoto.id][j].id]; !ok {
			newScore := CalcScoreBetweenTwo(*currentPhoto, *samePhotosMap[currentPhoto.id][j])
			if newScore > maxScore {
				maxScore = newScore
				maxPhoto = samePhotosMap[currentPhoto.id][j]
			}

			if newScore >= currentPhoto.nrOfTag/2 {
				break
			}
		}
	}

	if maxPhoto == nil {
		*storage = append(*storage, currentPhoto)
	}

	// Assign to the photo with max score
	if maxPhoto != nil {
		(*assigned)[maxPhoto.id] = struct{}{}
		*answer = append(*answer, *maxPhoto)
		// printSlideShow(answer)

		// Start on the assigned photo
		currentPhoto = maxPhoto
		solve(currentPhoto, assigned, samePhotosMap, answer, storage)
	}
}

func initializeNewTag(p *Photo) *tagWrap {
	startTag := &tag{
		photo:    nil,
		previous: nil,
		next:     nil,
	}

	endTag := &tag{
		photo:    nil,
		previous: nil,
		next:     nil,
	}

	tag := &tag{
		photo:    p,
		previous: startTag,
		next:     endTag,
	}

	startTag.next = tag
	endTag.previous = tag

	tagWrap := &tagWrap{
		start: startTag,
		end:   endTag,
	}

	return tagWrap
}

func appendNewTag(p *Photo, tw *tagWrap) {
	tag := &tag{
		photo:    p,
		previous: tw.end.previous,
		next:     tw.end,
	}

	tw.end.previous.next = tag
	tw.end.previous = tag
}

func getPhotosInTag(photo *Photo, tw *tagWrap) (photos []*Photo) {
	photoTag := tw.start.next
	for photoTag != tw.end {
		if photoTag.photo != photo {
			photos = append(photos, photoTag.photo)
		}

		photoTag = photoTag.next
	}

	return
}

func printTagMap(tagMap map[string]*tagWrap) {
	fmt.Println("Result:")
	for k := range tagMap {
		fmt.Println("Tag: ", k)
		printTag := tagMap[k].start.next
		for printTag != tagMap[k].end {
			fmt.Print(printTag.photo.id, " ")
			printTag = printTag.next
		}
		fmt.Println()
	}
}

func printSamePhotos(photos []*Photo) {
	fmt.Print("Same photo: ")
	for k := range photos {
		fmt.Print(photos[k].id, " ")
	}
	fmt.Println()
}

func slideShowNoDuplicate(photos []Photo) bool {
	duplicate := make(map[int]struct{})

	for k := range photos {
		if _, ok := duplicate[photos[k].id]; ok {
			return false
			// fmt.Println("Duplicate:", photos[k].id)
		} else {
			duplicate[photos[k].id] = struct{}{}
		}
	}

	return true
}
