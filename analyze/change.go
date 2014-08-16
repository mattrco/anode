package analyze

import (
	"fmt"
	"math"

	"github.com/dgryski/go-change"
	"github.com/golang/glog"
	"github.com/mattrco/anode.exp/data"
)

type Change struct {
	stream *change.Stream
	input  chan data.Datapoint
	output chan data.Datapoint
}

func (c *Change) Init(input chan data.Datapoint, output chan data.Datapoint,
	windowSize int, minSample int, blockSize int) error {
	c.stream = change.NewStream(windowSize, minSample, blockSize, 0.995)
	c.input = input
	c.output = output
	return nil
}

func (c *Change) Run() {
	for d := range c.input {
		cp := c.stream.Push(d.Value)
		if cp != nil {
			diff := math.Abs(cp.Difference / cp.Before.Mean())
			if cp.Difference != 0 && diff > 0.06 {
				glog.Info(cp.Difference)
				c.output <- data.Datapoint{
					MetricName: fmt.Sprintf("anode.change.%s.difference", d.MetricName),
					Timestamp:  d.Timestamp,
					Value:      cp.Difference,
				}
			}
		}
	}
}
