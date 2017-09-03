package jobtracker

import (
	"testing"
)

func TestPut(t *testing.T) {
	err := Put(struct{}{})
	if err != nil {
		t.Errof("test failed")
	}
}
