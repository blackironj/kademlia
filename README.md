# Kademlia-like
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/blackironj/kademlia/master/LICENSE)
[![Build Status](https://travis-ci.org/blackironj/kademlia.svg?branch=master)](https://travis-ci.org/blackironj/kademlia)

This is a Go implementation of a `Kademlia-like` dht using [grpc](https://github.com/grpc/grpc-go)

_Currently, This project is experimental. So, I would not recommned using it in production enviroment._

## Quick start
```go
package main

import kad "github.com/blackironj/kademlia"

func main() {
	routingTable := kad.NewRoutingTable(
		&kad.Options{
			BucketSize: 10,
			ID:         "your unique id",
			IP:         "127.0.0.1",// your IP 
			Port:       "50051",    // your port number
		})
	// if you don't enter an ip, it set an ip automatically

	kadNet := kad.NewKademliaNet(routingTable)

    // if you do not want to bootrap. it's okay to skip this step
    // generate a bootstrap node using node kad.NewNode(id, ip, port)
	bootstrapNodes := []kad.Node{}
	kadNet.Bootstrap(bootstrapNodes)

	kadNet.Start()
}
```
