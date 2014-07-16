package analyze

import (
	"testing"

	"github.com/mattrco/anode.exp/data"
)

func TestThreeSigmaInit(t *testing.T) {
	ts := ThreeSigma{}

	var input chan data.Datapoint
	var outputs [4]chan data.Datapoint
	if err := ts.Init(input, outputs); err != nil {
		t.Fatalf("Error initializing three sigma: %s", err)
	}
}

func TestThreeSigmaRun(t *testing.T) {

	input := make(chan data.Datapoint)
	o1 := make(chan data.Datapoint, 1)
	o2 := make(chan data.Datapoint, 1)
	o3 := make(chan data.Datapoint, 1)
	o4 := make(chan data.Datapoint, 1)
	outputs := [4]chan data.Datapoint{o1, o2, o3, o4}

	ts := ThreeSigma{}
	if err := ts.Init(input, outputs); err != nil {
		t.Fatalf("Error initializing three sigma: %s", err)
	}
	go ts.Run()

	input <- data.Datapoint{Value: 1}
	result := <-o2

	if result.Value != 1 {
		t.Fatalf("Expected 1, got %d", result.Value)
	}
}
