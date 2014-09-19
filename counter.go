package counter

import "sync"

// Interface to counter objects
type Interface interface {
	Add(int) int
	Set(int)
	Reset(int) int
	Value() int
}

// Counter is an atomic counter
type Counter struct {
	sync.Mutex
	v int
}

// NewCounter creates a new counter with 0 value
func NewCounter(v int) *Counter {
	return &Counter{v: v}
}

// Add increments the counter by delta
func (c *Counter) Add(delta int) int {
	c.Lock()
	defer c.Unlock()
	c.v += delta
	return c.v
}

// Set sets the counter to a value
func (c *Counter) Set(v int) {
	c.Lock()
	defer c.Unlock()
	c.v = v
}

// Reset sets the counter to a value and return the old value
func (c *Counter) Reset(v int) int {
	c.Lock()
	defer c.Unlock()
	old := c.v
	c.v = v
	return old
}

// Value return the current counter value
func (c *Counter) Value() int {
	return c.v
}

// DeltaCounter is a Counter that accumulates a Delta
type DeltaCounter struct {
	value Interface
	delta Interface
}

// NewDeltaCounter creates a new counter that keeps track of the delta
func NewDeltaCounter(v int) *DeltaCounter {
	return &DeltaCounter{&Counter{v: v}, &Counter{}}
}

// Add increment the counter by a delta
func (d *DeltaCounter) Add(delta int) int {
	d.delta.Add(delta)
	return d.value.Add(delta)
}

// Set sets the counter to a value and clean the delta
func (d *DeltaCounter) Set(v int) {
	d.delta.Add(v - d.value.Value())
	d.value.Set(v)
}

// Value return the current counter value
func (d *DeltaCounter) Value() int {
	return d.value.Value()
}

// Delta return the counter changes since last reset
func (d *DeltaCounter) Delta() int {
	return d.delta.Value()
}

// Reset will empty the delta, set the counter to v and return the old value
func (d *DeltaCounter) Reset(v int) int {
	d.delta.Reset(0)
	d.value.Set(v)
	return v
}

// ObservableCounter creates a channel to notify counter changes
type ObservableCounter struct {
	Counter Interface
	C       chan bool
}

// NewObservableCounter creates an counter that notifies changes to a channel
func NewObservableCounter(counter Interface) *ObservableCounter {
	d := &ObservableCounter{counter, make(chan bool, 1)}
	return d
}

// Add and notify
func (d *ObservableCounter) Add(delta int) int {
	v := d.Counter.Add(delta)
	select {
	case d.C <- true:
	default:
	}
	return v
}

// Set and notify
func (d *ObservableCounter) Set(v int) {
	d.Counter.Set(v)
	select {
	case d.C <- true:
	default:
	}
}

// Reset and notify
func (d *ObservableCounter) Reset(v int) int {
	v = d.Counter.Reset(v)
	select {
	case d.C <- true:
	default:
	}
	return v
}

// Value return the current counter value
func (d *ObservableCounter) Value() int {
	return d.Counter.Value()
}
