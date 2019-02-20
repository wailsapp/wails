package cmd

import (
	"fmt"
	"strconv"
)

// Prompt asks the user for a value
func Prompt(question string, defaultValue ...string) string {
	var answer string
	haveDefault := len(defaultValue) > 0 && defaultValue[0] != ""

	if haveDefault {
		question = fmt.Sprintf("%s (%s)", question, defaultValue[0])
	}
	fmt.Printf(question + ": ")
	fmt.Scanln(&answer)
	if haveDefault {
		if len(answer) == 0 {
			answer = defaultValue[0]
		}
	}
	return answer
}

// PromptRequired calls Prompt repeatedly until a value is given
func PromptRequired(question string, defaultValue ...string) string {
	for {
		result := Prompt(question, defaultValue...)
		if result != "" {
			return result
		}
	}
}

// PromptSelection asks the user to choose an option
func PromptSelection(question string, options []string) int {

	fmt.Println(question + ":")
	for index, option := range options {
		fmt.Printf("  %d: %s\n", index+1, option)
	}

	fmt.Println()
	selectedValue := -1

	for {
		choice := Prompt("Please choose an option")

		// index
		number, err := strconv.Atoi(choice)
		if err == nil {
			if number > 0 && number <= len(options) {
				selectedValue = number - 1
				break
			} else {
				continue
			}
		}

	}

	return selectedValue
}
