package analyze

import (
	"github.com/myrid/anode.exp/data"
)

type Debug struct {
	input   chan data.Datapoint
	outputs []chan data.Datapoint
}

func (d *Debug) Init(input chan data.Datapoint, outputs []chan data.Datapoint) error {
	d.input = input
	d.outputs = outputs
	return nil
}

func (d *Debug) Run() {
	for dp := range d.input {
		for _, output := range d.outputs {
			output <- dp
		}
	}
}
