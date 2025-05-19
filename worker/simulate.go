package worker

import (
	"log"
	"math/rand"
	"time"
)

func Simulate() bool {
	sim := rand.Intn(4) // 0: crash, 1: slow, 2/3: normal
	if sim == 0 {
		log.Println("Simulated crash")
		return false
	} else if sim == 1 {
		log.Println("Simulated delay")
		time.Sleep(12 * time.Second)
	}
	return true
}
