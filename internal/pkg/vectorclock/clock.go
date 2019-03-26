package vectorclock

import (
	"math/rand"
)

type Clock struct {
	id int
	tick int
}

func New() Clock {
	return Clock{
		id: rand.Int(),
		tick: 0,
	}
}

func (c Clock) Next() Clock {
	return Clock{
		id: c.id,
		tick: c.tick + 1,
	}
}

func (c Clock) IsBefore(b Clock) bool {
	return c.id == b.id && c.tick < b.tick
}

func (c Clock) IsEqual(b Clock) bool {
	return c.id == b.id && c.tick == b.tick
}