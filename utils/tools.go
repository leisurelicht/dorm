package utils

// DuplicateString 去重
func DuplicateString(s []string) (result []string) {
	temp := map[string]bool{}
	for _, item := range s {
		if _, ok := temp[item]; !ok {
			temp[item] = false
			result = append(result, item)
		}
	}
	return result
}
