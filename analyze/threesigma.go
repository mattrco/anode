package analyze

import (
	"errors"
	"math"

	"github.com/eclesh/welford"
	"github.com/myrid/anode.exp/data"
)

type ThreeSigma struct {
	stats   *welford.Stats
	input   chan data.Datapoint
	outputs []chan data.Datapoint
	// Stores the three latest values for calculating a moving average.
	tailbuf [3]float64
}

func (t *ThreeSigma) Init(input chan data.Datapoint, outputs []chan data.Datapoint) error {
	if len(outputs) != 4 {
		return errors.New("Must supply 4 output channels")
	}
	t.stats = welford.New()
	t.input = input
	t.outputs = outputs
	return nil
}

func (t *ThreeSigma) Run() {
	for d := range t.input {
		// Add value to distribution, update mean and stddev.
		t.stats.Add(d.Value)
		stddev := t.stats.Stddev()
		mean := t.stats.Mean()

		// If difference between MA and mean > 3 sigma, send to output.
		ma := t.movingAvg(d.Value)
		if math.Abs(ma-mean) > 3*stddev {
			t.outputs[0] <- d
		}

		// Output mean.
		t.outputs[1] <- data.Datapoint{Timestamp: d.Timestamp, Value: mean}

		// Output mean +/- 3 standard deviations.
		upper := t.stats.Mean() + 3*stddev
		lower := t.stats.Mean() - 3*stddev
		t.outputs[2] <- data.Datapoint{Timestamp: d.Timestamp, Value: upper}
		t.outputs[3] <- data.Datapoint{Timestamp: d.Timestamp, Value: lower}
	}
}

// movinAvg returns the mean of the latest three values.
// TODO: tidy up handle fewer than 3 values correctly.
func (t *ThreeSigma) movingAvg(latest float64) float64 {
	t.tailbuf[0] = t.tailbuf[1]
	t.tailbuf[1] = t.tailbuf[2]
	t.tailbuf[2] = latest
	ma := (t.tailbuf[0] + t.tailbuf[1] + t.tailbuf[2]) / 3
	return ma
}
