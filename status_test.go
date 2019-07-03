package status_test

import (
	"log"
	"testing"

	status "github.com/dlarssonse/status"
)

func TestGetStatus(t *testing.T) {
	stat, err := status.GetStatus("test_status.json")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	if stat == nil {
		t.Errorf("Unable to get testing status.")
		return
	}

	log.Printf("%+v", stat)
}
