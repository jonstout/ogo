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

func (s *OfpSwitchConfig) GetHeader() *OfpHeader {
	return &s.Header
}

func (s *OfpSwitchConfig) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpSwitchConfig) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = s.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &s.Flags); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &s.MissSendLen); err != nil {
		return
	}
	n += 2
	return
}
