package server

import (
	"net"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"fmt"
	"github.com/mrasu/malsf/discover"
	"github.com/mrasu/malsf/structs"
	"github.com/mrasu/malsf/util"
)

type Server struct {
	grpcServer *grpc.Server
	serverAct structs.ServerAct
	*sender
}

func StartServer(port int, serverAct structs.ServerAct) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := &Server{
		sender: &sender{
			Name: serverAct.Name(),
			ServiceName: serverAct.Service(),
		},
		grpcServer: grpc.NewServer(),
		serverAct: serverAct,
	}
	structs.RegisterActionServiceServer(s.grpcServer, s)
	s.grpcServer.Serve(lis)
}

func(s *Server) Notify(ctx context.Context, action *structs.Action) (*structs.Reaction, error) {
	util.LogActionReceived(s.Name, action.NodeName, action.Id, fmt.Sprintf("Get (%s): %s", action.Type, action.Message))

	message, err := s.serverAct.Receive(ctx, action)
	if err != nil {
		return nil, err
	} else if message != nil {
		name := action.NodeName
		member, err := discover.NewNodeDiscoverer().GetMember(name, action.Service)
		if err != nil {
			panic(err)
		}
		s.send(member, message)
	}

	r := &structs.Reaction{
		Id: s.incrementId(),
		FromId: action.Id,
		Code: 0,
		Message: "Success",
	}

	util.LogReaction(s.Name, r.Id, action.NodeName, r.FromId, r.Message)
	return r, nil
}