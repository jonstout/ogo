package ofp10

import (
	"errors"
)

func Parse(b []byte) (message *ofp.Message, err error) {
	switch buf[1] {
	case ofp10.T_PACKET_IN:
		message = new(ofp10.PacketIn)
		message.Write(buf)
	case ofp10.T_HELLO:
		message = new(ofp10.Header)
		message.Write(buf)
	case ofp10.T_ECHO_REPLY:
		message = new(ofp10.Header)
		message.Write(buf)
	case ofp10.T_ECHO_REQUEST:
		message = new(ofp10.Header)
		message.Write(buf)
	case ofp10.T_ERROR:
		message = new(ofp10.ErrorMsg)
		message.Write(buf)
	case ofp10.T_VENDOR:
		message = new(ofp10.VendorHeader)
		message.Write(buf)
	case ofp10.T_FEATURES_REPLY:
		message = new(ofp10.SwitchFeatures)
		message.Write(buf)
	case ofp10.T_GET_CONFIG_REPLY:
		message = new(ofp10.SwitchConfig)
		message.Write(buf)
	case ofp10.T_FLOW_REMOVED:
		message = new(ofp10.FlowRemoved)
		message.Write(buf)
	case ofp10.T_PORT_STATUS:
		message = new(ofp10.PortStatus)
		message.Write(buf)
	case ofp10.T_STATS_REPLY:
		message = new(ofp10.StatsReply)
		message.Write(buf)
	case ofp10.T_BARRIER_REPLY:
		message = new(ofp10.Header)
		message.Write(buf)
	default:
		err = errors.New("An unknown v1.0 packet type was received. Parse function will discard data.")
	}
	return
}
