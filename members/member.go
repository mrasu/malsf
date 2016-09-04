package members

import (
	"google.golang.org/grpc"
	"fmt"
)

type Status int

const (
	ALIVE Status = iota
	SUSPECT
)

type Member struct {
	Name string
	Addr grpc.Address
	IncarnationNumber int
	Status Status
}

func NewMember(name string, addr string, in int) (*Member, error) {
	fmt.Printf("||||%s\n", addr)
	return &Member{
		Name: name,
		Addr: grpc.Address{
			Addr: addr,
		},
		IncarnationNumber: in,
		Status: ALIVE,
	}, nil
}

func (m *Member) Address() string {
	return m.Addr.Addr
}

func (m *Member) Connect() (*grpc.ClientConn, error) {
	return grpc.Dial(m.Address(), grpc.WithInsecure())
}

func (m *Member) String() string {
	return fmt.Sprintf("Name(%s), Addr(%s), IncarnationNumber(%d), Status(%s)", m.Name, m.Addr, m.IncarnationNumber, m.Status)
}