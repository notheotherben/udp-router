package netmap

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/soniakeys/graph"
	"github.com/spartan563/udp-router/internal/pkg/protocol"
)

type Route struct {
	Source protocol.Address
	Dest   protocol.Address
	Cost   int
	Port   int
}

type ViaRoute struct {
	Source  protocol.Address
	Dest    protocol.Address
	Via     Route
	EstCost int
}

func (r *ViaRoute) Route() Route {
	return Route{
		Source: r.Source,
		Dest:   r.Dest,
		Cost:   r.EstCost,
		Port:   r.Via.Port,
	}
}

type Map struct {
	root       protocol.Address
	routes     []Route
	routeTable map[protocol.Address]ViaRoute
	hostmap    *hostmap
	m          sync.RWMutex
}

func New(root protocol.Address) *Map {
	return &Map{
		root:       root,
		routes:     []Route{},
		routeTable: map[protocol.Address]ViaRoute{},
		hostmap:    newHostmap(),
	}
}

func (m *Map) Neighbours() []Route {
	out := []Route{}

	for _, route := range m.routes {
		if route.Source == m.root {
			out = append(out, route)
		}
	}

	return out
}

func (m *Map) Route(target protocol.Address) (ViaRoute, bool) {
	m.m.RLock()
	defer m.m.RUnlock()

	route, ok := m.routeTable[target]
	return route, ok
}

func (m *Map) Update(route Route) map[protocol.Address]ViaRoute {
	m.m.Lock()
	defer m.m.Unlock()

	for i, r := range m.routes {
		if r.Source == route.Source && r.Dest == route.Dest {
			m.routes[i] = route
			break
		}
	}

	m.hostmap.Get(route.Source)
	m.hostmap.Get(route.Dest)

	m.routes = append(m.routes, route)

	g := graph.LabeledDirected{
		LabeledAdjacencyList: m.buildLAL(),
	}

	w := func(label graph.LI) float64 {
		return float64(m.routes[int(label)].Cost)
	}

	start := m.hostmap.Get(m.root)
	f, _, dist, end := g.BellmanFord(w, start)

	if end >= 0 {
		// Invalid routing table - negative cycle from this node
		// Keep the existing routing table intact until this can be
		// resolved.
		return map[protocol.Address]ViaRoute{}
	}

	rt := map[protocol.Address]ViaRoute{}
	p := make([]graph.NI, f.MaxLen)

	for n, e := range f.Paths {
		if e.Len < 2 {
			// Ignore entries which are unreachable
			continue
		}

		target, ok := m.hostmap.Reverse(graph.NI(n))
		if !ok {
			// We should probably panic here...
			panic("could not find host address for graph node")
		}

		nextStep := f.PathTo(graph.NI(n), p)[1]

		nextHopAddress, ok := m.hostmap.Reverse(nextStep)
		if !ok {
			// We should probably panic here...
			panic("could not find host address for graph node")
		}

		route, ok := m.findRoute(m.root, nextHopAddress)
		if !ok {
			// We don't have a route which allows this traversal
			continue
		}

		logrus.WithFields(logrus.Fields{
			"src":     m.root,
			"dest":    target,
			"via":     route,
			"estCost": dist[n],
		}).Info("route table updated")

		rt[target] = ViaRoute{
			Source:  m.root,
			Dest:    target,
			Via:     route,
			EstCost: int(dist[n]),
		}
	}

	diff := m.diffTable(rt)

	m.routeTable = rt

	return diff
}

func (m *Map) diffTable(newTable map[protocol.Address]ViaRoute) map[protocol.Address]ViaRoute {
	diff := map[protocol.Address]ViaRoute{}

	for addr, route := range newTable {
		oldRoute, ok := m.routeTable[addr]
		if !ok {
			diff[addr] = route
			continue
		}

		if route.Via.Dest != oldRoute.Via.Dest || route.EstCost != oldRoute.EstCost {
			diff[addr] = route
			continue
		}
	}

	return diff
}

func (m *Map) findRoute(src, dest protocol.Address) (Route, bool) {
	for _, route := range m.routes {
		if route.Source == src && route.Dest == dest {
			return route, true
		}
	}

	return Route{
		Source: src,
		Dest:   dest,
		Cost:   -1,
		Port:   -1,
	}, false
}

func (m *Map) buildLAL() graph.LabeledAdjacencyList {
	// Keeps track of the route indices which originate at a given source
	srcMap := map[protocol.Address][]int{}

	for i, route := range m.routes {
		sm := srcMap[route.Source]
		if sm == nil {
			sm = []int{}
		}

		sm = append(sm, i)

		srcMap[route.Source] = sm
	}

	lal := make(graph.LabeledAdjacencyList, m.hostmap.Size())
	for src, routes := range srcMap {
		al := []graph.Half{}

		for _, rix := range routes {
			route := m.routes[rix]

			al = append(al, graph.Half{
				To:    m.hostmap.Get(route.Dest),
				Label: graph.LI(rix),
			})
		}

		lal[int(m.hostmap.Get(src))] = al
	}

	return lal
}
