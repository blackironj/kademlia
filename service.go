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
func (s *kademliaNet) FindNode(ctx context.Context, target *Target) (*Nodes, error) {
	hashedTargetID := ConvertPeerID(target.GetTargetId())

	senderID := target.Sender.GetId()
	senderIP := target.Sender.GetIp()
	senderKadPort := target.Sender.GetPort()

	sender := NewNode(senderID, senderIP, senderKadPort)

	s.table.Update(sender)

	nodes := s.table.NearestPeers(hashedTargetID, ParallelismAlpha)

	var neighbors []*NodeInfo

	for _, n := range nodes {
		neighbor := &NodeInfo{
			Id:   n.ID,
			Ip:   n.IP,
			Port: n.Port,
		}
		neighbors = append(neighbors, neighbor)
	}
	return &Nodes{
		Nodes: neighbors,
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
		TargetId: s.table.selfID,
		Sender: &NodeInfo{
			Id:   s.table.selfID,
			Ip:   s.table.selfIP,
			Port: s.table.selfPort,
		},
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

			for _, info := range res.GetNodes() {
				foundNode := NewNode(info.Id, info.Ip, info.Port)
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
