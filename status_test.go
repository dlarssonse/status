package status_test

import (
	"encoding/json"
	"log"
	"testing"

	status "github.com/dlarssonse/status"
)

func TestGetStatus(t *testing.T) {

	stat, err := status.ReadFile("test_status.json")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	if stat == nil {
		t.Errorf("Unable to get testing status.")
		return
	}
	if err = status.GetStatus(stat); err != nil {
		t.Errorf("%s", err)
		return
	}

	data, _ := json.Marshal(stat)
	log.Printf("%s", data)
}
