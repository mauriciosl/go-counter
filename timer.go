package counter

import "time"

// FibTimer will sleep in a fibonnaci sequence time at max sequence index until reset is called or stop
// f := NewFibTimer(5, time.Second)
// start := time.Now()
// for t := range f.Timer() {
//   1, 2, 3, 5, 8, 8, 8 seconds...
// }
// After a Reset, the timer go back to 1
type FibTimer struct {
	i          int
	max        int
	resolution time.Duration
	evis       bool
	c          chan time.Time
}

func NewFibTimer(max int, resolution time.Duration) *FibTimer {
	f := &FibTimer{1, max, resolution, false, make(chan time.Time)}
	go f.Run()
	return f
}

func (f *FibTimer) Reset() {
	f.i = 1
}

func (f *FibTimer) Stop() {
	f.evis = false
}

func (f *FibTimer) Timer() <-chan time.Time {
	return f.c
}

func (f *FibTimer) Run() {
	f.evis = true
	for f.evis {
		t := <-time.After(time.Duration(fib(f.i+1)) * f.resolution)
		f.c <- t
		f.i = min(f.max, f.i+1)
	}
	close(f.c)
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func fib(i int) int {
	if i <= 1 {
		return 1
	}
	return i + fib(i-1)
}
