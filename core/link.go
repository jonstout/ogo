package core

import "time"
import "net"

// Internal representation of a network link. Can be used to
// describe the state of the link. Each switch maintains its own
// set of links.
type Link struct {
	DPID net.HardwareAddr
	Port uint16
	Latency time.Duration
	Bandwidth int
}
