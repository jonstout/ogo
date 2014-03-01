package ofp10

import (
	"errors"
	
	"github.com/jonstout/ogo/protocol/ofpxx"
	"github.com/jonstout/ogo/protocol/util"
)

func Parse(b []byte) (message util.Message, err error) {
	switch b[1] {
	/*case Type_Packet_In:
		message = new(PacketIn)
		message.Write(b)*/
	case Type_Hello:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	case Type_Echo_Reply:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	case Type_Echo_Request:
		message = new(ofpxx.Header)
		message.UnmarshalBinary(b)
	/*case Type_Error:
		message = new(ErrorMsg)
		message.Write(b)
	case Type_Vendor:
		message = new(VendorHeader)
		message.Write(b)
	 case Type_Features_Reply:
		message = new(SwitchFeatures)
		message.Write(b)
	case Type_Get_Config_Reply:
		message = new(SwitchConfig)
		message.Write(b)
	case Type_Flow_Removed:
		message = new(FlowRemoved)
		message.Write(b)
	case Type_Port_Status:
		message = new(PortStatus)
		message.Write(b)
	case Type_Stats_Reply:
		message = new(StatsReply)
		message.Write(b)
	case Type_Stats_Request:
		message = new(StatsRequest)
		message.Write(b)
	 case Type_Barrier_Reply:
		message = new(ofpxx.Header)
		message.UnmarshelBinary(b)
	 case Type_Barrier_Request:
		message = new(ofpxx.Header)
		message.UnmarshelBinary(b)*/
	default:
		err = errors.New("An unknown v1.0 packet type was received. Parse function will discard data.")
	}
	return
}
