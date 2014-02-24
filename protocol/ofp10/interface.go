package ofp10

import (
	"net"
	
	"github.com/jonstout/ogo/protocol/ofpxx"
)

type ConnectionUpReactor interface {
	ConnectionUp(dpid net.HardwareAddr)
}

type ConnectionDownReactor interface {
	ConnectionDown(dpid net.HardwareAddr, err error)
}

type HelloReactor interface {
	Hello(hello ofpxx.Hello)
}

type ErrorReactor interface {
	Error(dpid net.HardwareAddr, err *ErrorMsg)
}

type VendorReactor interface {
	VendorHeader(dpid net.HardwareAddr, v *VendorHeader)
}

type SwitchFeaturesRequestReactor interface {
	FeaturesRequest(features *ofpxx.Header)
}

type SwitchFeaturesReplyReactor interface {
	FeaturesReply(dpid net.HardwareAddr, features *SwitchFeatures)
}

type ConfigRequestReactor interface {
	GetConfigRequest(config *ofpxx.Header)
}

type ConfigReplyReactor interface {
	GetConfigReply(dpid net.HardwareAddr, config *SwitchConfig)
}

type SetConfigReactor interface {
	SetConfigRequest(config *SwitchConfig)
}

type EchoReplyReactor interface {
	EchoReply(dpid net.HardwareAddr)
}

type EchoRequestReactor interface {
	EchoRequest(dpid net.HardwareAddr)
}

type PacketInReactor interface {
	PacketIn(dpid net.HardwareAddr, packet *PacketIn)
}

type PacketOutReactor interface {
	PacketOut(dpid net.HardwareAddr, packet *PacketOut)
}
