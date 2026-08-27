package ts2

import "sort"

// gnssSolverPorts are the GNSS NAV grouping ports whose position is produced by
// the solver rather than computed on the device. Decode dispatches on this set,
// and it is exported through IsGnssSolverPort/GnssSolverPorts so consumers do
// not have to restate the port numbers.
var gnssSolverPorts = map[uint8]struct{}{
	192: {},
	193: {},
	194: {},
	195: {},
	199: {},
	210: {},
	211: {},
}

// IsGnssSolverPort reports whether a GNSS fix arriving on this port is
// solver-derived. A GNSS fix on any other port is computed on the device.
func IsGnssSolverPort(port uint8) bool {
	_, ok := gnssSolverPorts[port]
	return ok
}

// GnssSolverPorts returns the solver-backed GNSS ports in ascending order.
func GnssSolverPorts() []uint8 {
	ports := make([]uint8, 0, len(gnssSolverPorts))
	for port := range gnssSolverPorts {
		ports = append(ports, port)
	}
	sort.Slice(ports, func(i, j int) bool { return ports[i] < ports[j] })
	return ports
}
