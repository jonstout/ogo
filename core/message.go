package core

import (
	"io"
	"bytes"
	"net"
	"time"
	"encoding/binary"
	"github.com/jonstout/pacit"
)

type LinkDiscovery struct {
	eth pacit.Ethernet
	src net.HardwareAddr
	nsec uint64 /* Number of nanoseconds elapsed since Jan 1, 1970. */
}

func NewListDiscovery(s string) (d *LinkDiscovery, err error) {
	d = new(LinkDiscovery)
	d.eth = *new(pacit.Ethernet)
	mac := *new(net.HardwareAddr)

	if mac, err = net.ParseMAC(s); err != nil {
		return nil, err
	}
	d.src = mac
	d.nsec = uint64(time.Now().UnixNano())
	return
}

func (d *LinkDiscovery) Len() uint16 {
	return d.eth.Len() + 16
}

func (d *LinkDiscovery) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	if n, err = d.eth.Read(b); err != nil {
		return
	}
	//binary.Write(buf, binary.BigEndian, d.src)
	b = append(b, []byte(d.src)...)
	binary.Write(buf, binary.BigEndian, d.nsec)
	n, err = buf.Read(b)
	return n, io.EOF
}

func (d *LinkDiscovery) Write(b []byte) (n int, err error) {
	n, err = d.eth.Write(b)
	d.src = make([]byte, 8)
	d.src = b[n:n+8]
	n += 8
	d.nsec = binary.BigEndian.Uint64(b[n:n+8])
	n += 8
	return
}
