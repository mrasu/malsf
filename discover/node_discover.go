package discover

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/mrasu/malsf/members"
)

type NodeDiscoverer struct {
	client *api.Client
}

func NewNodeDiscoverer() *NodeDiscoverer {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	return &NodeDiscoverer{
		client: client,
	}
}

func (n *NodeDiscoverer) GetMember(name string, service string) (*members.Member, error) {
	node, _, err := n.client.Catalog().Node(name, nil)
	if err != nil {
		return nil, err
	}

	if service, ok := node.Services[service]; ok {
		m, err := members.NewMember(
			node.Node.Node,
			node.Node.Address,
			service.Port,
		)
		if err != nil {
			return nil, err
		}
		return m, nil
	} else {
		return nil, errors.New("service not found")
	}
}

func (n *NodeDiscoverer) GetMembersByTag(service string) ([]*members.Member, error) {
	service_members, _, err := n.client.Catalog().Service(service, "", nil)

	if err != nil {
		return nil, err
	}

	result := []*members.Member{}
	for _, member := range service_members {
		m, err := members.NewMember(
			member.Node,
			member.Address,
			member.ServicePort,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}
