package util

func InSlice(slice []interface{}, target interface{}) bool {

	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
