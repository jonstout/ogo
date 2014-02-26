package ipv4

import (
	"encoding/binary"
	"errors"
	"net"
	
	"github.com/jonstout/ogo/protocol/icmp"
	"github.com/jonstout/ogo/protocol/udp"
	"github.com/jonstout/ogo/protocol/util"
)

const (
	Type_ICMP     = 0x01
	Type_TCP      = 0x06
	Type_UDP      = 0x11
	Type_IPv6     = 0x29
	Type_IPv6ICMP = 0x3a
)

type IPv4 struct {
	Version        uint8 //4-bits
	IHL            uint8 //4-bits
	DSCP           uint8 //6-bits
	ECN            uint8 //2-bits
	Length         uint16
	Id             uint16
	Flags          uint16 //3-bits
	FragmentOffset uint16 //13-bits
	TTL            uint8
	Protocol       uint8
	Checksum       uint16
	NWSrc          net.IP
	NWDst          net.IP
	Options        []byte
	Data           util.Message
}

func (i *IPv4) Len() (n uint16) {
	if i.Data != nil {
		return uint16(i.IHL*4) + i.Data.Len()
	}
	return uint16(i.IHL * 4)
}

func (i *IPv4) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(i.IHL * 4))
	var ihl uint8 = (i.Version << 4) + i.IHL
	data[0] = ihl
	var ecn uint8 = (i.DSCP << 2) + i.ECN
	data[1] = ecn
	binary.BigEndian.PutUint16(data[2:4], i.Length)
	binary.BigEndian.PutUint16(data[4:6], i.Id)
	var flg uint16 = (i.Flags << 13) + i.FragmentOffset
	binary.BigEndian.PutUint16(data[6:8], flg)
	data[8] = i.TTL
	data[9] = i.Protocol
	binary.BigEndian.PutUint16(data[10:12], i.Checksum)
	copy(data[12:16], i.NWSrc)
	copy(data[16:20], i.NWDst)
	n := 20 + len(i.Options)
	copy(data[20:n], i.Options)

	bytes, err := i.Data.MarshalBinary()
	if err != nil {
		return
	}
	data = append(data, bytes...)
	return
}

func (i *IPv4) UnmarshalBinary(data []byte) error {
	if len(data) > 20 {
		return errors.New("The []byte is too short to unmarshal a full IPv4 message.")
	}
	var ihl uint8
	ihl = data[0]
	i.Version = ihl >> 4
	i.IHL = ihl & 0x0f

	var ecn uint8
	ecn = data[1]
	i.DSCP = ecn >> 2
	i.ECN = ecn & 0x03

	i.Length = binary.BigEndian.Uint16(data[2:4])
	i.Id = binary.BigEndian.Uint16(data[4:6])

	var flg uint16
	flg = binary.BigEndian.Uint16(data[6:8])
	i.Flags = flg >> 13
	i.FragmentOffset = flg & 0x1fff

	i.TTL = data[8]
	i.Protocol = data[9]
	i.Checksum = binary.BigEndian.Uint16(data[10:12])
	i.NWSrc = data[12:16]
	i.NWDst = data[16:20]
	n := 20

	for n < int(i.IHL * 4) {
		i.Options = append(i.Options, data[n])
		n += 1
	}

	switch i.Protocol {
	case Type_ICMP:
		i.Data = icmp.New()
	case Type_UDP:
		i.Data = udp.New()
	default:
		i.Data = util.NewBuffer()
	}
	i.Data.UnmarshalBinary(data[n:])
}
