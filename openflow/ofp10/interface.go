package ofp10

import (
	"net"
	)

type ConnectionUpReactor interface {
	ConnectionUp(dpid net.HardwareAddr)
}

type ConnectionDownReactor interface {
	ConnectionDown(dpid net.HardwareAddr, err error)
}

type SwitchFeaturesReactor interface {
	FeaturesReply(dpid net.HardwareAddr, features *SwitchFeatures)
}

type PacketInReactor interface {
	PacketIn(dpid net.HardwareAddr, packet *PacketIn)
}