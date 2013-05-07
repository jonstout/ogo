package ofp10

import (
	"testing"
)

func TestActionOutput(t *testing.T) {
	act := NewActionOutput()
	if act.Type != OFPAT_OUTPUT {
		t.Error("Action type was:", act.Type, "instead of:", OFPAT_OUTPUT)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(OfpActionOutput)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionEnqueue(t *testing.T) {
	act := NewActionEnqueue()
	if act.Type != OFPAT_ENQUEUE {
		t.Error("Action type was:", act.Type, "instead of:", OFPAT_ENQUEUE)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(OfpActionEnqueue)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionVLANVID(t *testing.T) {
	act := NewActionVLANVID()
	if act.Type != OFPAT_SET_VLAN_VID {
		t.Error("Action type was:", act.Type, "instead of:", OFPAT_SET_VLAN_VID)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(OfpActionVLANVID)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionVLANVPCP(t *testing.T) {
	act := NewActionVLANPCP()
	if act.Type != OFPAT_SET_VLAN_PCP {
		t.Error("Action type was:", act.Type, "instead of:", OFPAT_SET_VLAN_PCP)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(OfpActionVLANPCP)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionDLAddr(t *testing.T) {
	act := NewActionDLSrc()
	if act.Type != OFPAT_SET_DL_SRC {
		t.Error("Action type was:", act.Type, "instead of:", OFPAT_SET_DL_SRC)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(OfpActionDLAddr)
	act2.Write(b)
	if act.Type != act2.Type {
		t.Error("Encode / Decode - Type:", act, act2)
	}
	if act.Length != act2.Length {
		t.Error("Encode / Decode - Length:", act, act2)
	}
	if act.DLAddr.String() != act2.DLAddr.String() {
		t.Error("Encode / Decode - DLAddr:", act, act2)
	}
	if act.Pad != act2.Pad {
		t.Error("Encode / Decode - Pad:", act, act2)
	}
}

func TestActionNWAddr(t *testing.T) {
	act := NewActionNWSrc()
	if act.Type != OFPAT_SET_NW_SRC {
		t.Error("Action type was:", act.Type, "instead of:", OFPAT_SET_NW_SRC)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(OfpActionNWAddr)
	act2.Write(b)
	if act.Type != act2.Type {
		t.Error("Encode / Decode - Type:", act, act2)
	}
	if act.Length != act2.Length {
		t.Error("Encode / Decode - Length:", act, act2)
	}
	if act.NWAddr.String() != act2.NWAddr.String() {
		t.Error("Encode / Decode - DLAddr:", act, act2)
	}
}

func TestActionNWTOS(t *testing.T) {
	act := NewActionNWTOS()
	if act.Type != OFPAT_SET_NW_TOS {
		t.Error("Action type was:", act.Type, "instead of:", OFPAT_SET_NW_TOS)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(OfpActionNWTOS)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionTPPort(t *testing.T) {
	act := NewActionTPSrc()
	if act.Type != OFPAT_SET_TP_SRC {
		t.Error("Action type was:", act.Type, "instead of:", OFPAT_SET_TP_SRC)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(OfpActionTPPort)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionVendorPort(t *testing.T) {
	act := NewActionVendorPort()
	if act.Type != OFPAT_VENDOR {
		t.Error("Action type was:", act.Type, "instead of:", OFPAT_VENDOR)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(OfpActionVendorPort)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}
