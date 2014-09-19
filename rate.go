package counter

import "time"

// Rate keeps track of hits
type Rate struct {
	Hits   Interface
	DeltaT time.Duration
}

// HitsPerDuration return the rate of hits in the duration d
func (r Rate) HitsPerDuration(d time.Duration) float64 {
	if d < r.DeltaT {
		d = r.DeltaT
	}
	f := (float64(d) / float64(r.DeltaT))
	if f == 0 {
		return float64(r.Hits.Value())
	}
	return float64(r.Hits.Value()) / f
}

// RateLimit tracks the rate of hits in time
type RateLimit struct {
	R     Rate
	Limit int
	T0    time.Time
}

// NewRateLimit creates a new rate limit tracker
func NewRateLimit(limit int, deltat time.Duration) *RateLimit {
	return &RateLimit{Rate{NewCounter(0), deltat}, limit, time.Now()}
}

// Exceeded return if the rate limit is exceeded
func (r RateLimit) Exceeded() bool {
	return r.R.HitsPerDuration(time.Since(r.T0)) > float64(r.Limit)
}

// Hit increments the hit count and return if the rate is exceeded
func (r *RateLimit) Hit() bool {
	return r.timeHit(time.Now())
}

func (r *RateLimit) timeHit(t time.Time) bool {
	r.R.Hits.Add(1)
	delta := t.Sub(r.T0)
	if delta > r.R.DeltaT {
		r.R.Hits.Set(1)
		r.T0 = t
	}
	return r.Exceeded()
}

// HitsPerDuration return the current hit rate
func (r RateLimit) HitsPerDuration() float64 {
	return r.R.HitsPerDuration(time.Since(r.T0))
}
