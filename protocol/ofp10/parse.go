package ofp10

import (
	"errors"
	
	"github.com/jonstout/ogo/protocol/util"
)

func Parse(b []byte) (message util.Message, err error) {
	switch b[1] {
	case T_PACKET_IN:
		message = new(PacketIn)
		message.Write(b)
	case T_HELLO:
		message = new(Header)
		message.Write(b)
	case T_ECHO_REPLY:
		message = new(Header)
		message.Write(b)
	case T_ECHO_REQUEST:
		message = new(Header)
		message.Write(b)
	case T_ERROR:
		message = new(ErrorMsg)
		message.Write(b)
	case T_VENDOR:
		message = new(VendorHeader)
		message.Write(b)
	case T_FEATURES_REPLY:
		message = new(SwitchFeatures)
		message.Write(b)
	case T_GET_CONFIG_REPLY:
		message = new(SwitchConfig)
		message.Write(b)
	case T_FLOW_REMOVED:
		message = new(FlowRemoved)
		message.Write(b)
	case T_PORT_STATUS:
		message = new(PortStatus)
		message.Write(b)
	case T_STATS_REPLY:
		message = new(StatsReply)
		message.Write(b)
	case T_BARRIER_REPLY:
		message = new(Header)
		message.Write(b)
	default:
		err = errors.New("An unknown v1.0 packet type was received. Parse function will discard data.")
	}
	return
}
