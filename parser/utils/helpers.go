package utils

import "unicode"

func CountSpaces(line string) int {
	count := 0
	for i := 0; i < len(line); i++ {
		if unicode.IsSpace(rune(line[i])) {
			count++
		}else {
			break
		}
	}
	return count
}

func CheckForCharacter(line string) bool {
	for i := 0; i < len(line); i++ {
		if !unicode.IsSpace(rune(line[i])) {
			return true
		}
	}
	return false
}