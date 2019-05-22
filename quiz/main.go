package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	csvFile := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer' (default \"problems.csv\"")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFile) // For read access.
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	problems := parseProblemsCsv(file)
	numberOfCorrectAnwers := answerProblems(problems, *timeLimit)
	fmt.Printf("You scored %d of %d.\n", numberOfCorrectAnwers, len(problems))
}

func answerProblems(problems []problem, timeLimit int) int {
	quit := make(chan bool)
	numberOfCorrectAnwers := 0
	go sleep(quit, timeLimit)
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanln(&answer)
			answerCh <- answer
		}()

		select {
		case <-quit:
			fmt.Println()
			return numberOfCorrectAnwers
		case answer := <-answerCh:
			if problem.answer == answer {
				numberOfCorrectAnwers++
			}
		}
	}

	return numberOfCorrectAnwers
}

func sleep(quit chan bool, timeLimit int) {
	time.Sleep(time.Duration(timeLimit) * time.Second)
	quit <- true
}

func parseProblemsCsv(file *os.File) []problem {
	r := csv.NewReader(file)
	problems := []problem{}
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		problem := problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
		problems = append(problems, problem)
	}
	return problems
}

type problem struct {
	question string
	answer   string
}
