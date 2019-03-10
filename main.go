package main

import (
	"fmt"
	"time"
)

// Rollar coaster Capacity
const capacity = 10
const passengerNum = 35

func openCoaster() chan bool {
	time.Sleep(2 * time.Second)
	open := make(chan bool)
	go func() {
		open <- true
	}()
	return open
}

func lineCustomer() []chan int {
	var customers []chan int
	for i := 0; i < passengerNum; i++ {
		customers = append(customers, waitingInLine(i))
	}
	return customers
}

func waitingInLine(id int) chan int {
	ch := make(chan int)
	go func() {
		ch <- id
	}()
	return ch
}

// Generator Go Concurrency Pattern
func canRide(seatsNeeded int) chan bool {
	canRide := make(chan bool)
	go func() {
		for i := 0; i < seatsNeeded; i++ {
			canRide <- true
		}
	}()
	return canRide
}

// use a fanin go concurrecy pattern
func getInRollar(customers []chan int, canRide chan bool) chan int {
	fmt.Println("New Rounds")

	chToCoaster := make(chan int, capacity)

	// This does not impose an order... like a mob
	for _, chRider := range customers {
		// Block the goChannel Worker to get into the Coaster
		// if the Coaster is running
		go func(chRider chan int) {
			<-canRide
			id := <-chRider
			fmt.Printf("Passenger id %d gets into the coaster\n", id)
			chToCoaster <- id
		}(chRider)
		//This part is really important because if not using this case,
		// each go function will just take the variable they see instead of the one we want to assign to them
	}

	// This codes actually impose an order
	// go func() {
	// 	for _, chRider := range customers {
	// 		// Block the goChannel Worker to get into the Coaster
	// 		// if the Coaster is running
	// 		<-canRide
	// 		id := <-chRider
	// 		fmt.Printf("Passenger id %d gets into the coaster\n", id)
	// 		chToCoaster <- id
	// 	}
	// }()
	return chToCoaster
}

// handling the actuall play of the rollar coaster
func rollarCoasterHandler(coasterChan chan int) {
	for {
		finished := time.After(3 * time.Second)
		select {
		case pid := <-coasterChan:
			fmt.Printf("Passenger %d is playing\n", pid)
		case <-finished:
			fmt.Println("The RollarCoaster return...")
			return
		}
	}
}

// Check if this will be the last Ride
func lastRide(finishRide int) bool {
	return passengerNum-finishRide < capacity
}

func main() {
	// Open the RollarCoaster
	open := openCoaster()

	// Put the customer inline
	customers := lineCustomer()

	fmt.Println("Waiting for coaster to open.....")
	// block until the Coaster is open
	<-open

	finishRide := 0
	// Each for loop is a round
	for finishRide < passengerNum {
		var chCanRide chan bool

		if lastRide(finishRide) {
			seatsNeeded := passengerNum - finishRide
			chCanRide = canRide(seatsNeeded)
		} else {
			chCanRide = canRide(capacity)
		}
		// Try to move the customers into the rollar coaster
		chToCoaster := getInRollar(customers, chCanRide)

		rollarCoasterHandler(chToCoaster)

		// Incrementing the count of how many people have finished riding
		if lastRide(finishRide) {
			break
		} else {
			// Kicked off those who have finished riding the coaster
			customers = customers[capacity:]
			finishRide += capacity
		}
	}
	fmt.Println("All passengers have finished their rides, I will close")
}
