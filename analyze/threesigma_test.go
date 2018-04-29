package analyze

import (
	"sync"
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
	var wg sync.WaitGroup
	input := make(chan data.Datapoint, 3)
	output := make(chan data.Datapoint)

	ts := ThreeSigma{}
	if err := ts.Init(input, output); err != nil {
		t.Fatalf("Error initializing three sigma: %s", err)
	}
	go ts.Run()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			input <- data.Datapoint{Value: 1}
		}
	}()

	wg.Wait()
	result := <-output

	if result.Value != 1 {
		t.Fatalf("Expected 1, got %f", result.Value)
	}
}
