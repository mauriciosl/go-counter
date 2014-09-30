package counter

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *MySuite) TestFibTimerMax(c *C) {
	t := NewFibTimer(4, time.Millisecond)
	for i := 1; i <= 3; i++ {
		_ = <-t.Timer()
		c.Assert(t.i, Equals, i+1)
	}
	_ = <-t.Timer()
	c.Assert(t.i, Equals, 4)
	t.Reset()
	c.Assert(t.i, Equals, 1)
	_ = <-t.Timer()
	c.Assert(t.i, Equals, 2)
}
