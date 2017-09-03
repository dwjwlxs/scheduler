package tasktracker

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"scheduler/tracker"
	"scheduler/tracker/jobtracker"
	"scheduler/worker"
)

const (
	MaxGoroutines    = 20
	TrackerSleepTime = 1
)

var (
	mu        sync.Mutex                      // guards a job only matches a task tracker
	limitChan = make(chan int, MaxGoroutines) //total goroutine in task tracker scope.
)

/**
each task tracker has a private task queue,
which make sure tasks in a job is done in right sequence.
*/
type Tracker struct {
	name string
	wg   sync.WaitGroup
}

func NewTracker() (*Tracker, error) {
	host, _ := os.Hostname()
	return &Tracker{
		name: host,
	}, nil
}

func (this *Tracker) Run() {
	this.wg.Add(1)
	go this.Reserve()
	this.wg.Wait()
}

/**
a goroutine, try to reserve a job from job queue.
*/
func (this *Tracker) Reserve() {
	defer this.wg.Done()
	for {
		fmt.Printf("%v: tasktracker@%v is trying to reserve a job from job queue\n", time.Now().String(), this.name)
		//one job, one goroutine
		limitChan <- 1 //blocked when reach the max num of goroutine
		mu.Lock()      //serialize Reserve() operation, not a must
		job, rerr := jobtracker.Reserve()
		mu.Unlock()
		if rerr != nil {
			<-limitChan
			continue
		}
		go this.Deliver(job)
		time.Sleep(TrackerSleepTime * time.Second)
	}
}

/**
a goroutine to deliver each task to a specified worker, e.g. a mailer.
*/
func (this *Tracker) Deliver(job interface{}) {
	fmt.Printf("%v: tasktracker@%v has delivered a job\n", time.Now().String(), this.name)
	ts, uerr := this.unpackJob(job)
	if uerr == nil {
		tasks := ts.([]tracker.Task)
		for _, task := range tasks {
			w := worker.Instance(task.Worker, task.Fields)
			_, _ = w.Execute() //synchronized operation guards task is done in right sequence.
		}
	}
	<-limitChan
}

func (this *Tracker) unpackJob(job interface{}) (interface{}, error) {
	j, ok := job.(tracker.Job)
	if ok {
		if j.OrderReady == false {
			//do ordering thing
			//i think it is better to do task sorting in job tracker.
			return j.Tasks, nil
		}
		return j.Tasks, nil
	}
	return nil, errors.New("error occured when assert type:Job")
}
