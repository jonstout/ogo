package ofp10

import (
	"bytes"
	"encoding/binary"
	"io"
)

func NewConfigRequest() *OfpHeader {
	h := NewHeader()
	h.Type = OFPT_GET_CONFIG_REQUEST
	return h
}

// ofp_config_flags 1.0
const (
	OFPC_FRAG_NORMAL = 0
	OFPC_FRAG_DROP = 1
	OFPC_FRAG_REASM = 2
	OFPC_FRAG_MASK = 3
)

// ofp_switch_config 1.0
type OfpSwitchConfig struct {
	Header OfpHeader
	Flags uint16 // OFPC_* flags
	MissSendLen uint16
}

func NewSetConfig() *OfpSwitchConfig {
	h := NewHeader()
	h.Type = OFPT_SET_CONFIG

	c := new(OfpSwitchConfig)
	c.Header = *h
	c.Flags = 0
	c.MissSendLen = 0
	return c
}

func (c *OfpSwitchConfig) Len() (n uint16) {
	return 12
}

func (c *OfpSwitchConfig) GetHeader() *OfpHeader {
	return &c.Header
}

func (c *OfpSwitchConfig) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, c)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (c *OfpSwitchConfig) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = c.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &c.Flags); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &c.MissSendLen); err != nil {
		return
	}
	n += 2
	return
}
