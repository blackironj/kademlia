package kademlia

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	MaxTryConnCount            = 3
	NodeFailedConnsBeforeStale = 2
)

type Node struct {
	ID       string
	HashedID []byte

	IP           string
	KademliaPort string
	ServicePort  string

	Conn *grpc.ClientConn

	FailedReqCount     int32
	FailedTryConnCount int32
}

func NewNode(id string, ip string, kadPort string, servPort string) Node {
	n := Node{
		ID:                 id,
		IP:                 ip,
		KademliaPort:       kadPort,
		ServicePort:        servPort,
		FailedReqCount:     0,
		FailedTryConnCount: 0,
	}
	n.HashedID = ConvertPeerID(n.ID)
	return n
}

func (n *Node) makeConnection() {
	var err error
	n.Conn, err = grpc.Dial(n.IP+":"+n.KademliaPort, grpc.WithInsecure())

	if err == nil {
		n.Conn.GetState()
	}
}

func (n *Node) IsAlive() bool {
	if MaxTryConnCount <= n.FailedTryConnCount {
		return false
	}

	connState := n.Conn.GetState()
	if connState == connectivity.Ready {
		n.FailedTryConnCount = 0
		return true
	} else if connState == connectivity.Connecting {
		n.FailedTryConnCount++
	}
	return false
}

func (n *Node) IsStale() bool {
	return NodeFailedConnsBeforeStale == n.FailedReqCount
}
