package jobtracker

import (
	"testing"
)

func TestPush(t *testing.T) {
	s := struct{}{}
	err := queue.Push(s)
	if err != nil {
		t.Errorf("test failed")
	}
}

func TestPop(t *testing.T) {
	s := struct{}{}
	_ = queue.Push(s)
	e, err := queue.Pop()
	if err != nil {
		t.Errorf("test failed")
	}
}
