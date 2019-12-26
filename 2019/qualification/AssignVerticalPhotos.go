package main

// AssignVertical assign vertical photos
func AssignVertical(photos []Photo) (answer []Photo) {
	// Store unassigned vertical photos
	var singleVertical []Photo
	var currentOverlap int
	var smallestOverlapPhotoPosition int
	var smallestOverlap int

	// Filter out all horizontal images
	for _, p := range photos {
		if p.orientation != 'V' {
			answer = append(answer, p)
		} else {
			singleVertical = append(singleVertical, p)
		}
	}

	lenSingleVertical := len(singleVertical)

	// Return if no vertical image or only one
	if lenSingleVertical <= 1 {
		return answer
	}

	// Process all vertical images
	for k, v := range singleVertical {
		// Process current image
		// If no vertical image left, discard this image from use
		// Pick another vertical photo to form a slide
		if !v.isUsedAsVertical && k+1 < lenSingleVertical {
			// Find the one with least overlap
			smallestOverlapPhotoPosition = k + 1
			smallestOverlap = CalcNumberOfOverlapTags(v, singleVertical[k+1])

			for j, v1 := range singleVertical[k+1:] {
				if !v1.isUsedAsVertical {
					currentOverlap = CalcNumberOfOverlapTags(v, v1)
					if currentOverlap < smallestOverlap {
						smallestOverlapPhotoPosition = k + 1 + j
						smallestOverlap = currentOverlap

						// Improve performance by ending faster when best solution is found
						if smallestOverlap == 0 {
							break
						}
					}
				}
			}

			// Append smallestOverlapPhoto into v
			AppendVerticalPhoto(&singleVertical[k], &singleVertical[smallestOverlapPhotoPosition])
			answer = append(answer, singleVertical[k])
		}
	}

	return
}

func assignEasyVertical(photos *[]Photo) (answer []Photo) {
	// Store unassigned vertical photos
	var singleVertical []Photo

	for _, p := range *photos {
		if p.orientation != 'V' {
			answer = append(answer, p)
		} else {
			singleVertical = append(singleVertical, p)
		}
	}

	for i := 0; i < len(singleVertical); i += 2 {
		if i+1 < len(singleVertical) {
			AppendVerticalPhoto(&singleVertical[i], &singleVertical[i+1])
			answer = append(answer, singleVertical[i])
		}
	}

	return
}
