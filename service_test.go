package kademlia

import (
	"testing"
)

func TestKademliaNet_RefreshBuckets(t *testing.T) {
	rt := NewRoutingTable(20, selfID, selfIP, selfKadPort, selfServPort)
	nodes := genRandomNode(1000)

	for _, node := range nodes {
		rt.Update(node)
	}

	kadNet := NewKademliaNet(rt)
	kadNet.RefreshBuckets()

	if rt.Size() > 0 {
		t.Fatal("should not have peer")
	}
}

func TestKademliaNet_ReqFindNeighborsQuery(t *testing.T) {
	rtFirst := NewRoutingTable(20, selfID, selfIP, selfKadPort, selfServPort)
	nodes := genRandomNode(1000)

	for _, node := range nodes {
		rtFirst.Update(node)
	}

	kadNetFirst := NewKademliaNet(rtFirst)
	go kadNetFirst.Start(selfKadPort)

	firstNode := NewNode(selfID, selfIP, selfKadPort, selfServPort)
	rtTest := NewRoutingTable(20, selfID+"test", selfIP, selfKadPort+"0", selfServPort + "0")
	rtTest.Update(firstNode)

	kadNetTest := NewKademliaNet(rtTest)

	foundNodes := kadNetTest.ReqFindNeighborsQuery()

	if len(foundNodes) == 0 {
		t.Fatal("should have at least one node")
	}
}