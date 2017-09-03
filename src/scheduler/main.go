package main

import (
	"fmt"
	"sync"
	"time"

	"scheduler/tracker/jobtracker"
	"scheduler/tracker/tasktracker"
)

var wg sync.WaitGroup

func main() {
	supvsr, _ := jobtracker.NewSupervisor()
	go supvsr.Run()

	consumer, _ := tasktracker.NewTracker()
	go consumer.Run()
	wg.Add(1)
	go func() {
		for {
			fmt.Printf("%v: Main goroutine is running\n", time.Now().String())
			time.Sleep(time.Second * 10)
		}
		wg.Done()
	}()
	wg.Wait()
}
