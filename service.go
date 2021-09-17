package kademlia

import (
	"context"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	_defaultParallelismAlpha      = 3
	_defaultReqFindNodeDeadline   = 1 * time.Second
	_defaultRefreshBucketInterval = 5 * time.Minute
)

type KademliaOpt struct {
	ParallelismAlpha      int
	ReqFindNodeDeadline   time.Duration
	RefreshBucketInterval time.Duration
}

type kademliaNet struct {
	table *RoutingTable

	parallelismAlpha      int
	reqFindNodeDeadline   time.Duration
	refreshBucketInterval time.Duration
}

func NewKademliaNet(routingTable *RoutingTable, opt ...KademliaOpt) *kademliaNet {
	ks := &kademliaNet{
		table:                 routingTable,
		parallelismAlpha:      _defaultParallelismAlpha,
		reqFindNodeDeadline:   _defaultReqFindNodeDeadline,
		refreshBucketInterval: _defaultRefreshBucketInterval,
	}

	if len(opt) != 0 {
		if opt[0].ParallelismAlpha == 0 {
			ks.parallelismAlpha = opt[0].ParallelismAlpha
		}
		if opt[0].ReqFindNodeDeadline != 0 {
			ks.reqFindNodeDeadline = opt[0].ReqFindNodeDeadline
		}
		if opt[0].RefreshBucketInterval != 0 {
			ks.refreshBucketInterval = opt[0].RefreshBucketInterval
		}
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

	refreshTicker := time.NewTicker(s.refreshBucketInterval)

	go func() {
		for range refreshTicker.C {
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
	ctx, cancel := context.WithTimeout(context.Background(), s.reqFindNodeDeadline)

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
	nodes := s.table.NearestPeers(hashedTargetID, s.parallelismAlpha)

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
