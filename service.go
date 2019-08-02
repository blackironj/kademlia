package kademlia

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	ParallelismAlpha      = 3
	ReqFindNodeDeadline   = 1 * time.Second
	RefreshBucketInterval = 5 * time.Minute
)

type kademliaNet struct {
	table               *RoutingTable
	bucketRefreshTicker *time.Ticker
}

func NewKademliaNet(routingTable *RoutingTable) *kademliaNet {
	ks := &kademliaNet{
		table:               routingTable,
		bucketRefreshTicker: time.NewTicker(RefreshBucketInterval),
	}
	return ks
}

// `FIND NODE` RPC
func (s *kademliaNet) FindNode(ctx context.Context, target *Target) (*Neighbors, error) {
	hashedTargetID := ConvertPeerID(target.GetTargetId())

	senderID := target.GetSenderId()
	senderIP := target.GetSenderIp()
	senderKadPort := target.GetSenderKadPort()
	senderServPort := target.GetSenderServPort()

	sender := NewNode(senderID, senderIP, senderKadPort, senderServPort)

	s.table.Update(sender)

	nodes := s.table.NearestPeers(hashedTargetID, ParallelismAlpha)

	var neighbors []*NeighborInfo

	for _, n := range nodes {
		neighbor := &NeighborInfo{
			Id:       n.ID,
			Ip:       n.IP,
			KadPort:  n.KademliaPort,
			ServPort: n.ServicePort,
		}
		neighbors = append(neighbors, neighbor)
	}
	return &Neighbors{
		Neighbors: neighbors,
	}, nil
}

func (s *kademliaNet) Start(kadPort string) {
	lis, err := net.Listen("tcp", ":"+kadPort)
	if err != nil {
		log.Fatal(err)
	}

	rpcServer := grpc.NewServer()
	RegisterKademliaServiceServer(rpcServer, s)

	go func() {
		for range s.bucketRefreshTicker.C {
			s.RefreshBuckets()
		}
	}()

	err = rpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *kademliaNet) ReqFindNeighborsQuery() []Node {
	var nodes []Node

	target := &Target{
		TargetId:       s.table.selfID,
		SenderId:       s.table.selfID,
		SenderIp:       s.table.selfIP,
		SenderKadPort:  s.table.selfKadPort,
		SenderServPort: s.table.selfServPort,
	}

	for _, bucket := range s.table.Buckets {
		recentlySeenNode := bucket.list.Front()
		if recentlySeenNode != nil {
			client := NewKademliaServiceClient(recentlySeenNode.Value.(Node).Conn)
			ctx, cancel := context.WithTimeout(context.Background(), ReqFindNodeDeadline)

			res, err := client.FindNode(ctx, target)
			if err != nil {
				log.Fatal(err)
			}

			for _, info := range res.GetNeighbors() {
				foundNode := NewNode(info.Id, info.Ip, info.KadPort, info.ServPort)
				nodes = append(nodes, foundNode)
			}
			cancel()
		}
	}

	return nodes
}

func (s *kademliaNet) RefreshBuckets() {
	s.table.RemoveDeadNodes()
	foundNode := s.ReqFindNeighborsQuery()

	for _, n := range foundNode {
		s.table.Update(n)
	}
}
