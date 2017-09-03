package dbsvc

import (
	"testing"
	"time"

	"scheduler/tracker"
)

func TestListEntity(t *testing.T) {
	jobs, err := ListEntity("", "", "")
	js := jobs.([]tracker.JobObject)
	if err != nil || len(js) < 0 {
		t.Errorf("test failed: %v", err)
	}
}

func TestUpdateEntity(t *testing.T) {
	set := map[string]interface{}{
		"type":   "2",
		"status": "0",
		"body":   "[]",
		"tag":    time.Now().Unix(),
	}
	err := UpdateEntity(1, set)
	if err != nil {
		t.Errorf("test failed: %v", err)
	}
}

func TestGetEntity(t *testing.T) {
	_, err := GetEntity(1)
	if err != nil {
		t.Errorf("test failed: %v", err)
	}
}
