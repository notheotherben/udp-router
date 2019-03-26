package netmap

import (
	"github.com/soniakeys/graph"
	"github.com/spartan563/udp-router/internal/pkg/protocol"
)

// hostmap provides a consistent store of host addresses to
// graph node ids.
type hostmap struct {
	mapping map[protocol.Address]graph.NI
	reverse map[graph.NI]protocol.Address
	nextID  int
}

func newHostmap() *hostmap {
	return &hostmap{
		mapping: map[protocol.Address]graph.NI{},
		reverse: map[graph.NI]protocol.Address{},
		nextID:  0,
	}
}

func (m *hostmap) Get(host protocol.Address) graph.NI {
	if id, ok := m.mapping[host]; ok {
		return id
	}

	m.mapping[host] = graph.NI(m.nextID)
	m.reverse[graph.NI(m.nextID)] = host
	m.nextID++

	return m.mapping[host]
}

func (m *hostmap) Reverse(id graph.NI) (protocol.Address, bool) {
	addr, ok := m.reverse[id]
	return addr, ok
}

func (m *hostmap) Size() int {
	return len(m.mapping)
}
