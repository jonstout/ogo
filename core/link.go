package core

import "time"
import "net"

/*
In general. Each switch should keep track of any switches
connected to itself.
*/
type Link struct {
	DPID net.HardwareAddr
	Port uint16
	Latency time.Duration
	Bandwidth int
}
