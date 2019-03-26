package vectorclock_test

import (
	"testing"

	"github.com/spartan563/udp-router/internal/pkg/vectorclock"
	"github.com/stretchr/testify/assert"
)

func TestClock(t *testing.T) {
	a := vectorclock.New()
	b := vectorclock.New()

	assert.False(t, a.IsEqual(b), "incomparable clocks should report false equality")
	assert.False(t, a.IsBefore(b), "incomparable clocks should report false precedence")

	c := a.Next()
	assert.False(t, c.IsBefore(a), "newer clocks should not report that they precede another")
	assert.True(t, a.IsBefore(c), "older clocks should report that they precede another equatable one")
	assert.False(t, b.IsBefore(c), "inequitable clocks should not report precedence")
}
