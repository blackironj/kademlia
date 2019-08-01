package kademlia

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	ParalelismAlpha = 3
)

type kademliaServer struct {
	table *RoutingTable
}

func NewKademliaServer(routingTable *RoutingTable) *kademliaServer {
	ks := &kademliaServer{
		table : routingTable,
	}
	return ks
}

// `FIND NODE` RPC
func (s *kademliaServer) FindNode(ctx context.Context, target *Target) (*Neighbors, error) {
	hashedTargetID := ConvertPeerID(target.GetTargetId())

	senderID := target.GetSenderId()
	senderIP := target.GetSenderIp()
	senderKadPort := target.GetSenderKadPort()
	senderServPort := target.GetSenderServPort()

	sender := NewNode(senderID, senderIP, senderKadPort, senderServPort)

	s.table.Update(sender)

	nodes := s.table.NearestPeers(hashedTargetID, ParalelismAlpha)

	var neighbors []*NeighborInfo

	for _, n := range nodes {
		neighbor := &NeighborInfo{
			Id:                   n.ID,
			Ip:                   n.IP,
			KadPort:              n.KademliaPort,
			ServPort:             n.ServicePort,
		}
		neighbors = append(neighbors, neighbor)
	}
	return &Neighbors{
		Neighbors: neighbors,
	}, nil
}

func (s *kademliaServer) Start(kadPort string){
	lis, err := net.Listen("tcp", ":"+kadPort)
	if err != nil {
		log.Fatal(err)
	}

	rpcServer := grpc.NewServer()
	RegisterKademliaServiceServer(rpcServer, s)

	err = rpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}