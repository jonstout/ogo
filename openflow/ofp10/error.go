package ofp10

import (
	"io"
	"bytes"
	"encoding/binary"
)

// BEGIN: ofp10 - 5.4.4
// ofp_error_msg 1.0
type OfpErrorMsg struct {
	Header OfpHeader
	Code uint16
	Data []uint8
}

func (e *OfpErrorMsg) GetHeader() *OfpHeader {
	return &e.Header
}

func (e *OfpErrorMsg) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, e)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (e *OfpErrorMsg) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = e.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &e.Code); err != nil {
		return
	}
	n += 2
	e.Data = make([]uint8, buf.Len())
	m := buf.Len()
	if err = binary.Read(buf, binary.BigEndian, &e.Data); err != nil {
		return
	}
	n += m
	return
}

// ofp_error_type 1.0
const (
	OFPET_HELLO_FAILED = iota
	OFPET_BAD_REQUEST
	OFPET_BAD_ACTION
	OFPET_FLOW_MOD_FAILED
	OFPET_PORT_MOD_FAILED
	OFPET_QUEUE_OP_FAILED
)

// ofp_hello_failed_code 1.0
const (
	OFPHFC_INCOMPATIBLE = iota
	OFPHFC_EPERM
)

// ofp_bad_request_code 1.0
const (
	OFPBRC_BAD_VERSION = iota
	OFPBRC_BAD_TYPE
	OFPBRC_BAD_STAT
	OFPBRC_BAD_VENDOR

	OFPBRC_BAD_SUBTYPE
	OFPBRC_EPERM
	OFPBRC_BAD_LEN
	OFPBRC_BUFFER_EMPTY
	OFPBRC_BUFFER_UNKNOWN
)

// ofp_bad_action_code 1.0
const (
	OFPBAC_BAD_TYPE = iota
	OFPBAC_BAD_LEN
	OFPBAC_BAD_VENDOR
	OFPBAC_BAD_VENDOR_TYPE
	OFPBAC_BAD_OUT_PORT
	OFPBAC_BAD_ARGUMENT
	OFPBAC_EPERM
	OFPBAC_TOO_MANY
	OFPBAC_BAD_QUEUE
)

// ofp_flow_mod_failed_code 1.0
const (
	OFPFMFC_ALL_TABLES_FULL = iota
	OFPFMFC_OVERLAP
	OFPFMFC_EPERM
	OFPFMFC_BAD_EMERG_TIMEOUT
	OFPFMFC_BAD_COMMAND
	OFPFMFC_UNSUPPORTED
)

// ofp_port_mod_failed_code 1.0
const (
	OFPPMFC_BAD_PORT = iota
	OFPPMFC_BAD_HW_ADDR
)

// ofp_queue_op_failed_code 1.0
const (
	OFPQOFC_BAD_PORT = iota
	OFPQOFC_BAD_QUEUE
	OFPQOFC_EPERM
)
// END: ofp10 - 5.4.4
// END: ofp10 - 5.4
