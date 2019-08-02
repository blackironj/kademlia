package kademlia

import (
	"github.com/btcsuite/btcutil/base58"
	"math/rand"
	"strconv"
	"time"
)

const (
	selfID       = "self-id"
	selfIP       = "127.0.0.1"
	selfKadPort  = "50051"
	selfServPort = "50052"
	testID       = "test-id"
	testKadPort  = "50011"
	testServPort = "50012"
)

// Test basic features of the bucket struct
func genRandomNode(num int) []Node {
	seed := rand.NewSource(time.Now().UnixNano())
	ranGen := rand.New(seed)

	peers := make([]Node, num)

	for i := 0; i < num; i++ {
		var ip string
		for j := 0; j < 3; j++ {
			n := 1 + ranGen.Intn(255)
			ip += strconv.Itoa(n) + "."
		}
		lastNum := 1 + ranGen.Intn(255)
		ip += strconv.Itoa(lastNum)

		kadPortNum := 9000 + ranGen.Intn(60000)
		kadPort := strconv.Itoa(kadPortNum)

		servPortNum := 9000 + ranGen.Intn(60000)
		servPort := strconv.Itoa(servPortNum)

		peers[i] = NewNode(ip+":"+kadPort, ip, kadPort, servPort)
	}

	return peers
}

func genRandomID() string {
	t := time.Now().UnixNano()
	strconv.FormatInt(t, 10)
	return base58.Encode([]byte(strconv.FormatInt(t, 10)))
}
