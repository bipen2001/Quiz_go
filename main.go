package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type problem struct {
	question string
	answer   string
}

var (
	score           = 0
	defautTime      = 1
	quesTime        = 1
	quesTimer       *time.Timer
	fileName        = flag.String("csv", "problems.csv", "a csv file in format of questions,answers")
	timeLimit       = flag.Int("limit", 30, "this is the time limit for the quiz in seconds")
	m               sync.Mutex
	checkArr        []problem
	totalNumberQues int = 5
	start           string
)

func main() {

	flag.Parse()

	fmt.Print("Are u ready to give the Quiz? If yes press Y OR YES  \n")
	fmt.Scan(&start)

	switch start {
	case "Y":
		organizeQuiz()
	case "YES":
		organizeQuiz()
	case "yes":
		organizeQuiz()
	case "y":
		organizeQuiz()

	default:
		fmt.Println("Entered wrong input")
	}

}

func organizeQuiz() {
	var (
		lines              [][]string
		innerSlice         []string
		answerindex, index int
		questionPart       string
	)
	f, err := os.Open(*fileName)

	if err != nil {
		log.Fatal(err)
	}

	sampleRegexp := regexp.MustCompile(`\d`)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		for i, _ := range line {
			if line[i] == '?' {
				index = i
				break
			}
		}
		questionPart = strings.TrimSpace(line[:index+1])

		line = line[index+1:]
		for i, _ := range line {
			if sampleRegexp.MatchString(string(line[i])) {
				answerindex = i
				break
			}
		}

		answerpart := strings.TrimSpace(line[answerindex:])

		innerSlice = append(innerSlice, questionPart, answerpart)
		lines = append(lines, innerSlice)
		defer f.Close()

	}

	problems := parseLines(lines)
	answerCh := make(chan string)
	startQuiz(problems, answerCh)
	fmt.Printf("you answered %d Questions correctly out of %d\n", score, totalNumberQues)
}

func startQuiz(problems []problem, answerCh chan string) {
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemloop:
	for i := 0; i < totalNumberQues; i++ {

		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(problems))
		fmt.Printf("Problem #%d: %s\n", i+1, problems[r].question)

		quesTimer = time.NewTimer(time.Duration(quesTime) * time.Second)

		go func(answerCh chan string) {
			m.Lock()

			defer m.Unlock()
			var answer string

			_, err := fmt.Scanf("%s\n", &answer)

			if err != nil {
				panic(err)
			}

			answerCh <- answer

		}(answerCh)
		select {
		case answer := <-answerCh:

			if answer == problems[r].answer {
				score++
			}
			problems = append(problems[:r], problems[r+1:]...)

		case <-timer.C:
			break problemloop
		case <-quesTimer.C:
			problems = append(problems[:r], problems[r+1:]...)

			continue

		}

	}
}
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret
}
