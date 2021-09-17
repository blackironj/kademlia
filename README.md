# Kademlia-like
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/blackironj/kademlia/master/LICENSE)
<p align="left">
  <a href="https://github.com/blackironj/kademlia/actions"><img alt="GitHub Actions status" src="https://github.com/actions/setup-go/workflows/build-test/badge.svg"></a>
</p>

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
	// if Options are empty. Set default value

	kadNet := kad.NewKademliaNet(routingTable)

    // if you do not want to bootrap. it's okay to skip this step
    // generate a bootstrap node using node kad.NewNode(id, ip, port)
	bootstrapNodes := []kad.Node{}
	kadNet.Bootstrap(bootstrapNodes)

	kadNet.Start()
}
```
