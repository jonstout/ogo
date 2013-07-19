package core

import (
	"io"
	"bytes"
	"net"
	"encoding/binary"
)

type LinkDiscovery struct {
	src net.HardwareAddr
	nsec int64 /* Number of nanoseconds elapsed since Jan 1, 1970. */
}

func NewListDiscovery(s string) (d LinkDiscovery, err error) {
	d := new(LinkDiscovery)
	mac := new(net.HardwareAddr)

	if mac, err = net.ParseMAC(s); err != nil {
		return nil, err
	} else {
		d.src = mac
		d.nsec = time.Now().UnixNano()
		return
	}
}

func (d *LinkDiscovery) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, d.src)
	binary.Write(buf, binary.BigEndian, d.nsec)
	n, err = buf.Read(b)
	return n, io.EOF
}

func (d *LinkDiscovery) Write(b []byte) (n int, err error) {
	d.src = make([]byte, 8)
	d.src = b[:8]
	n += 8
	d.nsec = binary.BigEndian.int64(b[8:16])
	n += 8
	return
}
