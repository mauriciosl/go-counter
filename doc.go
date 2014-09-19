/*
	Package counter implements some types of counters and a rate limit tracker

	Usage:

	A basic atomic(thread safe) counter:

		c := NewCounter(0)
		c.Add(1) // 1
		c.Value() // 1
		c.Set(5)
		c.Value() // 5
		c.Reset(1) // 5
		c.Value() // 1

	A DeltaCounter keeps track of the changes to the counter since last Reset

		d := NewDeltaCounter(0)
		d.Add(5) // 5
		d.Delta() // 5
		d.Reset(1) // 5
		d.Delta() // 0
		d.Value() // 1
		d.Add(-3)
		d.Delta() // -3
		d.Value() // -2

	A RateLimit tracks the rate of use of some resource.
	To limit the rate of use to let's say, 5 per minute.

		r := NewRateLimit(5, time.Minute)
		if r.Hit() {
			// rate limit exceeded
		}
		// normal code

	You will achieve grater precision with greater time deltas, so use
	60 / Minute instead of 1 / Second

		// Worse
		r := NewRateLimite(1, time.Second)
		// Better
		r := NewRateLimit(60, time.Minute)

*/
package counter
