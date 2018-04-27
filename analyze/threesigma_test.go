package analyze

import (
	"testing"

	"github.com/mattrco/anode/data"
)

func TestThreeSigmaInit(t *testing.T) {
	ts := ThreeSigma{}

	var input chan data.Datapoint
	var output chan data.Datapoint
	if err := ts.Init(input, output); err != nil {
		t.Fatalf("Error initializing three sigma: %s", err)
	}
}

func TestThreeSigmaRun(t *testing.T) {
	input := make(chan data.Datapoint, 1)
	output := make(chan data.Datapoint, 1)

	ts := ThreeSigma{}
	if err := ts.Init(input, output); err != nil {
		t.Fatalf("Error initializing three sigma: %s", err)
	}
	go ts.Run()

	input <- data.Datapoint{Value: 1}
	result := <-output

	if result.Value != 1 {
		t.Fatalf("Expected 1, got %d", result.Value)
	}
}
