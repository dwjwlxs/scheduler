package utils

import "testing"

func TestNearestFuture(t *testing.T) {
	clock := " 15,0,30 14,13 * * * "
	seconds, err := NearestFuture(clock)
	if err != nil || seconds < 0 {
		t.Errorf("test failed: %#v", seconds)
	}
}
