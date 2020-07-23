package kademlia

import (
	"testing"
	"time"
)

func TestKademliaNet_RefreshBuckets(t *testing.T) {
	rt := NewRoutingTable(20, selfID, selfIP, selfPort)
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
	rtFirst := NewRoutingTable(20, selfID, selfIP, selfPort)
	nodes := genRandomNode(100)

	for _, node := range nodes {
		rtFirst.Update(node)
	}

	kadNetFirst := NewKademliaNet(rtFirst)
	go kadNetFirst.Start(selfPort)

	time.Sleep(time.Second * 2)

	firstNode := NewNode(selfID, selfIP, selfPort)
	rtTest := NewRoutingTable(20, testID, selfIP, testPort)
	rtTest.Update(firstNode)

	kadNetTest := NewKademliaNet(rtTest)

	foundNodes := kadNetTest.ReqFindNeighborsQuery(kadNetTest.table.selfID)

	if len(foundNodes) == 0 {
		t.Fatal("should have at least one node")
	}
}
