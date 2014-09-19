package counter

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *MySuite) TestRatePerDuration(c *C) {
	r := Rate{100, time.Second}
	c.Assert(r.HitsPerDuration(0*time.Second), Equals, float64(100))
	c.Assert(r.HitsPerDuration(500*time.Millisecond), Equals, float64(200))
	c.Assert(r.HitsPerDuration(2*time.Second), Equals, float64(50))
	c.Assert(r.HitsPerDuration(3*time.Second), Equals, float64(100)/3)
}

func (s *MySuite) TestRateLimitExceeded(c *C) {
	r := Rate{0, time.Second}
	rl := RateLimit{r, 100, time.Now().Add(-1 * time.Second)}
	c.Assert(rl.Exceeded(), Equals, false)
	rl.R.Hits = 100
	c.Assert(rl.Exceeded(), Equals, false)
	rl.R.Hits = 101
	c.Assert(rl.Exceeded(), Equals, true)
}

func (s *MySuite) TestRateLimitHit(c *C) {
	t := time.Now().Add(-2 * time.Second)
	ra := RateLimit{Rate{0, time.Second}, 100, t}
	t = time.Now()
	ra.timeHit(t)
	c.Assert(ra.R.Hits, Equals, 1)
	c.Assert(ra.T0, Equals, t)
}

func (s *MySuite) TestRateTracking(c *C) {
	t := time.Now().Add(-1 * time.Minute)
	ra := RateLimit{Rate{0, time.Minute}, 5, t}
	c.Assert(ra.timeHit(t), Equals, false)
	c.Assert(ra.timeHit(t.Add(10*time.Second)), Equals, false)
	c.Assert(ra.timeHit(t.Add(20*time.Second)), Equals, false)
	c.Assert(ra.timeHit(t.Add(30*time.Second)), Equals, false)
	c.Assert(ra.timeHit(t.Add(40*time.Second)), Equals, false)
	c.Assert(ra.timeHit(t.Add(50*time.Second)), Equals, true)
	c.Assert(ra.timeHit(t.Add(70*time.Second)), Equals, false)
}
