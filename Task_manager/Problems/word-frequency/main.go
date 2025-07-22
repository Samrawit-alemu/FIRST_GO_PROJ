package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func WordFrequencyCounter(text string) map[string]int {
	var cleanedText strings.Builder
	for _, char := range text {
		if unicode.IsLetter(char) || unicode.IsSpace(char) {
			cleanedText.WriteRune(unicode.ToLower(char))
		}
	}
	words := strings.Fields(cleanedText.String())
	frequencyMap := make(map[string]int)
	for _, word := range words {
		frequencyMap[word]++
	}
	return frequencyMap
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a sentence to analyze: ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	wordFreq := WordFrequencyCounter(userInput)

	fmt.Println("\n--- Word Frequencies ---")
	if len(wordFreq) == 0 {
		fmt.Println("(No valid words found in your input)")
	} else {
		for word, count := range wordFreq {
			fmt.Printf("  - %s: %d\n", word, count)
		}
	}
}
