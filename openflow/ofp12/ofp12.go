/*
 * Jonathan M. Stout 2012
 * ofp.go
 * OpenFlow Specification 1.2-rc2
*/

package ofp12
//-------------------------------------------------
// BEGIN: OpenFlow 1.2-rc2 const
//-------------------------------------------------
// ofp_flow_mod_command
const (
      OFPFC_ADD = iota // OFPFC_ADD == 0
      OFPFC_MODIFY
      OFPFC_MODIFY_STRICT
      OFPFC_DELETE
      OFPFC_DELETE_STRICT
)

// ofp_group_mod_command
const (
      OFPGC_ADD = iota
      OFPGC_MODIFY
      OFPGC_DELETE
)

// ofp_type
const (
	/* Immutable messages. */
	OFPT_HELLO = iota
	OFPT_ERROR
	OFPT_ECHO_REQUEST
	OFPT_ECHO_REPLY
	OFPT_EXPERIMENTER

	/* Switch configuration messages. */
	OFPT_FEATURES_REQUEST
	OFPT_FEATURES_REPLY
	OFPT_GET_CONFIG_REQUEST
	OFPT_GET_CONFIG_REPLY
	OFPT_SET_CONFIG

	/* Asynchronous messages. */
	OFPT_PACKET_IN
	OFPT_FLOW_REMOVED
	OFPT_PORT_STATUS

	/* Controller command messages. */
	OFPT_PACKET_OUT
	OFPT_FLOW_MOD
	OFPT_GROUP_MOD
	OFPT_PORT_MOD
	OFPT_TABLE_MOD

	/* Statistics messages. */
	OFPT_STATS_REQUEST
	OFPT_STATS_REPLY

	/* Barrier messages. */
	OFPT_BARRIER_REQUEST
	OFPT_BARRIER_REPLY

	/* Queue Configuration messages. */
	OFPT_QUEUE_GET_CONFIG_REQUEST
	OFPT_QUEUE_GET_CONFIG_REPLY

	/* Controller role change request messages. */
	OFPT_ROLE_REQUEST
	OFPT_ROLE_REPLY
)

// ofp_port_config
const (
	OFPPC_PORT_DOWN = 1 << 0
	OFPPC_NO_RECV = 1 << 2
	OFPPC_NO_FWD = 1 << 5
	OFPPC_NO_PACKET_IN = 1 << 6
)

// ofp_port_state
const (
	OFPPS_LINK_DOWN = 1 << 0
	OFPPS_BLOCKED = 1 << 1
	OFPPS_LIVE = 1 << 2
)

// ofp_port_no
const (
	OFPP_MAX = 0Xffffff00
	OFPP_IN_PORT = 0xfffffff8
	OFPP_TABLE = 0xfffffff9
	OFPP_NORMAL = 0xfffffffa
	OFPP_FLOOD = 0xfffffffb
	OFPP_ALL = 0xfffffffc
	OFPP_CONTROLLER = 0xfffffffd
	OFPP_LOCAL = 0xfffffffe
	OFPP_ANY = 0xffffffff
)

// ofp_port_features
const (
	OFPPF_10MB_HD = 1 << 0
	OFPPF_10MB_FD = 1 << 1
	OFPPF_100MB_HD = 1 << 2
	OFPPF_100MB_FD = 1 << 3
	OFPPF_1GB_HD = 1 << 4
	OFPPF_1GB_FD = 1 << 5
	OFPPF_10GB_FD = 1 << 6
	OFPPF_40GB_FD = 1 << 7
	OFPPF_100GB_FD = 1 << 8
	OFPPF_1TB_FD = 1 << 9
	OFPPF_OTHER = 1 << 10

	OFPPF_COPPER = 1 << 11
	OFPPF_FIBER = 1 << 12
	OFPPF_AUTONEG = 1 << 13
	OFPPF_PAUSE = 1 << 14
	OFPPF_PAUSE_ASYM = 1 << 15
)

// ofp_queue properties
const (
	OFPQT_MIN_RATE = iota
	OFPQT_MAX_RATE
	OFPQT_EXPERIMENTER = 0xffff
)

// ofp_match_type
const (
	OFPMT_STANDARD = iota
	OFPMT_OXM
)

// ofp_oxm_class
const (
	OFPXMC_NXM_0 = 0x0000
	OFPXMC_NXM_1 = 0x0001
	OFPXMC_OPENFLOW_BASIC = 0x8000
	OFPXMC_EXPERIMENTER = 0xffff
)

// enum oxm_ofb_match_fields
const (
	OFPXMT_OFB_IN_PORT = 0
	OFPXMT_OFB_IN_PHY_PORT = 1
	OFPXMT_OFB_METADATA = 2
	OFPXMT_OFB_ETH_DST = 3
	OFPXMT_OFB_ETH_SRC = 4
	OFPXMT_OFB_ETH_TYPE = 5
	OFPXMT_OFB_VLAN_VID = 6
	OFPXMT_OFB_VLAN_PCP = 7
	OFPXMT_OFB_IP_TOS = 8
	OFPXMT_OFB_IP_ECN = 9
	OFPXMT_OFB_IP_PROTO = 10
	OFPXMT_OFB_IPV4_SRC = 11
	OFPXMT_OFB_IPV4_DST = 12
	OFPXMT_OFB_TCP_SRC = 13
	OFPXMT_OFB_TCP_DST = 14
	OFPXMT_OFB_UDP_SRC = 15
	OFPXMT_OFB_UDP_DST = 16
	OFPXMT_OFB_SCTP_SRC = 17
	OFPXMT_OFB_SCTP_DST = 18
	OFPXMT_OFB_ICMPV4_TYPE = 19
	OFPXMT_OFB_ICMPV4_CODE = 20
	OFPXMT_OFB_ARP_OP = 21
	OFPXMT_OFB_ARP_SPA = 22
	OFPXMT_OFB_ARP_TPA = 23
	OFPXMT_OFB_ARP_SHA = 24
	OFPXMT_OFB_ARP_THA = 25
	OFPXMT_OFB_IPV6_SRC = 26
	OFPXMT_OFB_IPV6_DST = 27
	OFPXMT_OFB_IPV6_FLABEL = 28
	OFPXMT_OFB_ICMPV6_TYPE = 29
	OFPXMT_OFB_ICMPV6_CODE = 30
	OFPXMT_OFB_IPV6_ND_TARGET = 31
	OFPXMT_OFB_IPV6_ND_SLL = 32
	OFPXMT_OFB_IPV6_ND_TLL = 33
	OFPXMT_OFB_MPLS_LABEL = 34
	OFPXMT_OFB_MPLS_TC = 35
)

// ofp_vlan_id
const (
	OFPVID_PRESENT = 0x1000
	OFPVID_NONE = 0x000
)

// enum ofp_instruction_type
const (
	OFPIT_GOTO_TABLE = 1
	OFPIT_WRITE_METADATA = 2
	OFPIT_WRITE_ACTIONS = 3
	OFPIT_APPLY_ACTIONS = 4
	OFPIT_CLEAR_ACTIONS = 5
	OFPIT_EXPERIMENTER = 0xFFFF
)

// enum ofp_action_type
const (
	OFPAT_OUTPUT = iota
	OFPAT_OBSOLETE_1
	OFPAT_OBSOLETE_2
	OFPAT_OBSOLETE_3
	OFPAT_OBSOLETE_4
	OFPAT_OBSOLETE_5
	OFPAT_OBSOLETE_6
	OFPAT_OBSOLETE_7
	OFPAT_OBSOLETE_8
	OFPAT_OBSOLETE_9
	OFPAT_OBSOLETE_10
	OFPAT_COPY_TTL_OUT
	OFPAT_COPY_TTL_IN
	OFPAT_OBSOLETE_12
	OFPAT_SET_MPLS_TTL
	OFPAT_DEC_MPLS_TTL
	OFPAT_PUSH_VLAN
	OFPAT_POP_VLAN
	OFPAT_PUSH_MPLS
	OFPAT_POP_MPLS
	OFPAT_SET_QUEUE
	OFPAT_GROUP
	OFPAT_SET_NW_TTL
	OFPAT_DEC_NW_TTL
	OFPAT_SET_FIELD
	OFPAT_EXPERIMENTER = 0xffff
)

// ofp_controller_max_len
const (
	OFPCML_MAX = 0xffe5
	OFPCML_NO_VUFFER = 0xffff
)

// ofp_config_flags
const (
	OFPC_FRAG_NORMAL = 0
	OFPC_FRAG_DROP = 1 << 0
	OFPC_FRAG_REASM = 1 << 1
	OFPC_FRAG_MASK = 3
	OFPC_INVALID_TTL_TO_CONTROLLER 1 << 2
)

// ofp_table
const (
	OFPTT_MAX = 0xfe
	OFPTT_ALL = 0xff
)

// ofp_table_config
const (
	OFPTC_TABLE_MISS_CONTROLLER = 0
	OFPTC_TABLE_MISS_COUNTER = 1 << 0
	OFPTC_TABLE_MISS_DROP = 1 << 1
	OFPTC_TABLE_MISS_MASK = 3
)

// ofp_flow_mod_flags
const (
	OFPFF_SEND_FLOW_REM = 1 << 0
	OFPFF_CHECK_OVERLAP = 1 << 1
	OFPFF_RESET_COUNTS = 1 << 2
)

// ofp_group_type
const (
	OFPGT_ALL = iota
	OFPGT_SELECT
	OFPGT_INDIRECT
	OFPGT_FF
)

// ofp_stats_types
const (
	/* Description of this OpenFlow switch.
	* The request body is empty.
	* The reply body is struct ofp_desc_stats. */
	OFPST_DESC = iota
	/* Individual flow statistics.
	* The request body is struct ofp_flow_stats_request.
	* The reply body is an array of struct ofp_flow_stats. */
	OFPST_FLOW
	/* Aggregate flow statistics.
	* The request body is struct ofp_aggregate_stats_request.
	* The reply body is struct ofp_aggregate_stats_reply. */
	OFPST_AGGREGATE
	/* Flow table statistics.
	* The request body is empty.
	* The reply body is an array of struct ofp_table_stats. */
	OFPST_TABLE
	/* Port statistics.
	* The request body is struct ofp_port_stats_request.
	* The reply body is an array of struct ofp_port_stats. */
	OFPST_PORT
	/* Queue statistics for a port
	* The request body is struct ofp_queue_stats_request.
	* The reply body is an array of struct ofp_queue_stats */
	OFPST_QUEUE
	/* Group counter statistics.
	* The request body is struct ofp_group_stats_request.
	* The reply is an array of struct ofp_group_stats. */
	OFPST_GROUP
	/* Group description statistics.
	* The request body is empty.
	* The reply body is an array of struct ofp_group_desc_stats. */
	OFPST_GROUP_DESC
	/* Group features.
	* The request body is empty.
	* The reply body is struct ofp_group_features_stats. */
	OFPST_GROUP_FEATURES
	/* Experimenter extension.
	* The request and reply bodies begin with
	* struct ofp_experimenter_stats_header.
	* The request and reply bodies are otherwise experimenter-defined. */
	OFPST_EXPERIMENTER = 0xffff
)

// ofp_group_capabilities
const (
	OFPGFC_SELECT_WEIGHT = iota
	OFPGFC_SELECT_LIVENESS
	OFPGFC_CHAINING
	OFPGFC_CHAINING_CHECKS
)

// ofp_controller_role
const (
	OFPCR_ROLE_NOCHANGE = iota
	OFPCR_ROLE_EQUAL
	OFPCR_ROLE_MASTER
	OFPCR_ROLE_SLAVE
)

// ofp_packet_in_reason
const (
	OFPR_NO_MATCH = iota
	OFPR_ACTION
	OFPR_INVALID_TTL
)

// ofp_flow_removed_reason
const (
	OFPRR_IDLE_TIMEOUT = iota
	OFPRR_HARD_TIMEOUT
	OFPRR_DELETE
	OFPRR_GROUP_DELETE
)

// ofp_port_reason
const (
	OFPPR_ADD = iota
	OFPPR_DELETE
	OFPPR_MODIFY
)

// ofp_error_type
const (
	OFPET_HELLO_FAILED
	OFPET_BAD_REQUEST
	OFPET_BAD_ACTION
	OFPET_BAD_INSTRUCTION
	OFPET_BAD_MATCH
	OFPET_FLOW_MOD_FAILED
	OFPET_GROUP_MOD_FAILED
	OFPET_PORT_MOD_FAILED
	OFPET_TABLE_MOD_FAILED
	OFPET_QUEUE_OP_FAILED
	OFPET_SWITCH_CONFIG_FAILED
	OFPET_ROLE_REQUEST_FAILED
	OFPET_EXPERIMENTER = 0xffff
)

// ofp_hello_failed_code
const (
	OFPHFC_INCOMPATIBLE = iota
	OFPHFC_EPERM
)

// ofp_bad_request_code
const (
	OFPBRC_BAD_VERSION = iota
	OFPBRC_BAD_TYPE
	OFPBRC_BAD_STAT
	OFPBRC_BAD_EXPERIMENTER
	OFPBRC_BAD_EXP_TYPE
	OFPBRC_EPERM
	OFPBRC_BAD_LEN
	OFPBRC_BUFFER_EMPTY
	OFPBRC_BUFFER_UNKNOWN
	OFPBRC_BAD_TABLE_ID
	OFPBRC_IS_SLAVE
)

// ofp_bad_action_code
const (
	OFPBAC_BAD_TYPE = iota
	OFPBAC_BAD_LEN
	OFPBAC_BAD_EXPERIMENTER
	OFPBAC_BAD_EXP_TYPE
	OFPBAC_BAD_OUT_PORT
	OFPBAC_BAD_ARGUMENT
	OFPBAC_EPERM
	OFPBAC_TOO_MANY
	OFPBAC_BAD_QUEUE
	OFPBAC_BAD_OUT_GROUP
	OFPBAC_MATCH_INCONSISTENT
	OFPBAC_UNSUPPORTED_ORDER
	OFPBAC_BAD_TAG
	OpenFlow Switch Specification
	OFPBAC_BAD_SET_TYPE
	OFPBAC_BAD_SET_LEN
	OFPBAC_BAD_SET_ARGUMENT
)

// ofp_bad_instruction_code
const (
	OFPBIC_UNKNOW_INST = iota
	OFPBIC_UNSUP_INST
	OFPBIC_BAD_TABLE_ID
	OFPBIC_UNSUP_METADATA
	OFPBIC_UNSUP_METADATA_MASK
	OFPBIC_BAD_EXPERIMENTER
	OFPBIC_BAD_EXP_TYPE
)

// ofp_bad_match_code
const (
	OFPBMC_BAD_TYPE = iota /* Unsupported match type specified by the match */
	OFPBMC_BAD_LEN /* Length problem in match. */
	OFPBMC_BAD_TAG /* Match uses an unsupported tag/encap. */
	OFPBMC_BAD_DL_ADDR_MASK /* Unsupported datalink addr mask - switch does not support arbitrary datalink address mask. */
	OFPBMC_BAD_NW_ADDR_MASK /* Unsupported network addr mask - switch does not support arbitrary network address mask. */
	OFPBMC_BAD_WILDCARDS /* Unsupported combination of fields masked or omitted in the match. */
	OFPBMC_BAD_FIELD /* Unsupported field type in the match. */
	OFPBMC_BAD_VALUE /* Unsupported value in a match field. */
	OFPBMC_BAD_MASK /* Unsupported mask specified in the match, field is not dl-address or nw-address. */
	OFPBMC_BAD_PREREQ /* A prerequisite was not met. */
	OFPBMC_DUP_FIELD /* A field type was duplicated. */
)

// ofp_flow_mod_failed_code
const (
	OFPFMFC_UNKNOWN = iota
	OFPFMFC_TABLE_FULL
	OFPFMFC_BAD_TABLE_ID
	OFPFMFC_OVERLAP
	OFPFMFC_EPERM
	OFPFMFC_BAD_TIMEOUT
	OFPFMFC_BAD_COMMAND
)

// ofp_group_mod_failed_code
const (
	OFPGMFC_GROUP_EXISTS = iota
	OFPGMFC_INVALID_GROUP
	OFPGMFC_WEIGHT_UNSUPPORTED
	OFPGMFC_OUT_OF_GROUPS
	OFPGMFC_OUT_OF_BUCKETS
	OFPGMFC_CHAINING_UNSUPPORTED
	OFPGMFC_WATCH_UNSUPPORTED
	OFPGMFC_LOOP
	OFPGMFC_UNKNOWN_GROUP
	OFPGMFC_CHAINED_GROUP
)

// ofp_port_mod_failed_code
const (
	OFPPMFC_BAD_PORT = iota
	OFPPMFC_BAD_HW_ADDR
	OFPPMFC_BAD_CONFIG
	OFPPMFC_BAD_ADVERTISE
)

// ofp_table_mod_failed_code
const (
	OFPTMFC_BAD_TABLE = iota
	OFPTMFC_BAD_CONFIG
)

// ofp_queue_op_failed_code
const (
	OFPQOFC_BAD_PORT = iota
	OFPQOFC_BAD_QUEUE
	OFPQOFC_EPERM
)

// ofp_switch_config_failed_code
const (
	OFPSCFC_BAD_FLAGS = iota
	OFPSCFC_BAD_LEN
)

// ofp_role_request_failed_code
const (
	OFPRRFC_STALE = iota
)
//-------------------------------------------------
// END: OpenFlow 1.2-rc2 const
//-------------------------------------------------

//-------------------------------------------------
// BEGIN: OpenFlow 1.2-rc2 struct
//-------------------------------------------------
type OfpHeader struct {
     version uint8
     type_ uint8
     length uint16
     xid uint32
}

const (
	OFP_ETH_ALEN = 6
	OFP_MAX_PORT_NAME_LEN = 16
)
type OfpPort struct {
     port_no uint32
     pad [4]uint8
     hw_addr [OFP_ETH_ALEN]uint8
     pad2 [2]uint8
     name [OFP_MAX_PORT_NAME_LEN]char
     config uint32
     state uint32
     curr uint32
     advertised uint32
     peer uint32
     curr_speed uint32
     max_speed uint32
}

type OfpPacketQueue struct {
	queue_id uint32
	port uint32
	length uint16
	pad [6]uint8
	properties [0]OfpQueuePropHeader
}

type OfpQueuePropHeader struct {
	property uint16
	length uint16
	pad [4]uint8
}

type OfpQueuePropMinRate struct {
	prop_header OfpQueuePropHeader
	rate uint16
	pad [6]uint8
}

type OfpQueuePropMaxRate struct {
	propHeader OfpQueuePropHeader
	rate uint16
	pad [6]uint8
}

type OfpQueuePropExperimenter struct {
	prop_header OfpQueuePropHeader
	experimenter uint32
	pad [4]uint8
	data [0]uint8
}

type OfpMatch struct {
	type_ uint16
	length uint16
	/* Followed by
	 * -Exactly (length - 4) (possible 0) bytes containing OXM TLVs, then
	 * -Exactly ((length + 7)/8*8 - length) (between 0 and 7) bytes of
	 *  all-zero bytes
	 * In summary, ofp_match is padded as needed, to make its overall size
	 * a multiple of 8, to preserve alignement in structures using it.
	 */
	oxm_fields [4]uint8
}

type OfpOxmExperimenterHeader struct {
	oxm_header uint32
	experimenter uint32
}

type OfpInstructionGotoTable struct {
	type_ uint16
	length uint16
	table_id uint8
	pad [3]uint8
}

type OfpInstructionWriteMetadata struct {
	type_ uint16
	length uint16
	pad [4]uint8
	metadata uint64
	metadata_mask uint64
}

type OfpInstructionActions struct {
	type_ uint16
	length uint16
	pad [4]uint8
	actions [0]OfpActionHeader
}

type OfpActionHeader struct {
	type_ uint16
	length uint16
	pad [4]uint8
}

type OfpActionGroup struct {
	type_ uint16
	length uint16
	group_id uint32
}

type OfpActionSetQueue struct {
	type_ uint16
	length uint16
	queue_id uint32
}

type OfpActionMplsTtl struct {
	type_ uint16
	length uint16
	mpls_ttl uint8
	pad [3]uint8
}

type OfpActionNewTtl struct {
	type_ uint16
	length uint16
	nw_ttl uint8
	pad [3]uint8
}

type OfpActionPush struct {
	type_ uint16
	length uint16
	ethertype uint16
	pad [2]uint8
}

type OfpActionPopMpls struct {
	type_ uint16
	length uint16
	ethertype uint16
	pad [2]uint8
}
//-------------------------------------------------
// END: OpenFlow 1.2-rc2 struct
//-------------------------------------------------