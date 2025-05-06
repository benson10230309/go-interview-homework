package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// question holds the data for one math problem
type Question struct {
	a, b    int
	op      string
	correct float32
}

var once sync.Once // only be executed once

// generateQuestion creates a random math question
func generateQuestion() Question {

	ops := []string{"+", "-", "*", "/"}
	op := ops[rand.Intn(len(ops))]
	a := rand.Intn(101)
	b := rand.Intn(101)
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
			b = rand.Intn(100) + 1 // avoid dividing by zero
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

func raiseYourHandToAnswer(name string, answer Question, chAnswer chan string, chOther chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	didWin := false
	time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
	once.Do(func() { // only be executed once
		chAnswer <- name
		fmt.Println("Student", name, ": The answer ", answer.a, answer.op, answer.b, " is =", answer.correct)
		didWin = true
	})
	// prevent winners from continuing to execute
	if !didWin {
		congratulation(name, chOther)
	}
}

func congratulation(name string, chOther chan string) {
	winner := <-chOther
	fmt.Println("Student", name, ":", winner, ", you win")
}

func teacherAsk(student []string) {
	chAnswer := make(chan string, 1)             // record winner
	chOther := make(chan string, len(student)-1) // give everyone the name of the winner

	var wg sync.WaitGroup

	fmt.Println("Teacher: Guys, are you ready?")

	time.Sleep(3 * time.Second)
	q := generateQuestion()
	fmt.Println("Teacher:", q.a, q.op, q.b, "= ?")

	for _, v := range student {
		wg.Add(1)
		go raiseYourHandToAnswer(v, q, chAnswer, chOther, &wg)
	}

	winner := <-chAnswer
	fmt.Println("Teacher:", winner, ", you are right!")

	for i := 0; i < len(student)-1; i++ {
		chOther <- winner // give everyone the name of the winner
	}

	wg.Wait()
}

func main() {
	student := []string{"A", "B", "C", "D", "E"}
	teacherAsk(student)
}
