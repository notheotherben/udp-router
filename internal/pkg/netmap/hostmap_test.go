package netmap

import (
	"testing"

	"github.com/soniakeys/graph"
	"github.com/spartan563/udp-router/internal/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestHostmap(t *testing.T) {
	hm := newHostmap()

	assert.NotNil(t, hm)

	assert.Equal(t, graph.NI(0), hm.Get(protocol.Address(1)), "should start index at 0")
	assert.Equal(t, hm.Get(protocol.Address(1)), hm.Get(protocol.Address(1)), "should return consistent indices")

	assert.Equal(t, graph.NI(1), hm.Get(protocol.Address(2)), "should increment indices")

	addr, ok := hm.Reverse(graph.NI(0))
	assert.True(t, ok, "should find reversed addresses")
	assert.Equal(t, protocol.Address(1), addr, "should reverse indices correctly")
}
