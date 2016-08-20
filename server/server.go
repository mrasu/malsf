package server

import (
	"github.com/mrasu/malsf/members"
	"net"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"fmt"
	"github.com/mrasu/malsf/discover"
)

type Server struct {
	grpcServer *grpc.Server
}

func StartServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := &Server{
		grpcServer: grpc.NewServer(),
	}
	members.RegisterActionServiceServer(s.grpcServer, s)
	s.grpcServer.Serve(lis)
}

func(s *Server) Notify(ctx context.Context, action *members.Action) (*members.Reaction, error) {
	fmt.Printf("%s(%s), Message: %s\n", action.NodeName, action.Type, action.Message)
	if action.Type == "Memory" {
		fmt.Println("Memory is leeking!!")
		name := action.NodeName
		member, err := discover.NewNodeDiscoverer().GetMember(name, "client")
		if err != nil {
			panic(err)
		}
		member.Send("kill process!", "Kill Process")
	}

	return &members.Reaction{
		Id: action.Id + 1000,
		FromId: action.Id,
		Code: 0,
		Message: "Success",
	}, nil
}