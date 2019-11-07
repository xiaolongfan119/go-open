package util

func InSlice(slice []interface{}, target interface{}) bool {

	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func InSliceInt(slice []int, target int) bool {

	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func InSliceString(slice []string, target string) bool {

	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
