package utils

func FindNOccurance(substr, entry string, occurance int) int {
	substrLen := len(substr)
	entryLen := len(entry)
	count := 0
	for i := 0; i < substrLen; i++ {
		if i+entryLen > substrLen {
			return -1
		} else if substr[i:i+entryLen] == entry {
			count += 1
		}
		if count == occurance {
			return i
		}
	}
	return -1
}
