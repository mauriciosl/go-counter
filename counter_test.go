package counter

import (
	"runtime"
	"sync"

	. "gopkg.in/check.v1"
)

func (s *MySuite) TestCounterAdd(c *C) {
	x := NewCounter()
	x.Add(1)
	c.Assert(x.Value(), Equals, 1)
	x.Add(3)
	c.Assert(x.Value(), Equals, 4)
}

func (s *MySuite) TestSharedCounterAdd(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	x := NewCounter()
	w := &sync.WaitGroup{}
	w.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			x.Add(1)
			w.Done()
		}()
	}
	w.Wait()
	c.Assert(x.Value(), Equals, 100)
}

func (s *MySuite) TestDeltaCounter(c *C) {
	x := NewDeltaCounter()
	x.Add(2)
	c.Assert(x.Value(), Equals, 2)
	c.Assert(x.Delta(), Equals, 2)
	x.Reset(4)
	c.Assert(x.Value(), Equals, 4)
	c.Assert(x.Delta(), Equals, 0)
	x.Add(3)
	c.Assert(x.Value(), Equals, 7)
	c.Assert(x.Delta(), Equals, 3)
	x.Set(3)
	c.Assert(x.Value(), Equals, 3)
	c.Assert(x.Delta(), Equals, -1)
}

func (s *MySuite) TestObservableCounter(c *C) {
	c1 := NewObservableCounter(NewCounter())
	select {
	case <-c1.C:
		c.Fatal("Channel should start empty")
	default:
		c.Succeed()
	}
	c1.Add(1)
	select {
	case <-c1.C:
		c.Succeed()
	default:
		c.Fatal("Channel should have value")
	}
}
