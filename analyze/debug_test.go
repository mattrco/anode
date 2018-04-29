package analyze

import (
	"sync"
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
	var wg sync.WaitGroup
	input := make(chan data.Datapoint, 3)
	output := make(chan data.Datapoint)
	outputs := []chan data.Datapoint{output}

	d := Debug{}
	if err := d.Init(input, outputs); err != nil {
		t.Fatalf("Error initializing debug: %s", err)
	}
	go d.Run()

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
