package ofp10

import "net"

type ConnectionUpReactor interface {
	ConnectionUp(dpid net.HardwareAddr, features SwitchFeatures)
}

type ConnectionDownReactor interface {
	ConnectionDown(dpid net.HardwareAddr)
}