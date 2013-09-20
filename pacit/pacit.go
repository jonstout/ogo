package pacit

import (
	//"fmt"
)

type Packet struct {
	Preamble [7]uint8
	Delimiter uint8
	HWDst [6]uint8
	HWSrc [6]uint8
	VLANHeader VLAN
	Ethertype uint16
	Payload []uint8
}
