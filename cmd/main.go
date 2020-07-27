package main

import (
	kad "github.com/blackironj/kademlia"
)

func main() {
	routingTable := kad.NewRoutingTable(
		&kad.Options{
			BucketSize: 10,
			ID:         "your unique id",
			IP:         "127.0.0.1", // your IP
			Port:       "50051",     // your port number
		})

	kadNet := kad.NewKademliaNet(routingTable)

	bootstrapNodes := []kad.Node{}
	kadNet.Bootstrap(bootstrapNodes)

	kadNet.Start()
}
