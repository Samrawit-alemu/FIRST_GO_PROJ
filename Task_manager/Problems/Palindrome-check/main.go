package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func IsPalindrome(text string) bool {
	var sanitizedBuilder strings.Builder
	for _, char := range text {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			sanitizedBuilder.WriteRune(unicode.ToLower(char))
		}
	}
	sanitizedText := sanitizedBuilder.String()

	if len(sanitizedText) <= 1 {
		return true
	}

	runes := []rune(sanitizedText)
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
		if runes[left] != runes[right] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a word or phrase to check if it's a palindrome: ")

	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	isPal := IsPalindrome(userInput)

	fmt.Println("\n--- Result ---")
	if isPal {
		fmt.Printf("Yes, \"%s\" is a palindrome!\n", userInput)
	} else {
		fmt.Printf("No, \"%s\" is not a palindrome.\n", userInput)
	}
}
