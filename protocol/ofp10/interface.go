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
	Hello(hello *ofpxx.Header)
}

type ErrorReactor interface {
	Error(dpid net.HardwareAddr, err *ErrorMsg)
}

type EchoRequestReactor interface {
	EchoRequest(dpid net.HardwareAddr)
}

type EchoReplyReactor interface {
	EchoReply(dpid net.HardwareAddr)
}

type VendorReactor interface {
	VendorHeader(dpid net.HardwareAddr, v *VendorHeader)
}

type FeaturesRequestReactor interface {
	FeaturesRequest(features *ofpxx.Header)
}

type FeaturesReplyReactor interface {
	FeaturesReply(dpid net.HardwareAddr, features *SwitchFeatures)
}

type GetConfigRequestReactor interface {
	GetConfigRequest(config *ofpxx.Header)
}

type GetConfigReplyReactor interface {
	GetConfigReply(dpid net.HardwareAddr, config *SwitchConfig)
}

type SetConfigReactor interface {
	SetConfig(config *SwitchConfig)
}

type PacketInReactor interface {
	PacketIn(dpid net.HardwareAddr, packet *PacketIn)
}

type FlowRemovedReactor interface {
	FlowRemoved(dpid net.HardwareAddr, flow *FlowRemoved)
}

type PortStatusReactor interface {
	PortStatus(dpid net.HardwareAddr, status *PortStatus)
}

type PacketOutReactor interface {
	PacketOut(packet *PacketOut)
}

type FlowModReactor interface {
	FlowMod(flowMod *FlowMod)
}

type PortModReactor interface {
	PortMod(portMod *PortMod)
}

type StatsRequestReactor interface {
	StatsRequest(req *StatsRequest)
}

type StatsReplyReactor interface {
	StatsReply(dpid net.HardwareAddr, rep *StatsReply)
}

type BarrierRequestReactor interface {
	BarrierRequest(req *ofpxx.Header)
}

type BarrierReplyReactor interface {
	BarrierReply(dpid net.HardwareAddr, msg *ofpxx.Header)
}
