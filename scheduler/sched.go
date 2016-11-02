package scheduler

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Schedule(duration time.Duration) <-chan bool {
	timeout := make(chan bool, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	c := make(chan bool, 0)
	go func() {
		for {
			go func() {
				time.Sleep(duration)
				timeout <- true
			}()
			select {
			case <-sig:
				close(c)
				break
			case <-timeout:
				c <- true
				continue
			}
		}
	}()
	return c
}

func Tasks(master <-chan bool, num_tasks int) []<-chan bool {
	var chans []chan bool = make([]chan bool, num_tasks)
	var res_chans []<-chan bool = make([]<-chan bool, num_tasks)

	for i := 0; i < num_tasks; i++ {
		chans[i] = make(chan bool, 0)
		res_chans[i] = chans[i]
	}

	go func() {
		for _ = range master {
			for i := 0; i < num_tasks; i++ {
				chans[i] <- true
			}
		}
		for i := 0; i < num_tasks; i++ {
			close(chans[i])
		}
	}()

	return res_chans
}
