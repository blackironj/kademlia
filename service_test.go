package kademlia

import (
	"testing"
	"time"
)

func TestKademliaNet_RefreshBuckets(t *testing.T) {

	rt := NewRoutingTable(
		&Options{
			BucketSize: 20,
			ID:         myID,
			IP:         myIP,
			Port:       myPort,
		})
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
	rtFirst := NewRoutingTable(
		&Options{
			BucketSize: 20,
			ID:         myID,
			IP:         myIP,
			Port:       myPort,
		})
	nodes := genRandomNode(100)

	for _, node := range nodes {
		rtFirst.Update(node)
	}

	kadNetFirst := NewKademliaNet(rtFirst)
	go kadNetFirst.Start(myPort)

	time.Sleep(time.Second * 2)

	firstNode := NewNode(myID, myIP, myPort)
	rtTest := NewRoutingTable(
		&Options{
			BucketSize: 20,
			ID:         testID,
			IP:         myIP,
			Port:       testPort,
		})
	rtTest.Update(firstNode)

	kadNetTest := NewKademliaNet(rtTest)

	foundNodes := kadNetTest.ReqFindNodesFromRandom(kadNetTest.table.selfID)

	if len(foundNodes) == 0 {
		t.Fatal("should have at least one node")
	}
}
