package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer rpio.Close()

	pin0 := rpio.Pin(17)
	pin1 := rpio.Pin(18)

	pin0.Input()
	pin1.Input()

	for {
		val0 := pin0.Read()
		val1 := pin1.Read()
		log.Printf("Value pin0:%v pin1:%v", val0, val1)
		time.Sleep(100 * time.Microsecond)
	}
}

func readGPIO()
