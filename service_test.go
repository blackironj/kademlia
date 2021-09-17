package kademlia

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKademliaNet_RefreshBuckets(t *testing.T) {

	rt := NewRoutingTable(
		&Options{
			BucketSize: 20,
		})
	nodes := genRandomNode(1000)

	for _, node := range nodes {
		rt.Update(node)
	}

	kadNet := NewKademliaNet(rt)
	kadNet.RefreshBuckets()

	assert.Equal(t, 0, rt.Size())
}

func TestKademliaNet_ReqFindNeighborsQuery(t *testing.T) {
	rtFirst := NewRoutingTable(
		&Options{
			BucketSize: 20,
		})
	nodes := genRandomNode(100)

	for _, node := range nodes {
		rtFirst.Update(node)
	}

	kadNetFirst := NewKademliaNet(rtFirst)
	go kadNetFirst.Start()

	time.Sleep(time.Second * 2)

	testID := NewUUIDv4()
	testIP := "127.0.0.1"
	testPort := "50051"

	firstNode := NewNode(testID, testIP, testPort)
	rtTest := NewRoutingTable(
		&Options{
			BucketSize: 20,
			ID:         testID,
			IP:         testIP,
			Port:       testPort,
		})
	rtTest.Update(firstNode)

	kadNetTest := NewKademliaNet(rtTest)

	foundNodes := kadNetTest.ReqFindNodesFromRandom(kadNetTest.table.selfID)

	assert.Greater(t, len(foundNodes), 0)
}
