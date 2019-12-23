package main

// AssignVertical assign vertical photos
func AssignVertical(photos []Photo) (answer []Photo) {
	// Store unassigned vertical photos
	var singleVertical []Photo
	// fmt.Println("Assign vertical: 1")
	// Filter out all horizontal images
	for _, p := range photos {
		if p.orientation != 'V' {
			answer = append(answer, p)
		} else {
			singleVertical = append(singleVertical, p)
		}
	}

	// fmt.Println("Assign vertical: 2")
	// Return if no vertical image or only one
	if len(singleVertical) <= 1 {
		return answer
	}

	// fmt.Println("Assign vertical: 3")
	// Process all vertical images
	for k, v := range singleVertical {
		// Process current image
		// If no vertical image left, discard this image from use
		if !v.isUsedAsVertical && k+1 < len(singleVertical) {
			// Pick another vertical photo to form a slide
			smallestOverlapPhotoPosition := k + 1
			smallestOverlap := CalcNumberOfOverlapTags(v, singleVertical[k+1])

			// fmt.Println("Assign vertical: 4 Image:", k)
			// Find the one with least overlap
			// TODO O(n^2) too slow
			for j, v1 := range singleVertical[k+1:] {
				if !v1.isUsedAsVertical && CalcNumberOfOverlapTags(v, v1) < smallestOverlap {
					// fmt.Println("Assign vertical: 5 Image:", j)
					smallestOverlapPhotoPosition = k + 1 + j
					smallestOverlap = CalcNumberOfOverlapTags(v, v1)

					// Improve performance by ending faster when best solution is found
					if smallestOverlap == 0 {
						break
					}
				}
			}

			// Append smallestOverlapPhoto into v
			AppendVerticalPhoto(&singleVertical[k], &singleVertical[smallestOverlapPhotoPosition])
			answer = append(answer, singleVertical[k])
		}
	}

	// fmt.Println("Assign vertical: 6")
	return
}

func assignEasyVertical(photos []Photo) (answer []Photo) {
	// Store unassigned vertical photos
	var singleVertical []Photo

	for _, p := range photos {
		if p.orientation != 'V' {
			answer = append(answer, p)
		} else {
			singleVertical = append(singleVertical, p)
		}
	}

	for i := 0; i < len(singleVertical); i += 2 {
		if i+1 < len(singleVertical) {
			AppendVerticalPhoto(&photos[i], &photos[i+1])
			answer = append(answer, photos[i])
		}
	}

	return
}
