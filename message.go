package ogo

import (
	"encoding/binary"
	"net"
	"time"
)

type LinkDiscovery struct {
	SrcDPID net.HardwareAddr
	Nsec    int64 /* Number of nanoseconds elapsed since Jan 1, 1970. */
	pad     []byte
}

func NewLinkDiscovery() *LinkDiscovery {
	d := new(LinkDiscovery)
	d.SrcDPID = make([]byte, 8)
	d.Nsec = time.Now().UnixNano()
	return d
}

func (d *LinkDiscovery) Len() uint16 {
	return 22
}

func (d *LinkDiscovery) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(d.Len()))
	
	next := 0
	copy(data[next:], d.SrcDPID)
	next += len(d.SrcDPID)
	binary.BigEndian.PutUint64(data[next:], uint64(d.Nsec))
	next += 8
	return
}

func (d *LinkDiscovery) UnmarshalBinary(data []byte) error {
	next := 0
	copy(d.SrcDPID, data[next:])
	next += len(d.SrcDPID)
	d.Nsec = int64(binary.BigEndian.Uint64(data[next:]))
	next += 8
	return nil
}
