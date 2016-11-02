package main

import (
	"fmt"
	"github.com/mildred/anemometre/scheduler"
	"github.com/stianeikeland/go-rpio"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer rpio.Close()

	pin0 := rpio.Pin(17)
	pin1 := rpio.Pin(18)

	var wg sync.WaitGroup

	log.Println("Wind direction on rPI GPIO_0")
	log.Println("Wind vane on rPI GPIO_1")

	sched := scheduler.Tasks(scheduler.Schedule(100*time.Microsecond), 2)
	val0 := readGPIO(sched[0], pin0)

	speedResource := NewSpeed(&wg, sched[1], pin1)
	go serveHTTP(speedResource)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range val0 {
			fmt.Printf("%v, pin0: %v\n", v.Time, v.Value)
		}
	}()

	wg.Wait()
	log.Println("Finished")
}

type Speed struct {
	lock  sync.Mutex
	cond  *sync.Cond
	speed atomic.Value
}

func NewSpeed(wg *sync.WaitGroup, sched <-chan bool, pin rpio.Pin) *Speed {
	var speedResource Speed
	speedResource.cond = sync.NewCond(&speedResource.lock)

	wg.Add(1)
	val := readGPIO(sched, pin)
	go func() {
		for s := range speed(val) {
			speedResource.speed.Store(s)
			speedResource.cond.Broadcast()
			fmt.Printf("speed: %f\n", s)
		}
		defer wg.Done()
	}()

	return &speedResource
}

func (s *Speed) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for {
		if speed, ok := s.speed.Load().(float64); ok {
			res.Write([]byte(fmt.Sprintf("%f\n", speed)))
			if f, ok := res.(http.Flusher); ok {
				f.Flush()
			} else {
				log.Println("Damn, no flush")
			}
		}
		s.cond.Wait()
	}
}

type Event struct {
	Time  time.Time
	Value rpio.State
}

func speed(events <-chan *Event) <-chan float64 {
	c := make(chan float64, 1)
	go func() {
		var last *Event
		for ev := range events {
			if last == nil {
				last = ev
			} else if last.Value == ev.Value {
				var duration time.Duration = ev.Time.Sub(last.Time)
				c <- 1.0 / duration.Seconds()
				last = ev
			}
		}
		close(c)
	}()
	return c
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
		close(c)
	}()

	return c
}

func serveHTTP(speed *Speed) {
	var s http.Server

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHTTPRequest)
	mux.Handle("/speed", speed)
	mux.HandleFunc("/fetch.js", serveFetch)

	s.Handler = mux
	s.Addr = ":80"
	log.Printf("Serve HTTP on %s", s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

const html string = `
<!DOCTYPE html>
<html>
<head>
<script src="fetch.js"></script>
<script type="text/javascript">
	var xhr = new XMLHttpRequest();
	xhr.previous_text = '';

	xhr.onreadystatechange = function()
	{
		if (this.readyState == 3 && this.status == 200)
		{
			var new_response = xhr.responseText.substring(xhr.previous_text.length);
			xhr.previous_text = xhr.responseText;
			document.querySelector("span.speed").textContent = new_response;
		}
	}
	xhr.open("get", "/speed", true);
	xhr.send();

</script>
</head>
<body>
<p>Speed: <span class="speed">N/A</span></p>
</body>
</html>
`

func serveHTTPRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.Write([]byte(html))
}

func serveFetch(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/javascript")
	res.Write([]byte(html_fetch))
}
