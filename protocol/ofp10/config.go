package ofp10

import (
	"encoding/binary"

	"github.com/jonstout/ogo/protocol/ofpxx"
)

func NewConfigRequest() *ofpxx.Header {
	h := ofpxx.NewOfp10Header()
	h.Type = Type_GetConfigRequest
	return &h
}

// ofp_config_flags 1.0
const (
	C_FRAG_NORMAL = 0
	C_FRAG_DROP   = 1
	C_FRAG_REASM  = 2
	C_FRAG_MASK   = 3
)

// ofp_switch_config 1.0
type SwitchConfig struct {
	ofpxx.Header
	Flags       uint16 // OFPC_* flags
	MissSendLen uint16
}

func NewSetConfig() *SwitchConfig {
	c := new(SwitchConfig)
	c.Header = ofpxx.NewOfp10Header()
	c.Header.Type = Type_SetConfig
	c.Flags = 0
	c.MissSendLen = 0
	return c
}

func (c *SwitchConfig) Len() (n uint16) {
	n = c.Header.Len()
	n += 4
	return
}

func (c *SwitchConfig) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(c.Len()))
	bytes := make([]byte, 0)
	next := 0

	c.Header.Length = c.Len()
	bytes, err = c.Header.MarshalBinary()
	copy(data[next:], bytes)
	next += len(bytes)
	binary.BigEndian.PutUint16(data[next:], c.Flags)
	next += 2
	binary.BigEndian.PutUint16(data[next:], c.MissSendLen)
	next += 2
	return
}

func (c *SwitchConfig) UnmarshalBinary(data []byte) error {
	var err error
	next := 0

	err = c.Header.UnmarshalBinary(data[next:])
	next += int(c.Header.Len())
	c.Flags = binary.BigEndian.Uint16(data[next:])
	next += 2
	c.MissSendLen = binary.BigEndian.Uint16(data[next:])
	next += 2
	return err
}
