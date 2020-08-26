package simple_logger

func findPadding(arr []Padding, elem Padding) int {
	if arr == nil {
		return -1
	}
	for index, element := range arr {
		if element == elem {
			return index
		}
	}

	return -1
}

func findStrings(arr []string, elem string) int {
	if arr == nil {
		return -1
	}
	for index, element := range arr {
		if element == elem {
			return index
		}
	}

	return -1
}

func findInts(arr []int, elem int) int {
	if arr == nil {
		return -1
	}
	for index, element := range arr {
		if element == elem {
			return index
		}
	}

	return -1
}
