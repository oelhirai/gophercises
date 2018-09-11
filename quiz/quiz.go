package main

import (
	"fmt"
	"encoding/csv"
	"os"
	"flag"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

const finalResult string = "You scored %d out of %d\n"
const problemFormat string = "Problem #%d: %s = "
const csvFlagDescription string = "a csv file in the format of 'question,answer'"
const limitFlagDescription string = "the time limit for the quiz in seconds"

var csvName string = "problems.csv"
var limit int = 5

func main() {
	readFlags()
	lines := readCSV()
	problems := buildQuestionnaire(lines)
	testUser(problems)
}

func readFlags() {
	flag.StringVar(&csvName, "csv", csvName, csvFlagDescription)
	flag.IntVar(&limit, "limit", limit, limitFlagDescription)
	flag.Parse()
}

func readCSV() [][]string {
	csvFile, _ := os.Open(csvName)
	reader := csv.NewReader(csvFile)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", csvName)
		os.Exit(1)
	}
	return lines
}

func buildQuestionnaire(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem { 
			q: line[0],
			a: strings.TrimSpace(line[1]),
		} 
	}
	return problems
}

func testUser(problems []problem) {
	timerChannel := time.After(time.Duration(limit) * time.Second)
	answerChannel := make(chan string)
	correctAnswers := 0

	for i, p := range problems {
		fmt.Printf(problemFormat, i + 1, p.q)
		go fetchAnswer(answerChannel)

		select {
		case <- timerChannel:
			fmt.Println()
			fmt.Printf(finalResult, correctAnswers, len(problems))
			os.Exit(0)
		case answer := <- answerChannel:
			if answer == p.a {
				correctAnswers++
			}
		}
	}

	fmt.Printf(finalResult, correctAnswers, len(problems))
}

func fetchAnswer(c chan string) {
	var answer string 
	fmt.Scanf("%s\n", &answer)
	c <- answer
}

