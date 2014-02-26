package ofp10

import (
	"bytes"
	"encoding/binary"
	"io"
	
	"github.com/jonstout/ogo/protocol/ofpxx"
)

// BEGIN: ofp10 - 5.4.4
// ofp_error_msg 1.0
type ErrorMsg struct {
	Header ofpxx.Header
	Code   uint16
	Data   []uint8
}

func (e *ErrorMsg) GetHeader() *ofpxx.Header {
	return &e.Header
}

func (e *ErrorMsg) Len() (n uint16) {
	n = e.Header.Len()
	n += 2
	n += uint16(len(e.Data))
	return
}

func (e *ErrorMsg) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, e)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (e *ErrorMsg) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = e.Header.UnmarshelBinary(buf.Next(8))

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
	ET_HELLO_FAILED = iota
	ET_BAD_REQUEST
	ET_BAD_ACTION
	ET_FLOW_MOD_FAILED
	ET_PORT_MOD_FAILED
	ET_QUEUE_OP_FAILED
)

// ofp_hello_failed_code 1.0
const (
	HFC_INCOMPATIBLE = iota
	HFC_EPERM
)

// ofp_bad_request_code 1.0
const (
	BRC_BAD_VERSION = iota
	BRC_BAD_TYPE
	BRC_BAD_STAT
	BRC_BAD_VENDOR

	BRC_BAD_SUBTYPE
	BRC_EPERM
	BRC_BAD_LEN
	BRC_BUFFER_EMPTY
	BRC_BUFFER_UNKNOWN
)

// ofp_bad_action_code 1.0
const (
	BAC_BAD_TYPE = iota
	BAC_BAD_LEN
	BAC_BAD_VENDOR
	BAC_BAD_VENDOR_TYPE
	BAC_BAD_OUT_PORT
	BAC_BAD_ARGUMENT
	BAC_EPERM
	BAC_TOO_MANY
	BAC_BAD_QUEUE
)

// ofp_flow_mod_failed_code 1.0
const (
	FMFC_ALL_TABLES_FULL = iota
	FMFC_OVERLAP
	FMFC_EPERM
	FMFC_BAD_EMERG_TIMEOUT
	FMFC_BAD_COMMAND
	FMFC_UNSUPPORTED
)

// ofp_port_mod_failed_code 1.0
const (
	PMFC_BAD_PORT = iota
	PMFC_BAD_HW_ADDR
)

// ofp_queue_op_failed_code 1.0
const (
	QOFC_BAD_PORT = iota
	QOFC_BAD_QUEUE
	QOFC_EPERM
)

// END: ofp10 - 5.4.4
// END: ofp10 - 5.4
