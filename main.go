package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

func main() {
	pin0 := rpio.Pin(0)
	pin1 := rpio.Pin(1)

	pin0.Input()
	pin1.Input()

	for {
		val0 := pin0.Read()
		val1 := pin1.Read()
		log.Printf("Value pin0:%v pin1:%v", val0, val1)
		time.Sleep(time.Second)
	}
}
