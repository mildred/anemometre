package main

import (
	"fmt"
	"github.com/mildred/anemometre/scheduler"
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

	sched := scheduler.Tasks(scheduler.Schedule(100*time.Microsecond), 2)
	val0 := readGPIO(sched[0], pin0)
	val1 := readGPIO(sched[1], pin1)

	go func() {
		for v := range val0 {
			fmt.Printf("%v, pin0: %v\n", v.Time, v.Value)
		}
	}()

	for v := range val1 {
		fmt.Printf("%v, pin1: %v\n", v.Time, v.Value)
	}
}

type Event struct {
	Time  time.Time
	Value rpio.State
}

func readGPIO(scheduler <-chan bool, pin rpio.Pin) <-chan *Event {
	c := make(chan *Event, 1)
	pin.Input()
	var oldval rpio.State = pin.Read()

	go func() {
		for _ = range scheduler {
			st := pin.Read()
			if st != oldval {
				oldval = st
				c <- &Event{time.Now(), st}
			}
		}
	}()

	return c
}
