package utils

func WordsGenerator(wordLength int) string {
	var word = ""
	for i := 0; i < wordLength; i++ {
		word += "a"
	}
	return word
}
