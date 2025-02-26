package app

import (
	"errors"
	"math"
	"strings"
)

// EstimateTokens estimates the number of tokens in a text using various methods
// method can be "average", "words", "chars", "max", "min", defaults to "max"
// "average" is the average of words and chars
// "words" is the word count divided by 0.75
// "chars" is the char count divided by 4
// "max" is the max of word and char
// "min" is the min of word and char
func EstimateTokens(text string, method string) (int, error) {
	wordCount := len(strings.Fields(text))
	charCount := len(text)
	tokensCountWordEst := float64(wordCount) / 0.75
	tokensCountCharEst := float64(charCount) / 4.0
	var output float64

	switch method {
	case "average":
		output = (tokensCountWordEst + tokensCountCharEst) / 2
	case "words":
		output = tokensCountWordEst
	case "chars":
		output = tokensCountCharEst
	case "max", "":
		output = math.Max(tokensCountWordEst, tokensCountCharEst)
	case "min":
		output = math.Min(tokensCountWordEst, tokensCountCharEst)
	default:
		return 0, errors.New("Invalid method. Use 'average', 'words', 'chars', 'max', or 'min'.")
	}

	return int(output), nil
}
