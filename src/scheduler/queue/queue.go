//ToDo
package queue

import (
	"fmt"
	"os"
)

type Queue struct {
	tube       string
	middleware interface{}
}

var q *Queue

func init() {
	var err error
	if q, err = NewQueue(""); err != nil {
		fmt.Println("err occured when New Queue: ", err)
		os.Exit(1)
	}
}

func NewQueue(name string) (*Queue, error) {
	return &Queue{
		tube: name,
	}, nil
}

func (this *Queue) Pop() string {
	return "middleware"
}

func (this *Queue) Push() error {
	return nil
}

func (this *Queue) List() map[string]interface{} {
	return map[string]interface{}{}
}
