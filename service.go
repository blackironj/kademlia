package kademlia

import (
	"context"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
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

func (s *kademliaNet) Bootstrap(bootstrapNodes []Node) {
	for _, bootNode := range bootstrapNodes {
		if err := s.table.Update(bootNode); err != nil {
			log.Debug(err)
			continue
		}

		neighborNodes := s.ReqFindNodesFromSpecific(bootNode, s.table.selfID)
		for _, neighbor := range neighborNodes {
			if err := s.table.Update(neighbor); err != nil {
				log.Debug(err)
				continue
			}
		}
	}
}

func (s *kademliaNet) Start() {
	lis, err := net.Listen("tcp", ":"+s.table.selfPort)
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

func (s *kademliaNet) ReqFindNodesFromSpecific(dest Node, targetID string) []Node {
	target := s.genTargetMsg(targetID)

	return s.reqFindNodes(dest, target)
}

func (s *kademliaNet) ReqFindNodesFromRandom(targetID string) []Node {
	var nodes []Node

	target := s.genTargetMsg(targetID)

	for _, bucket := range s.table.Buckets {
		elem := bucket.list.Front()
		if elem != nil {
			recentlySeenNode := elem.Value.(Node)
			foundNodes := s.reqFindNodes(recentlySeenNode, target)
			nodes = append(nodes, foundNodes...)
		}
	}
	return nodes
}

func (s *kademliaNet) genTargetMsg(targetID string) *Target {
	return &Target{
		TargetId: targetID,
		Sender: &NodeInfo{
			Id:   s.table.selfID,
			Ip:   s.table.selfIP,
			Port: s.table.selfPort,
		},
	}
}

func (s *kademliaNet) reqFindNodes(dest Node, target *Target) []Node {
	dest.makeConnection()
	client := NewKademliaServiceClient(dest.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), ReqFindNodeDeadline)

	nodes := make([]Node, 0, 10)

	res, err := client.FindNode(ctx, target)
	if err != nil {
		log.Debug(err)
		cancel()
		return nodes
	}

	for _, info := range res.GetNodes() {
		if info.Id == s.table.selfID {
			continue
		}
		foundNode := NewNode(info.Id, info.Ip, info.Port)
		nodes = append(nodes, foundNode)
	}
	cancel()
	return nodes
}

func (s *kademliaNet) RefreshBuckets() {
	s.table.RemoveDeadNodes()
	foundNode := s.ReqFindNodesFromRandom(s.table.selfID)

	for _, n := range foundNode {
		s.table.Update(n)
	}
}

func (s *kademliaNet) FindNode(ctx context.Context, target *Target) (*Nodes, error) {
	hashedTargetID := ConvertPeerID(target.GetTargetId())

	senderID := target.Sender.GetId()
	senderIP := target.Sender.GetIp()
	senderKadPort := target.Sender.GetPort()

	sender := NewNode(senderID, senderIP, senderKadPort)

	if err := s.table.Update(sender); err != nil {
		log.Debug(err)
	}
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
