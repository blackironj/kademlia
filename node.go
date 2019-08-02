package kademlia

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	MaxTryConnCount = 3
)

type Node struct {
	ID       string
	HashedID []byte

	IP           string
	KademliaPort string
	ServicePort  string

	Conn *grpc.ClientConn

	FailedTryConnCount int32
}

func NewNode(id string, ip string, kadPort string, servPort string) Node {
	n := Node{
		ID:                 id,
		IP:                 ip,
		KademliaPort:       kadPort,
		ServicePort:        servPort,
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
	if n.Conn == nil {
		return false
	}

	if MaxTryConnCount <= n.FailedTryConnCount {
		n.Conn.Close()
		return false
	}

	connState := n.Conn.GetState()
	if connState == connectivity.Ready {
		n.FailedTryConnCount = 0
		return true
	} else if connState == connectivity.Connecting {
		n.FailedTryConnCount++
		return n.IsAlive()
	}
	n.Conn.Close()
	return false
}
