package mysql

// Is element exists in array?
func inUint16Array(arr []uint16, ele uint16) bool {
	for _, val := range arr {
		if val == ele {
			return true
		}
	}
	return false
}

// Check exists, append if not exists.
func appendUint16Array(arr *[]uint16, ele uint16) {
	if inUint16Array(*arr, ele) {
		return
	}
	*arr = append(*arr, ele)
}
