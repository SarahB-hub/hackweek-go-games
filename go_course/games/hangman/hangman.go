package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/tjarratt/babble"
)

func main() {

	word := getRandomWord()
	guessedLetters := initialiseGuessedLetters(word)
	hangmanState := 0
	incorrectGuesses := []string{}
	guessesRemaining := 9

	greeting := `
Welcome to Hangman!

You have 9 attempts to guess the correct word, and we've started you off with the first and last letter filled in.

Good luck!`

	fmt.Println(greeting)

	for {
		printGameState(word, guessedLetters, hangmanState, guessesRemaining) //passes as args to printResults func
		inputLetter := readInput()
		if len(inputLetter) != 1 {
			fmt.Println("Please guess a single letter.")
			continue
		}

		if checkIfLetterInWord(word, inputLetter) {
			guessedLetters[rune(inputLetter[0])] = true //convert to for the map rune requires specifying which index in the string
			printSpace(2)
			fmt.Printf("Correct guess! The letter %s is in the word.", inputLetter)
			printSpace(4)
		} else {
			printSpace(2)
			fmt.Printf("Incorrect guess! The letter %s is not in the word.", inputLetter)
			printSpace(2)
			incorrectGuesses = append(incorrectGuesses, inputLetter)
			hangmanState++
			guessesRemaining--
			printSpace(4)

		}

		if isWordGuessed(word, guessedLetters) {
			fmt.Println("----------------------------------------------------")
			fmt.Printf("\nWell done, the word was %v, you won!\n", word)
			fmt.Println("----------------------------------------------------")
			printSpace(2)
			break
		} else if hangmanState >= 9 {
			fmt.Println("----------------------------------------------------")
			fmt.Printf("Out of tries: the word was %v. Game over!\n", word)
			fmt.Println("----------------------------------------------------")
			printSpace(2)
			break
		}

		if len(incorrectGuesses) != 0 {
			fmt.Printf("Letters not included: ")
			fmt.Print(incorrectGuesses)
			printSpace(2)
		}

	}
}

func printSpace(amount int) {
	fmt.Print(strings.Repeat("\n", amount))
}

func getRandomWord() string {
	babbler := babble.NewBabbler()
	babbler.Count = 1
	randomWord := babbler.Babble()
	for len(randomWord) > 8 || len(randomWord) < 5 {
		randomWord = babbler.Babble()
	}
	return randomWord
}

func initialiseGuessedLetters(word string) map[rune]bool {
	guessedLetters := map[rune]bool{} //map = dict, with keytype rune ( for iteration later), value type bool
	guessedLetters[unicode.ToLower(rune(word[0]))] = true
	guessedLetters[unicode.ToLower(rune(word[len(word)-1]))] = true //sets the first and last letter as found in the map for initialisation

	return guessedLetters
}

func isWordGuessed(word string, guessedLetters map[rune]bool) bool {

	for _, character := range strings.ToLower(word) { //or the first letter will never match
		if !guessedLetters[character] {
			return false
		} //for each index named character in the range of word
		// if the character isn't in the guessedLetters map, return false
	}

	return true //once all the chars in the word exist in the guessed letters, they won!
}

func readInput() string {
	var inputLetter string
	fmt.Scanln(&inputLetter)

	return inputLetter
}

func checkIfLetterInWord(word string, inputLetter string) bool {
	return strings.Contains(word, inputLetter)
}

func printGameState(
	word string, guessedLetters map[rune]bool, hangmanState int, guessesRemaining int) {
	fmt.Println(getHangmanDrawing(hangmanState))
	printSpace(2)
	if guessesRemaining > 1 {
		fmt.Printf("You have %d guesses remaining.", guessesRemaining)
	} else {
		fmt.Print("You have 1 guess remaining. ")
	}
	fmt.Println("Your word to guess is:")
	printSpace(1)
	fmt.Println(getWordProgress(word, guessedLetters))
}

func getWordProgress(
	word string,
	guessedLetters map[rune]bool,
) string {
	result := ""
	for _, char := range word {
		if char == ' ' { //in case the word has spaces somehow
			result += " "
		} else if guessedLetters[unicode.ToLower(char)] { //comparing set to lowercase so upper and lower don't count seperately
			result += fmt.Sprintf("%c ", char)
		} else {
			result += "_ "
		}
	}

	return result
}

func getHangmanDrawing(hangmanState int) string {
	data, err := os.ReadFile(fmt.Sprintf("hangman_states/hangman%d", hangmanState))
	if err != nil {
		panic(err)
	}

	return string(data)
}
