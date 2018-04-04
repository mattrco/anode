// Package stats implements Welford's one-pass algorithm for computing
// the mean and variance of a set of numbers. For more information see
// Knuth (TAOCP Vol 2, 3rd ed, pg 232).
package welford

import (
	"math"
)

// A Stats gathers statistics on values fed into it.
type Stats struct {
	count             uint64
	min, max, mean, s float64
}

// New returns a new Stats.
func New() *Stats {
	return &Stats{}
}

// Min returns the minimum value seen.
func (s *Stats) Min() float64 {
	return s.min
}

// Max returns the maximum value seen.
func (s *Stats) Max() float64 {
	return s.max
}

// Count returns the total number of values seen.
func (s *Stats) Count() uint64 {
	return s.count
}

// Mean returns the mean of all values seen.
func (s *Stats) Mean() float64 {
	return s.mean
}

// Variance returns the variance of all values seen.
func (s *Stats) Variance() float64 {
	if s.count > 1 {
		return s.s / float64(s.count-1)
	}
	return 0.0
}

// Stddev returns the standard deviation of all values seen.
func (s *Stats) Stddev() float64 {
	return math.Sqrt(s.Variance())
}

// Reset sets all counters to zero.
func (s *Stats) Reset() {
	s.count = 0
	s.min = 0.0
	s.max = 0.0
	s.mean = 0.0
	s.s = 0.0
}

// Add feeds a new value into the Stats.
func (s *Stats) Add(val float64) {
	if s.count == 0 {
		s.min = val
		s.max = val
	} else {
		if val < s.min {
			s.min = val
		}
		if val > s.max {
			s.max = val
		}
	}
	s.count++
	old_mean := s.mean
	s.mean += (val - old_mean) / float64(s.count)
	s.s += (val - old_mean) * (val - s.mean)
}
