package members

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

type ServiceServer interface {
	MemberServiceServer
	SwimServiceServer
}

func wait(f func() bool) {
	count := 0
	for {
		if f() || count > 10 {
			break
		}
		time.Sleep(10 * time.Millisecond)
		count += 1
	}
}

func startOtherServer(server ServiceServer) chan net.Listener {
	c := make(chan net.Listener)
	addr := ":10010"
	go func() {
		s := grpc.NewServer()
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			// Wait for close
			time.Sleep(100 * time.Millisecond)
			lis, err = net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}
		}
		RegisterMemberServiceServer(s, server)
		RegisterSwimServiceServer(s, server)
		c <- lis
		s.Serve(lis)
	}()
	return c
}

type validMemberServiceServer struct {
	argNodeInfo    *NodeInfo
	argAllNodeInfo *AllNodeInfo
	requirePingFn  func() (*Result, error)

	calledMethods map[string]bool
}

func newValidMemberServiceServer() *validMemberServiceServer {
	return &validMemberServiceServer{
		calledMethods: map[string]bool{},
	}
}

func (m *validMemberServiceServer) Join(c context.Context, ni *NodeInfo) (*AllNodeInfo, error) {
	m.argNodeInfo = ni
	m.calledMethods["Join"] = true
	return &AllNodeInfo{Nodes: makeNodes()}, nil
}
func (m *validMemberServiceServer) NotifyNode(c context.Context, ani *AllNodeInfo) (*AllNodeInfo, error) {
	m.argAllNodeInfo = ani
	m.calledMethods["NotifyNode"] = true
	return &AllNodeInfo{Nodes: makeNodes()}, nil
}
func (m *validMemberServiceServer) Ping(c context.Context, ni *NodeInfo) (*AckPing, error) {
	m.argNodeInfo = ni
	m.calledMethods["Ping"] = true
	return &AckPing{IsJoined: true}, nil
}
func (m *validMemberServiceServer) RequirePing(c context.Context, ni *NodeInfo) (*Result, error) {
	m.argNodeInfo = ni
	m.calledMethods["RequirePing"] = true
	if m.requirePingFn != nil {
		return m.requirePingFn()
	}
	return &Result{Success: true}, nil
}
func (m *validMemberServiceServer) Suspect(c context.Context, ni *NodeInfo) (*Empty, error) {
	m.argNodeInfo = ni
	m.calledMethods["Suspect"] = true
	return &Empty{}, nil
}
func (m *validMemberServiceServer) Alive(c context.Context, ni *NodeInfo) (*Empty, error) {
	m.argNodeInfo = ni
	m.calledMethods["Alive"] = true
	return &Empty{}, nil
}
func (m *validMemberServiceServer) Confirm(c context.Context, ni *NodeInfo) (*Empty, error) {
	m.argNodeInfo = ni
	m.calledMethods["Confirm"] = true
	return &Empty{}, nil
}

func makeNodes() []*NodeInfo {
	nodes := []*NodeInfo{}
	return append(nodes, &NodeInfo{Address: ":10010", IncarnationNumber: 0})
}

type timeoutMemberServiceServer struct{}

func (m *timeoutMemberServiceServer) Join(context.Context, *NodeInfo) (*AllNodeInfo, error) {
	time.Sleep(100 * time.Second)
	return nil, nil
}
func (m *timeoutMemberServiceServer) NotifyNode(context.Context, *AllNodeInfo) (*AllNodeInfo, error) {
	time.Sleep(100 * time.Second)
	return nil, nil
}
func (m *timeoutMemberServiceServer) Ping(c context.Context, ni *NodeInfo) (*AckPing, error) {
	time.Sleep(100 * time.Second)
	return nil, nil
}
func (m *timeoutMemberServiceServer) RequirePing(c context.Context, ni *NodeInfo) (*Result, error) {
	time.Sleep(100 * time.Second)
	return nil, nil
}
func (m *timeoutMemberServiceServer) Suspect(c context.Context, ni *NodeInfo) (*Empty, error) {
	time.Sleep(100 * time.Second)
	return nil, nil
}
func (m *timeoutMemberServiceServer) Alive(c context.Context, ni *NodeInfo) (*Empty, error) {
	time.Sleep(100 * time.Second)
	return nil, nil
}
func (m *timeoutMemberServiceServer) Confirm(c context.Context, ni *NodeInfo) (*Empty, error) {
	time.Sleep(100 * time.Second)
	return nil, nil
}

func TestNewSwimManager(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	if len(s.members) != 0 {
		t.Error("members already exists")
	}
	if s.myInfo == nil {
		t.Error("myInfo is not created")
	}
	if (s.myInfo.Addr.Addr == s.myInfo.Name && s.myInfo.Name == "TEST_ADDR") == false {
		t.Error("myInfo is not initialized with name and address")
	}
}

func TestSwimManager_Start(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	err := s.Start(grpc.NewServer(), "")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestSwimManager_StartWithAddress(t *testing.T) {
	c := startOtherServer(newValidMemberServiceServer())
	s := NewSwimManager("TEST_ADDR")
	lis := <-c
	defer lis.Close()
	err := s.Start(grpc.NewServer(), ":10010")
	if err != nil {
		t.Error(err)
	}
	if len(s.members) != 1 {
		t.Error("the number of members is not 1")
	}
	if s.members[":10010"] == nil {
		t.Error("member is not joined")
	}
}

func TestSwimManager_StartWithTimeoutMember(t *testing.T) {
	c := startOtherServer(&timeoutMemberServiceServer{})
	s := NewSwimManager("TEST_ADDR")
	s.pingInterval = 100 * time.Millisecond

	lis := <-c
	defer lis.Close()
	err := s.Start(grpc.NewServer(), ":10010")
	if err == nil {
		t.Error("timeout error doesn't occur")
	}
	if len(s.members) != 0 {
		t.Error("the number of members is not 0")
	}
}

func TestSwimManager_Join(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	ani, err := s.Join(nil, &NodeInfo{
		Address:           "addr",
		IncarnationNumber: 10,
	})

	if err != nil {
		t.Error(err)
	}
	if len(ani.Nodes) != 1 {
		t.Error("returns unknown node")
	}
	node := ani.Nodes[0]
	if (node.Address == "addr" && node.IncarnationNumber == 10) == false {
		t.Error("returns unknown member")
	}
}

func TestSwimManager_JoinKnowingOther(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	s.members["dummy"] = NewMember("dummy", "dummy", 5)
	ani, err := s.Join(nil, &NodeInfo{
		Address:           "addr",
		IncarnationNumber: 10,
	})

	if err != nil {
		t.Error(err)
	}
	if len(ani.Nodes) != 2 {
		t.Error("returns unknown node")
	}

	for _, node := range ani.Nodes {
		if node.Address == "dummy" {
			if node.IncarnationNumber != 5 {
				t.Error("returns unknown member")
			}
		} else if node.Address == "addr" {
			if node.IncarnationNumber != 10 {
				t.Error("returns unknown member")
			}
		} else {
			t.Error("returns unknown member")
		}
	}
}

func TestSwimManager_JoinDismiss(t *testing.T) {
	ss := newValidMemberServiceServer()
	c := startOtherServer(ss)
	s := NewSwimManager("TEST_ADDR")

	listener := <-c
	defer listener.Close()
	s.members[":10010"] = NewMember(":10010", ":10010", 5)

	s.Join(nil, &NodeInfo{
		Address:           "addr",
		IncarnationNumber: 10,
	})

	wait(func() bool {
		return ss.calledMethods["Join"] != false
	})
	if ss.argAllNodeInfo == nil {
		t.Error("not disseminated")
	}
	hasNewNode := false
	for _, node := range ss.argAllNodeInfo.Nodes {
		if node.Address == "addr" {
			hasNewNode = true
		}
	}
	if hasNewNode == false {
		t.Error("not including information of new node")
	}
	if len(ss.argAllNodeInfo.Nodes) != 2 {
		t.Error("the number of nodes are differnt")
	}
}

func TestSwimManager_NotifyNode(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	ani, err := s.NotifyNode(nil, &AllNodeInfo{[]*NodeInfo{{
		Address:           "addr",
		IncarnationNumber: 10,
	}}})

	if err != nil {
		t.Error(err)
	}
	if len(ani.Nodes) != 1 {
		t.Error("returns unknown node")
	}
	node := ani.Nodes[0]
	if (node.Address == "addr" && node.IncarnationNumber == 10) == false {
		t.Error("returns unknown member")
	}
}

func TestSwimManager_NotifyNodeAlreadyKnown(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	m := NewMember("addr", "addr", 1)
	s.members["addr"] = m

	ni := s.buildNodeInfo(m)
	ni.IncarnationNumber = 10
	ani, err := s.NotifyNode(nil, &AllNodeInfo{[]*NodeInfo{ni}})

	if err != nil {
		t.Error(err)
	}
	if len(ani.Nodes) != 1 {
		t.Error("returns unknown node")
	}
	node := ani.Nodes[0]
	if (node.Address == "addr" && node.IncarnationNumber == 10) == false {
		t.Error("returns unknown member")
	}
}

func TestSwimManager_Ping(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	s.members["addr"] = NewMember("addr", "addr", 1)

	ack, err := s.Ping(nil, &NodeInfo{Address: "addr"})
	if err != nil {
		t.Error(err)
	}
	if ack.IsJoined == false {
		t.Error("not registered")
	}
}

func TestSwimManager_PingFromUnknownNode(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	ack, err := s.Ping(nil, &NodeInfo{Address: "addr"})
	if err != nil {
		t.Error(err)
	}
	if ack.IsJoined == true {
		t.Error("lies node is registered")
	}
}

func TestSwimManager_execPing(t *testing.T) {
	ss := newValidMemberServiceServer()
	c := startOtherServer(ss)
	s := NewSwimManager("TEST_ADDR")

	listener := <-c
	defer listener.Close()
	err := s.execPing(NewMember(":10010", ":10010", 5))
	if err != nil {
		t.Error(err)
	}
}

func TestSwimManager_RequirePing(t *testing.T) {
	c := startOtherServer(newValidMemberServiceServer())
	s := NewSwimManager("TEST_ADDR")

	lis := <-c
	defer lis.Close()
	res, err := s.RequirePing(nil, &NodeInfo{Address: ":10010"})
	if err != nil {
		t.Error(err)
	}

	if res.Success != true {
		t.Error("not success")
	}
}

func TestSwimManager_RequirePingTimeout(t *testing.T) {
	c := startOtherServer(&timeoutMemberServiceServer{})
	s := NewSwimManager("TEST_ADDR")
	s.pingInterval = 100 * time.Millisecond

	lis := <-c
	defer lis.Close()
	res, err := s.RequirePing(nil, &NodeInfo{Address: ":10010"})
	if err != nil {
		t.Error(err)
	}

	if res.Success != false {
		t.Error("says success when timeout")
	}
}

func TestSwimManager_RequirePingInvalidNode(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	_, err := s.RequirePing(nil, &NodeInfo{Address: ":100"})
	if err == nil {
		t.Error("not returned error")
	}
}

func TestSwimManager_execRequirePing(t *testing.T) {
	ss := newValidMemberServiceServer()
	c := startOtherServer(ss)
	s := NewSwimManager("TEST_ADDR")
	s.addMember(NewMember(":10010", ":10010", 0))
	errMember := NewMember("addr", "addr", 3)
	s.addMember(errMember)

	lis := <-c
	defer lis.Close()
	s.execRequirePing(errMember)
	wait(func() bool {
		return ss.calledMethods["RequirePing"] != false
	})

	argNode := ss.argNodeInfo
	if (argNode.Address == "addr" && argNode.IncarnationNumber == 3) == false {
		t.Error("require ping not reached")
	}
	if s.members["addr"].Status != ALIVE {
		t.Error("require ping not alived")
	}
}

func TestSwimManager_execRequirePingFail(t *testing.T) {
	ss := newValidMemberServiceServer()
	ss.requirePingFn = func() (*Result, error) {
		return &Result{Success: false}, nil
	}
	c := startOtherServer(ss)
	s := NewSwimManager("TEST_ADDR")
	s.addMember(NewMember(":10010", ":10010", 0))
	errMember := NewMember("addr", "addr", 3)
	s.addMember(errMember)

	lis := <-c
	defer lis.Close()
	s.execRequirePing(errMember)
	wait(func() bool {
		return ss.calledMethods["RequirePing"] != false
	})

	argNode := ss.argNodeInfo
	if (argNode.Address == "addr" && argNode.IncarnationNumber == 3) == false {
		t.Error("require ping not reached")
	}
	if s.members["addr"].Status != SUSPECT {
		t.Error("require ping not suspect")
	}
}

func TestSwimManager_execRequirePingSuspectMemberFail(t *testing.T) {
	ss := newValidMemberServiceServer()
	ss.requirePingFn = func() (*Result, error) {
		return &Result{Success: false}, nil
	}
	c := startOtherServer(ss)
	s := NewSwimManager("TEST_ADDR")
	s.addMember(NewMember(":10010", ":10010", 0))
	errMember := NewMember("addr", "addr", 3)
	errMember.Status = SUSPECT
	s.addMember(errMember)

	lis := <-c
	defer lis.Close()
	s.execRequirePing(errMember)
	wait(func() bool {
		return ss.calledMethods["RequirePing"] != false
	})

	argNode := ss.argNodeInfo
	if (argNode.Address == "addr" && argNode.IncarnationNumber == 3) == false {
		t.Error("require ping not reached")
	}
	if s.members["addr"] != nil {
		t.Error("suspect member is not deleted")
	}
}

func TestSwimManager_Suspect(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	s.addMember(NewMember(":10010", ":10010", 0))
	_, err := s.Suspect(nil, &NodeInfo{Address: ":10010"})
	if err != nil {
		t.Error("not returned error")
	}
	if s.members[":10010"].Status != SUSPECT {
		t.Error("not set suspect")
	}
}

func TestSwimManager_SuspectUnknownMember(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	_, err := s.Suspect(nil, &NodeInfo{Address: ":10010"})
	if err != nil {
		t.Error("not returned error")
	}
}

func TestSwimManager_execSuspect(t *testing.T) {
	ss := newValidMemberServiceServer()
	c := startOtherServer(ss)
	s := NewSwimManager("TEST_ADDR")
	s.addMember(NewMember(":10010", ":10010", 0))

	lis := <-c
	defer lis.Close()
	s.execSuspect(NewMember("addr", "addr", 3))
	wait(func() bool {
		return ss.calledMethods["Suspect"] != false
	})

	argNode := ss.argNodeInfo
	if (argNode.Address == "addr" && argNode.IncarnationNumber == 3) == false {
		t.Error("require ping not reached")
	}
}

func TestSwimManager_Alive(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	s.addMember(NewMember(":10010", ":10010", 0))
	_, err := s.Alive(nil, &NodeInfo{Address: ":10010"})
	if err != nil {
		t.Error("not returned error")
	}
	if s.members[":10010"].Status != ALIVE {
		t.Error("not set alive")
	}
}

func TestSwimManager_AliveUnknownMember(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	_, err := s.Alive(nil, &NodeInfo{Address: ":10010"})
	if err != nil {
		t.Error("not returned error")
	}
}

func TestSwimManager_execAlive(t *testing.T) {
	ss := newValidMemberServiceServer()
	c := startOtherServer(ss)
	s := NewSwimManager("TEST_ADDR")
	s.addMember(NewMember(":10010", ":10010", 0))

	lis := <-c
	defer lis.Close()
	s.execAlive(NewMember("addr", "addr", 3))
	wait(func() bool {
		return ss.calledMethods["Alive"] != false
	})

	argNode := ss.argNodeInfo
	if (argNode.Address == "addr" && argNode.IncarnationNumber == 3) == false {
		t.Error("require ping not reached")
	}
}

func TestSwimManager_Confirm(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	s.addMember(NewMember(":10010", ":10010", 0))
	_, err := s.Confirm(nil, &NodeInfo{Address: ":10010"})
	if err != nil {
		t.Error("not returned error")
	}
	if _, ok := s.members[":10010"]; ok != false {
		t.Error("not deleted")
	}
}

func TestSwimManager_ConfirmUnknownMember(t *testing.T) {
	s := NewSwimManager("TEST_ADDR")
	_, err := s.Confirm(nil, &NodeInfo{Address: ":10010"})
	if err != nil {
		t.Error("not returned error")
	}
}

func TestSwimManager_execConfirm(t *testing.T) {
	ss := newValidMemberServiceServer()
	c := startOtherServer(ss)
	s := NewSwimManager("TEST_ADDR")
	s.addMember(NewMember(":10010", ":10010", 0))

	lis := <-c
	defer lis.Close()
	s.execConfirm(NewMember("addr", "addr", 3))
	wait(func() bool {
		return ss.calledMethods["Confirm"] != false
	})

	argNode := ss.argNodeInfo
	if (argNode.Address == "addr" && argNode.IncarnationNumber == 3) == false {
		t.Error("require ping not reached")
	}
}
