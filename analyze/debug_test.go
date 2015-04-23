package analyze

import (
	"testing"

	"github.com/mattrco/anode/data"
)

func TestDebugInit(t *testing.T) {
	d := Debug{}
	var input chan data.Datapoint
	var outputs []chan data.Datapoint
	if err := d.Init(input, outputs); err != nil {
		t.Fatalf("Error initializing debug: %s", err)
	}
}

func TestDebugRun(t *testing.T) {
	input := make(chan data.Datapoint)
	output := make(chan data.Datapoint)
	outputs := []chan data.Datapoint{output}

	d := Debug{}
	if err := d.Init(input, outputs); err != nil {
		t.Fatalf("Error initializing debug: %s", err)
	}
	go d.Run()

	input <- data.Datapoint{Value: 1}
	result := <-output

	if result.Value != 1 {
		t.Fatalf("Expected 1, got %d", result.Value)
	}
}
