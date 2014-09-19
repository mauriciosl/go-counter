package counter

import "time"

type fibTimer struct {
	i          int
	max        int
	resolution time.Duration
	evis       bool
	C          chan time.Time
}

func newFibTimer(max int, resolution time.Duration) *fibTimer {
	f := &fibTimer{1, max, resolution, false, make(chan time.Time)}
	go f.Run()
	return f
}

func (f *fibTimer) Reset() {
	f.i = 1
}

func (f *fibTimer) Stop() {
	f.evis = false
}

func (f *fibTimer) Run() {
	f.evis = true
	for f.evis {
		t := <-time.After(time.Duration(fib(f.i+1)) * f.resolution)
		f.C <- t
		f.i = min(f.max, f.i+1)
	}
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
