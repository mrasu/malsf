package members

import "testing"

func TestNewMember(t *testing.T) {
	m := NewMember("TEST_NAME", "TEST_ADDR", 10)
	if m.Name != "TEST_NAME" {
		t.Error("name is invalid")
	}
	if m.Addr.Addr != "TEST_ADDR" {
		t.Error("addr is invalid")
	}
	if m.IncarnationNumber != 10 {
		t.Error("IncarnationNumber is invalid")
	}
	if m.Status != ALIVE {
		t.Error("status is not alive")
	}
}

func TestMember_Address(t *testing.T) {
	m := NewMember("name", "addr", 10)
	if m.Address() != "addr" {
		t.Error("Address() not return address")
	}
}

func TestMember_Connect(t *testing.T) {
	m := NewMember("name", "addr", 10)
	conn, err := m.Connect()
	if err != nil {
		t.Errorf("error contains: %s", err)
	}
	if conn == nil {
		t.Error("Connection not returned")
	}
	conn.Close()
}

func TestMember_String(t *testing.T) {
	m := NewMember("name", "addr", 10)
	expected := "Name(name), Addr(addr), IncarnationNumber(10), Status(0)"
	if m.String() != expected {
		t.Errorf("Name() returns invalid text: %s", m.String())
	}
	m.Status = SUSPECT
	expected = "Name(name), Addr(addr), IncarnationNumber(10), Status(1)"
	if m.String() != expected {
		t.Errorf("Name() returns invalid text: %s", m.String())
	}

}
