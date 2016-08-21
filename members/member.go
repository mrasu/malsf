package members

import (
	"google.golang.org/grpc"
	"fmt"
)

type Member struct {
	Name string
	Addr grpc.Address
}

func NewMember(name string, addr string, port int) *Member {
	return &Member{
		Name: name,
		Addr: grpc.Address{
			Addr: fmt.Sprintf("%s:%d", addr, port),
		},
	}
}

func (m *Member) Address() string {
	return m.Addr.Addr
}