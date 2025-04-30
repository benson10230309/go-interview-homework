package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var ansLock sync.Mutex
var once sync.Once

// question holds the data for one math problem
type Question struct {
	a, b    int
	op      string
	correct float32
}

// generateQuestion creates a random math question
func generateQuestion() Question {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	ops := []string{"+", "-", "*", "/"}
	op := ops[r.Intn(len(ops))]
	a := r.Intn(101)
	b := r.Intn(101)

	var correct float32

	switch op {
	case "+":
		correct = float32(a + b)
	case "-":
		correct = float32(a - b)
	case "*":
		correct = float32(a * b)
	case "/":
		for b == 0 {
			b = r.Intn(101) // avoid dividing by zero
		}
		correct = float32(a) / float32(b)
	}

	return Question{
		a:       a,
		b:       b,
		op:      op,
		correct: correct,
	}

}

func raiseYourHandToAnswer(name string, answer Question, chRespondent chan string, wg *sync.WaitGroup, chAnswer chan float32, didwin *bool) {
	defer wg.Done()

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	time.Sleep(time.Duration(r.Intn(3)+1) * time.Second)
	hitRate := r.Intn(3) // Answer correctly (1/3 chance), otherwise answer incorrectly on purpose
	var myAnswer float32
	myAnswer = answer.correct

	ansLock.Lock()
	defer ansLock.Unlock()
	if hitRate != 0 || *didwin { // If the student misses the hit, then the answer is wrong.
		for myAnswer == answer.correct { // The second hit will also come here
			myAnswer = float32(r.Intn(101))
		}
		chRespondent <- name
		chAnswer <- myAnswer
		fmt.Println("Student", name, ": The answer ", answer.a, answer.op, answer.b, " is =", myAnswer)
	} else {
		once.Do(func() {
			fmt.Println("Student", name, ": The answer ", answer.a, answer.op, answer.b, " is =", answer.correct)
			*didwin = true
		})
		chRespondent <- name
		chAnswer <- answer.correct
	}
}

func main() {
	student := []string{"A", "B", "C", "D", "E"}

	chRespondent := make(chan string, len(student)) // record respondent
	chAnswer := make(chan float32, len(student))

	var wg sync.WaitGroup
	var wg1 sync.WaitGroup // Congratulations to those who answered correctly
	var winner string
	var haveCorrectAnswer bool
	didwin := false

	fmt.Println("Teacher: Guys, are you ready?")

	time.Sleep(3 * time.Second)
	q := generateQuestion()
	fmt.Println("Teacher:", q.a, q.op, q.b, "= ?")

	for _, v := range student {
		wg.Add(1)
		go raiseYourHandToAnswer(v, q, chRespondent, &wg, chAnswer, &didwin)
	}

	haveCorrectAnswer = false
	count := len(student)

	for range student {
		respondent := <-chRespondent
		answer := <-chAnswer
		if answer == q.correct && !haveCorrectAnswer { // Check if the returned answer is correct or not
			winner = respondent
			fmt.Println("Teacher:", winner, ", you are right!")
			haveCorrectAnswer = true
		} else if answer != q.correct {
			fmt.Println("Teacher:", respondent, ", you are wrong.")
			count--
		}
	}

	if count == 0 {
		fmt.Println("Teacher: Boooo~ Answer is ", q.correct)
	}

	if haveCorrectAnswer {
		for _, i := range student {
			if i != winner {
				wg1.Add(1)
				go congratulations(i, winner, &wg1) // Tell everyone who is the winner
			}
		}
		wg1.Wait()
	}
	wg.Wait()
}

func congratulations(name string, winner string, wg1 *sync.WaitGroup) {
	defer wg1.Done()
	fmt.Println("Student", name, ":", winner, ", you win.") // Say congratulations to the winner
}
