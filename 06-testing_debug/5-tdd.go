package testingdebug

func stringrev(s string) string {
	runes := []rune(s)
	for i := len(runes)/2 - 1; i >= 0; i-- {
		runes[i], runes[len(runes)-1-i] = runes[len(runes)-1-i], runes[i]
	}

	return string(runes)
}
