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
	serverAct structs.ReceiverAct
	*sender
}

func StartServer(port int, serverAct structs.ReceiverAct, mch chan(*structs.Message)) {
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
	go func() {
		s.ListenMessage(mch)
	}()
	structs.RegisterActionServiceServer(s.grpcServer, s)
	s.grpcServer.Serve(lis)
}

func(s *Server) Notify(ctx context.Context, action *structs.Action) (*structs.Reaction, error) {
	util.LogActionReceived(s.Name, action.NodeName, action.Id, fmt.Sprintf("Get (%s): %s", action.Type, action.Message))

	m, err := s.serverAct.Receive(ctx, action)
	fmt.Println("Notify#sendAsync")
	if err != nil {
		return nil, err
	} else if m != nil {
		name := action.NodeName
		member, err := discover.NewNodeDiscoverer().GetMember(name, action.Service)
		if err != nil {
			panic(err)
		}
		_, err = s.sendAsync(member, m)
		if err != nil {
			panic(err)
		}
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

func(s *Server) ListenMessage(mch chan(*structs.Message)) error {
	for {
		m := <- mch
		n := discover.NewNodeDiscoverer()
		for _, service := range m.ToServices {
			members, err := n.GetMembersByTag(service)
			if err != nil {
				return err
			}

			for _, member := range members {
				if _, err := s.sendAsync(member, m); err != nil {
					return err
				}
			}
		}
	}
	return nil
}