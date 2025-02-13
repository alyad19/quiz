package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q, a string
}

func problems(filename string) []problem {
	file, err := os.Open(filename)

	if err != nil {
		exit(fmt.Sprintf("Can not open file: %s\n", filename))
	}
	defer file.Close()
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	ret := make([]problem, len(lines))

	if err != nil {
		exit("Can not read file")
	}
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret

}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func ask(problem problem, timeLimit int) int {
	fmt.Printf("%s = ", problem.q)
	answerCh := make(chan string)

	go func() {

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		answerCh <- input
		temp := <-answerCh
		fmt.Println(temp)

	}()

	select {
	case userAnswer := <-answerCh:
		if problem.a == userAnswer {
			return 1
		}

	case <-time.After(time.Duration(timeLimit) * time.Second):
		fmt.Println("Time's up!")
		<-answerCh

	}
	return 0
}

func quiz(filename string, timeLimit int) (int, int) {
	problem_list := problems(filename)
	score := 0
	count := 0

	for _, problem := range problem_list {
		count++
		userAnswer := ask(problem, timeLimit)
		score += userAnswer

	}
	return score, count

}

func main() {
	filename := flag.String("csvfilename", "problems.csv", "A csv file in the format 'q,a'")
	timeLimit := flag.Int("limit", 5, "the time limit for the quiz in seconds")
	flag.Parse()
	score, count := quiz(*filename, *timeLimit)
	fmt.Printf("You score %d out of %d", score, count)

}
