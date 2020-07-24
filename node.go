package kademlia

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const MaxTryConnCount = 3

type Node struct {
	ID       string
	HashedID []byte

	IP   string
	Port string

	Conn *grpc.ClientConn

	FailedTryConnCount int32
}

func NewNode(id string, ip string, port string) Node {
	n := Node{
		ID:                 id,
		IP:                 ip,
		Port:               port,
		FailedTryConnCount: 0,
	}
	n.HashedID = ConvertPeerID(n.ID)

	return n
}

func (n *Node) makeConnection() {
	if n.Conn != nil {
		return
	}

	var err error
	n.Conn, err = grpc.Dial(n.IP+":"+n.Port, grpc.WithInsecure())

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
