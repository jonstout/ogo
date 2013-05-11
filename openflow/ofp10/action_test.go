package ofp10

import (
	"testing"
)

func TestActionOutput(t *testing.T) {
	act := NewActionOutput()
	if act.Type != AT_OUTPUT {
		t.Error("Action type was:", act.Type, "instead of:", AT_OUTPUT)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(ActionOutput)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionEnqueue(t *testing.T) {
	act := NewActionEnqueue()
	if act.Type != AT_ENQUEUE {
		t.Error("Action type was:", act.Type, "instead of:", AT_ENQUEUE)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(ActionEnqueue)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionVLANVID(t *testing.T) {
	act := NewActionVLANVID()
	if act.Type != AT_SET_VLAN_VID {
		t.Error("Action type was:", act.Type, "instead of:", AT_SET_VLAN_VID)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(ActionVLANVID)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionVLANVPCP(t *testing.T) {
	act := NewActionVLANPCP()
	if act.Type != AT_SET_VLAN_PCP {
		t.Error("Action type was:", act.Type, "instead of:", AT_SET_VLAN_PCP)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(ActionVLANPCP)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionDLAddr(t *testing.T) {
	act := NewActionDLSrc()
	if act.Type != AT_SET_DL_SRC {
		t.Error("Action type was:", act.Type, "instead of:", AT_SET_DL_SRC)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(ActionDLAddr)
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
	if act.Type != AT_SET_NW_SRC {
		t.Error("Action type was:", act.Type, "instead of:", AT_SET_NW_SRC)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(ActionNWAddr)
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
	if act.Type != AT_SET_NW_TOS {
		t.Error("Action type was:", act.Type, "instead of:", AT_SET_NW_TOS)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(ActionNWTOS)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionTPPort(t *testing.T) {
	act := NewActionTPSrc()
	if act.Type != AT_SET_TP_SRC {
		t.Error("Action type was:", act.Type, "instead of:", AT_SET_TP_SRC)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(ActionTPPort)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}

func TestActionVendorPort(t *testing.T) {
	act := NewActionVendorPort()
	if act.Type != AT_VENDOR {
		t.Error("Action type was:", act.Type, "instead of:", AT_VENDOR)
	}
	b := make([]byte, act.Len())
	act.Read(b)
	
	act2 := new(ActionVendorPort)
	act2.Write(b)
	if *act != *act2 {
		t.Error("Encode / Decode:", act, act2)
	}
}
