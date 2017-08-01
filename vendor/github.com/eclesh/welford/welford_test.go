package welford

import (
	"math"
	"math/rand"
	"testing"
)

func TestSimple(t *testing.T) {
	s := New()
	for i := 1; i < 6; i++ {
		s.Add(float64(i))
	}
	if s.Count() != 5 {
		t.Fatalf("count: expected 5, got %d\n", s.Count())
	}
	if s.Min() != 1 {
		t.Fatalf("min: expected 0.0, got %f\n", s.Min())
	}
	if s.Max() != 5 {
		t.Fatalf("max: expected 5.0, got %f\n", s.Max())
	}
}

func TestNormal(t *testing.T) {
	s := New()
	for i := 0; i < 1000000; i++ {
		s.Add(rand.NormFloat64())
	}
	// Allow .1% of error (.1% arbitrarily chosen)
	if math.Abs(s.Mean()) > 0.001 {
		t.Fatalf("mean: expected 0.0, got %f\n", s.Mean())
	}
	if s.Variance()-1 > 0.001 {
		t.Fatalf("variance: expected 1.0, got %f\n", s.Variance())
	}
	if s.Stddev()-1 > 0.001 {
		t.Fatalf("stddev: expected 1.0, got %f\n", s.Stddev())
	}
}

func TestVariance(t *testing.T) {
	s := New()
	if s.Variance() != 0 {
		t.Fatalf("variance: expected 0, got %f\n", s.Variance())
	}
}

func TestReset(t *testing.T) {
	s := New()
	for i := 0; i < 1000000; i++ {
		s.Add(rand.NormFloat64())
	}
	s.Reset()
	if s.Mean() != 0 {
		t.Fatalf("mean: expected 0.0, got %f\n", s.Mean())
	}
	if s.Variance() != 0 {
		t.Fatalf("variance: expected 0.0, got %f\n", s.Variance())
	}
	if s.Stddev() != 0 {
		t.Fatalf("stddev: expected 0.0, got %f\n", s.Stddev())
	}
}
