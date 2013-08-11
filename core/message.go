package core

import (
	"io"
	//"bytes"
	"net"
	"time"
	"encoding/binary"
	//"github.com/jonstout/pacit"
)

type LinkDiscovery struct {
	src net.HardwareAddr
	nsec uint64 /* Number of nanoseconds elapsed since Jan 1, 1970. */
}

func NewListDiscovery(srcDPID net.HardwareAddr) (d *LinkDiscovery, err error) {
	d = new(LinkDiscovery)
	d.src = srcDPID
	d.nsec = uint64(time.Now().UnixNano())
	return
}

func (d *LinkDiscovery) Len() uint16 {
	return 16
}

func (d *LinkDiscovery) Read(b []byte) (n int, err error) {
	//binary.Write(buf, binary.BigEndian, d.src)
	b = append(b, []byte(d.src)...)
	n += 8
	//binary.Write(buf, binary.BigEndian, d.nsec)
	binary.BigEndian.PutUint64(b, d.nsec)
	//n, err = buf.Read(b)
	n += 8
	return n, io.EOF
}

func (d *LinkDiscovery) Write(b []byte) (n int, err error) {
	//d.src = make([]byte, 8)
	d.src = b[n:n+8]
	n += 8
	d.nsec = binary.BigEndian.Uint64(b[n:n+8])
	n += 8
	return
}
