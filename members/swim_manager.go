package members

/**
  Not use IncarnationNumber to detect failure yet
*/

import (
	"github.com/mrasu/malsf/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
	"sync"
	"time"
)

const (
	PING_TIMEOUT  = 1
	PING_INTERVAL = 1
)

type SwimManager struct {
	myInfo      *Member
	members     map[string]*Member
	swimTargets []*Member

	mu           sync.Mutex
	pingInterval time.Duration
}

func NewSwimManager(addr string) *SwimManager {
	me := NewMember(addr, addr, 0)

	return &SwimManager{
		myInfo:       me,
		members:      map[string]*Member{},
		pingInterval: PING_TIMEOUT * time.Second,
	}
}

func (s *SwimManager) Start(server *grpc.Server, firstAddress string) error {
	RegisterMemberServiceServer(server, s)
	RegisterSwimServiceServer(server, s)

	if firstAddress != "" {
		m := NewMember(firstAddress, firstAddress, 0)
		err := s.joinMember(m)
		if err != nil {
			return err
		}
	}
	s.startSwim()
	return nil
}

func (s *SwimManager) startSwim() {
	go func() {
		t := time.NewTicker(PING_INTERVAL * time.Second)

		for {
			select {
			case <-t.C:
				go s.pingRandomMember()
			}
		}
	}()
}

func (s *SwimManager) joinMember(m *Member) error {
	conn, err := m.Connect()
	defer conn.Close()
	if err != nil {
		return err
	}
	c := NewMemberServiceClient(conn)

	myNi := &NodeInfo{
		Address:           s.myInfo.Addr.Addr,
		IncarnationNumber: 0,
	}
	util.LogSwimMethod(true, "JOIN", myNi.String())
	ctx, cancel := context.WithTimeout(context.Background(), s.pingInterval)
	defer cancel()
	ani, err := c.Join(ctx, myNi)
	if err != nil {
		return err
	}

	s.addMembers(ani.Nodes)
	return nil
}

func (s *SwimManager) buildMember(ni *NodeInfo) *Member {
	return NewMember(ni.Address, ni.Address, int(ni.IncarnationNumber))
}

func (s *SwimManager) addMembers(nis []*NodeInfo) []*Member {
	newMembers := []*Member{}
	for _, ni := range nis {
		if ni.Address == s.myInfo.Address() {
			continue
		}

		if m, ok := s.members[ni.Address]; ok {
			if int32(m.IncarnationNumber) < ni.IncarnationNumber {
				m.IncarnationNumber = int(ni.IncarnationNumber)
			}
			continue
		}

		m := s.buildMember(ni)
		s.addMember(m)
		newMembers = append(newMembers, m)
	}
	return newMembers
}

func (s *SwimManager) addMember(m *Member) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.members[m.Name] = m
}

func (s *SwimManager) getMember(ni *NodeInfo) *Member {
	return s.members[ni.Address]
}

func (s *SwimManager) deleteMember(m *Member) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.members, m.Name)
}

func (s *SwimManager) Join(ctx context.Context, ni *NodeInfo) (*AllNodeInfo, error) {
	util.LogSwimMethod(false, "JOIN", ni.String())
	if _, ok := s.members[ni.Address]; ok == false {
		s.disseminate(func(m *Member) {
			conn, err := m.Connect()
			defer conn.Close()
			if err != nil {
				return
			}
			c := NewMemberServiceClient(conn)
			ani := s.createAllNodeInfo()
			aniResponse, err := c.NotifyNode(context.Background(), ani)
			if err != nil {
				return
			}
			s.addMembers(aniResponse.Nodes)
		})

		m := s.buildMember(ni)
		s.addMember(m)
	}

	return s.createAllNodeInfo(), nil
}

func (s *SwimManager) disseminate(f func(*Member)) {
	for _, m := range s.members {
		go f(m)
	}
}

func (s *SwimManager) createAllNodeInfo() *AllNodeInfo {
	ani := &AllNodeInfo{
		Nodes: []*NodeInfo{},
	}
	for _, m := range s.members {
		ani.Nodes = append(ani.Nodes, s.buildNodeInfo(m))
	}

	return ani
}

func (s *SwimManager) NotifyNode(ctx context.Context, ani *AllNodeInfo) (*AllNodeInfo, error) {
	util.LogSwimMethod(false, "NotifyNode", ani.String())
	s.addMembers(ani.Nodes)

	return s.createAllNodeInfo(), nil
}

func (s *SwimManager) pingRandomMember() error {
	m := s.getRandomMember()
	if m == nil {
		return nil
	}
	return s.execPing(m)
}

func (s *SwimManager) execPing(m *Member) error {
	conn, err := m.Connect()
	defer conn.Close()
	if err != nil {
		go s.execRequirePing(m)
		return err
	}
	c := NewSwimServiceClient(conn)
	ni := s.buildNodeInfo(s.myInfo)

	util.LogSwimMethod(true, "Ping", ni.String())
	ctx, cancel := context.WithTimeout(context.Background(), s.pingInterval)
	defer cancel()
	ack, err := c.Ping(ctx, ni)
	if err != nil {
		go s.execRequirePing(m)
		return err
	}

	if ack.IsJoined == false {
		go func() {
			err := s.joinMember(m)
			if err != nil {
				s.deleteMember(m)
			}
		}()
	}
	if m.Status == SUSPECT {
		go s.execAlive(m)
	}
	return nil
}

func (s *SwimManager) buildNodeInfo(m *Member) *NodeInfo {
	return &NodeInfo{
		Address:           m.Address(),
		IncarnationNumber: int32(m.IncarnationNumber),
	}
}

func (s *SwimManager) getRandomMember() *Member {
	if len(s.members) == 0 {
		return nil
	}
	r := rand.Intn(len(s.members))
	i := 0
	for _, m := range s.members {
		if i == r {
			return m
		}
		i++
	}
	return nil
}

func (s *SwimManager) Ping(ctx context.Context, ni *NodeInfo) (*AckPing, error) {
	util.LogSwimMethod(false, "Ping", ni.String())
	if _, ok := s.members[ni.Address]; ok {
		util.LogSwimMethod(false, "Catch Ping", "true")
		return &AckPing{IsJoined: true}, nil
	} else {
		util.LogSwimMethod(false, "Catch Ping", "false")
		return &AckPing{IsJoined: false}, nil
	}
}

func (s *SwimManager) execRequirePing(m *Member) {
	util.Logf("Require Ping! %s\n", m)
	if len(s.members) == 1 {
		s.deleteMember(m)
		return
	}
	others := []*Member{}

	perm := rand.Perm(len(s.members))
	members := []*Member{}

	for _, member := range s.members {
		members = append(members, member)
	}

	mCount := 0
	for _, v := range perm {
		otherMember := members[v]
		if otherMember == m {
			continue
		}
		others = append(others, otherMember)
		mCount++

		if mCount > 4 {
			break
		}
	}

	var wg sync.WaitGroup
	alive := false
	for _, o := range others {
		wg.Add(1)
		go func() {
			aliveRes, err := s.execRequirePingToMember(o, m)
			if err == nil && aliveRes {
				alive = true
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if alive == false {
		if m.Status == ALIVE {
			m.Status = SUSPECT
			go s.execSuspect(m)
		} else {
			s.deleteMember(m)
			go s.execConfirm(m)
		}
	} else if m.Status == SUSPECT {
		m.Status = ALIVE
		go s.execAlive(m)
	}
}

func (s *SwimManager) execRequirePingToMember(m *Member, noResponseMember *Member) (bool, error) {
	conn, err := m.Connect()
	defer conn.Close()
	if err != nil {
		return false, err
	}
	c := NewSwimServiceClient(conn)

	util.LogSwimMethod(true, "RequirePing", noResponseMember.String())
	ctx, cancel := context.WithTimeout(context.Background(), PING_TIMEOUT*2*time.Second)
	defer cancel()
	ack, err := c.RequirePing(ctx, s.buildNodeInfo(noResponseMember))
	if err != nil {
		return false, err
	} else {
		return ack.Success, nil
	}
}

func (s *SwimManager) execSuspect(m *Member) {
	s.disseminate(func(target *Member) {
		conn, err := target.Connect()
		defer conn.Close()
		if err != nil {
			return
		}
		c := NewSwimServiceClient(conn)
		c.Suspect(context.Background(), s.buildNodeInfo(m))
	})
}

func (s *SwimManager) RequirePing(ctx context.Context, ni *NodeInfo) (*Result, error) {
	util.LogSwimMethod(false, "RequirePing", ni.String())
	m := s.buildMember(ni)

	conn, err := m.Connect()
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := NewSwimServiceClient(conn)

	timeout := make(chan bool)
	ackCh := make(chan *AckPing)
	errCh := make(chan error)
	go func() {
		time.Sleep(s.pingInterval)
		timeout <- true
	}()
	go func() {
		util.LogSwimMethod(true, "Ping", ni.String())
		// Not use WithTimeout to distinguish timeout from others
		ack, err := c.Ping(context.Background(), ni)
		if err != nil {
			errCh <- err
		} else {
			ackCh <- ack
		}
	}()

	select {
	case <-timeout:
		return &Result{Success: false}, nil
	case err := <-errCh:
		return nil, err
	case <-ackCh:
		return &Result{Success: true}, nil
	}
}

func (s *SwimManager) Suspect(ctx context.Context, ni *NodeInfo) (*Empty, error) {
	util.LogSwimMethod(false, "Suspect", ni.String())
	if m := s.getMember(ni); m != nil {
		m.Status = SUSPECT
	}
	return &Empty{}, nil
}

func (s *SwimManager) Alive(ctx context.Context, ni *NodeInfo) (*Empty, error) {
	util.LogSwimMethod(false, "Alive", ni.String())
	if m := s.getMember(ni); m != nil {
		m.Status = ALIVE
	}
	return &Empty{}, nil
}

func (s *SwimManager) execAlive(m *Member) {
	s.disseminate(func(target *Member) {
		conn, err := target.Connect()
		defer conn.Close()
		if err != nil {
			return
		}
		c := NewSwimServiceClient(conn)
		util.LogSwimMethod(true, "Alive", m.String())
		c.Alive(context.Background(), s.buildNodeInfo(m))
	})
}

func (s *SwimManager) Confirm(ctx context.Context, ni *NodeInfo) (*Empty, error) {
	util.LogSwimMethod(false, "Confirm", ni.String())
	if m := s.getMember(ni); m != nil {
		s.deleteMember(m)
	}
	return &Empty{}, nil
}

func (s *SwimManager) execConfirm(m *Member) {
	s.disseminate(func(target *Member) {
		conn, err := target.Connect()
		defer conn.Close()
		if err != nil {
			return
		}
		c := NewSwimServiceClient(conn)
		c.Confirm(context.Background(), s.buildNodeInfo(m))
	})
}
