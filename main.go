package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	PizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("received order #%d !\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++
		fmt.Printf("Making pizza #%d. it will take %d seconds .... \n", pizzaNumber, delay)
		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** w ran out of ingredients for pizza #%d", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** the cook quit while making pizza #%d", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("pizza order #%d is ready", pizzaNumber)
		}
		p := PizzaOrder{
			PizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
		return &p
	}
	return &PizzaOrder{
		PizzaNumber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	//  keep track of which pizza we are making
	var i = 0

	// run forever or until we receive a quit notification
	//try to make pizzas
	for {
		currentPizza := makePizza(i)
		//try to make a pizza
		//decision
		if currentPizza != nil {
			i = currentPizza.PizzaNumber
			select { //only for channels
			//we tried to make a pizza we sent something to the data channel
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <-pizzaMaker.quit:
				//close channels
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// print out a message
	color.Cyan("The Pizzeria is open for business")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzeria(pizzaJob)
	// create and run consumer
	for i := range pizzaJob.data {
		if i.PizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("order #%d is out for delivery", i.PizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("the curstomer is really mad")
			}
		} else {
			color.Cyan("done making pizzas")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** error closing channel", err)
			}
		}
	}

	//print out the ending message
	color.Cyan("------------------")
	color.Cyan("done for the day")

	color.Cyan("we made %d pizzas, but failed to make %d, with %d attempts in total ", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("it was an awful day")
	case pizzasFailed >= 6:
		color.Red("it was not a very good day")
	case pizzasFailed >= 4:
		color.Yellow(" it was an okay day")
	case pizzasFailed >= 2:
		color.Yellow(" it was a pretty good day")
	default:
		color.Green("it was a great day")
	}

}
