package analyze

import (
	"fmt"
	"math"

	"github.com/eclesh/welford"
	"github.com/mattrco/anode/data"
)

type ThreeSigma struct {
	stats  *welford.Stats
	input  chan data.Datapoint
	output chan data.Datapoint
	// Stores the three latest values for calculating a moving average.
	tailbuf [3]float64
}

func (t *ThreeSigma) Init(input chan data.Datapoint, output chan data.Datapoint) error {
	t.stats = welford.New()
	t.input = input
	t.output = output
	return nil
}

func (t *ThreeSigma) Run() {
	// It's useful to know how many metrics have been processed, e.g.
	// for skipping steps that require a certain (such as moving avg.)
	dataCount := 0
	for d := range t.input {
		// Add value to distribution, update mean and stddev.
		t.stats.Add(d.Value)
		stddev := t.stats.Stddev()
		mean := t.stats.Mean()

		t.UpdateTailbuf(d.Value)
		dataCount += 1
		if dataCount < 3 {
			continue
		}

		// If difference between MA and mean > 3 sigma, send to output.
		ma := (t.tailbuf[0] + t.tailbuf[1] + t.tailbuf[2]) / 3
		if math.Abs(ma-mean) > 3*stddev {
			t.output <- data.Datapoint{
				MetricName: fmt.Sprintf("anode.threesig.%s.anomalous", d.MetricName),
				Timestamp:  d.Timestamp,
				Value:      ma,
				IsAnamoly:  true,
			}
		}

		// Output mean.
		t.output <- data.Datapoint{
			MetricName: fmt.Sprintf("anode.threesig.%s.mean", d.MetricName),
			Timestamp:  d.Timestamp,
			Value:      mean,
			IsAnamoly:  false,
		}

		// Output mean +/- 3 standard deviations.
		t.output <- data.Datapoint{
			MetricName: fmt.Sprintf("anode.threesig.%s.upper", d.MetricName),
			Timestamp:  d.Timestamp,
			Value:      t.stats.Mean() + 3*stddev,
			IsAnamoly:  false,
		}
		t.output <- data.Datapoint{
			MetricName: fmt.Sprintf("anode.threesig.%s.lower", d.MetricName),
			Timestamp:  d.Timestamp,
			Value:      t.stats.Mean() - 3*stddev,
			IsAnamoly:  false,
		}
	}
}

func (t *ThreeSigma) UpdateTailbuf(latest float64) {
	t.tailbuf[0] = t.tailbuf[1]
	t.tailbuf[1] = t.tailbuf[2]
	t.tailbuf[2] = latest
}
