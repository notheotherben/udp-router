package netmap

import (
	"testing"

	"github.com/spartan563/udp-router/internal/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestNetmap(t *testing.T) {
	nm := New(protocol.Address(1))

	assert.NotNil(t, nm)

	rt, ok := nm.Route(protocol.Address(2))
	assert.False(t, ok, "should not find non-existent routes")

	diff := nm.Update(Route{
		Source: protocol.Address(1),
		Dest:   protocol.Address(2),
		Port:   10002,
		Cost:   8,
	})

	assert.Len(t, diff, 1)

	rt, ok = nm.Route(protocol.Address(2))
	assert.True(t, ok, "should find existing routes")
	assert.Equal(t, 8, rt.EstCost)

	diff = nm.Update(Route{
		Source: protocol.Address(1),
		Dest:   protocol.Address(3),
		Port:   10003,
		Cost:   2,
	})

	assert.Len(t, diff, 1)

	diff = nm.Update(Route{
		Source: protocol.Address(3),
		Dest:   protocol.Address(2),
		Port:   10002,
		Cost:   1,
	})

	assert.Len(t, diff, 1)

	rt, ok = nm.Route(protocol.Address(2))
	assert.True(t, ok, "should find shortest routes")
	assert.Equal(t, 3, rt.EstCost)
	assert.Equal(t, protocol.Address(3), rt.Via.Dest)

	ns := nm.Neighbours()
	assert.Len(t, ns, 2, "address 1 should have 2 known neighbours")
}
