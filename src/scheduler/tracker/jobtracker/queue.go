package jobtracker

import (
	"errors"

	"scheduler/tracker"
)

type Queue struct {
	queueChan chan interface{}
}

const QueueCache = 10

var queue *Queue

func init() {
	queue = &Queue{
		queueChan: make(chan interface{}, QueueCache),
	}
}

func (this *Queue) Pop() (interface{}, error) {
	job := <-this.queueChan
	return job, nil
}

//type assertion is needed only when you are a data user
func (this *Queue) Push(job interface{}) error {
	j, ok := job.(tracker.Job)
	if !ok {
		return errors.New("got error in type assertion")
	}
	this.queueChan <- j
	return nil
}

func (this *Queue) Length() (int, error) {
	l := len(this.queueChan)
	return l, nil
}
