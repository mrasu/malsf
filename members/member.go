package members

import (
	"fmt"
	"google.golang.org/grpc"
)

type Status int

const (
	ALIVE Status = iota
	SUSPECT
)

type Member struct {
	Name              string
	Addr              grpc.Address
	IncarnationNumber int
	Status            Status
}

func NewMember(name string, addr string, in int) *Member {
	return &Member{
		Name: name,
		Addr: grpc.Address{
			Addr: addr,
		},
		IncarnationNumber: in,
		Status:            ALIVE,
	}
}

func (m *Member) Address() string {
	return m.Addr.Addr
}

func (m *Member) Connect() (*grpc.ClientConn, error) {
	return grpc.Dial(m.Address(), grpc.WithInsecure())
}

func (m *Member) String() string {
	return fmt.Sprintf("Name(%s), Addr(%s), IncarnationNumber(%d), Status(%d)", m.Name, m.Addr.Addr, m.IncarnationNumber, m.Status)
}
