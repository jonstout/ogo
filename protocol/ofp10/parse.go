package ofp10

import (
	"errors"
	
	"github.com/jonstout/ogo/protocol/ofpxx"
	"github.com/jonstout/ogo/protocol/util"
)

func Parse(b []byte) (message util.Message, err error) {
	switch b[1] {
	case Type_Hello:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	case Type_Error:
		message = new(ErrorMsg)
		message.UnmarshalBinary(b)
	case Type_EchoRequest:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	case Type_EchoReply:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	case Type_Vendor:
		message = new(VendorHeader)
		message.UnmarshalBinary(b)
	 case Type_FeaturesRequest:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	 case Type_FeaturesReply:
		message = new(SwitchFeatures)
		message.UnmarshalBinary(b)
	case Type_GetConfigRequest:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	case Type_GetConfigReply:
		message = new(SwitchConfig)
		message.UnmarshalBinary(b)
	case Type_SetConfig:
		message = NewSetConfig()
		message.UnmarshalBinary(b)
	case Type_PacketIn:
		message = new(PacketIn)
		message.UnmarshalBinary(b)
	case Type_FlowRemoved:
		message = NewFlowRemoved()
		message.UnmarshalBinary(b)
	case Type_PortStatus:
		message = new(PortStatus)
		message.UnmarshalBinary(b)
	case Type_PacketOut:
		break
	case Type_FlowMod:
		message = NewFlowMod()
		message.UnmarshalBinary(b)
	case Type_PortMod:
		break
	case Type_StatsRequest:
		message = new(StatsRequest)
		message.UnmarshalBinary(b)
	case Type_StatsReply:
		message = new(StatsReply)
		message.UnmarshalBinary(b)
	 case Type_BarrierRequest:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	 case Type_BarrierReply:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	case Type_QueueGetConfigRequest:
		break
	case Type_QueueGetConfigReply:
		break
	default:
		err = errors.New("An unknown v1.0 packet type was received. Parse function will discard data.")
	}
	return
}
