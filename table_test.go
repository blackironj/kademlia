package kademlia

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {
	b := newBucket()

	testPeers := genRandomNode(100)
	for i := 0; i < 100; i++ {
		b.PushFront(testPeers[i])
	}

	testID := NewUUIDv4()
	hashedTestID := ConvertPeerID(testID)

	i := rand.Intn(len(testPeers))
	assert.True(t, b.Has(testPeers[i]))

	spl := b.Split(0, ConvertPeerID(testID))
	llist := b.list
	for e := llist.Front(); e != nil; e = e.Next() {
		p := ConvertPeerID(e.Value.(Node).ID)
		cpl := CommonPrefixLen(p, hashedTestID)
		assert.Equal(t, 0, cpl)
	}

	rlist := spl.list
	for e := rlist.Front(); e != nil; e = e.Next() {
		p := ConvertPeerID(e.Value.(Node).ID)
		cpl := CommonPrefixLen(p, hashedTestID)
		assert.NotEqual(t, 0, cpl)
	}
}

func TestTableCallbacks(t *testing.T) {
	rt := NewRoutingTable(&Options{})

	peers := genRandomNode(100)

	pset := make(map[string]struct{})
	rt.PeerAdded = func(p string) {
		pset[p] = struct{}{}
	}
	rt.PeerRemoved = func(p string) {
		delete(pset, p)
	}

	rt.Update(peers[0])
	_, ok := pset[peers[0].ID]
	assert.True(t, ok)

	rt.Remove(peers[0])
	_, ok = pset[peers[0].ID]
	assert.False(t, ok)

	for _, p := range peers {
		rt.Update(p)
	}

	out := rt.ListPeers()
	for _, outp := range out {
		_, ok := pset[outp.ID]
		assert.True(t, ok)

		delete(pset, outp.ID)
	}

	assert.Equal(t, 0, len(pset))
}

// Right now, this just makes sure that it doesnt hang or crash
func TestTableUpdate(t *testing.T) {
	rt := NewRoutingTable(&Options{
		BucketSize: 10,
	})

	peers := genRandomNode(100)
	// Testing Update
	for i := 0; i < 10000; i++ {
		rt.Update(peers[rand.Intn(len(peers))])
	}

	for i := 0; i < 100; i++ {
		id := ConvertPeerID(genRandomID())
		ret := rt.NearestPeers(id, 5)

		assert.NotEqual(t, len(ret), 0)
	}
}

func TestTableFind(t *testing.T) {
	rt := NewRoutingTable(&Options{
		BucketSize: 10,
	})

	peers := genRandomNode(100)
	for i := 0; i < 5; i++ {
		rt.Update(peers[i])
	}

	t.Logf("Searching for peer: '%s'", peers[2].ID)
	found := rt.NearestPeer(peers[2].HashedID)

	assert.Equal(t, found.ID, peers[2].ID)
}

func TestTableFindMultiple(t *testing.T) {
	rt := NewRoutingTable(&Options{
		BucketSize: 20,
	})

	peers := genRandomNode(100)
	for i := 0; i < 25; i++ {
		rt.Update(peers[i])
	}

	t.Logf("Searching for peer: '%s'", peers[2].ID)
	found := rt.NearestPeers(peers[2].HashedID, 15)

	assert.Equal(t, 15, len(found))
}

// Looks for race conditions in table operations. For a more 'certain'
// test, increase the loop counter from 1000 to a much higher number
// and set GOMAXPROCS above 1
func TestTableMultithreaded(t *testing.T) {
	tab := NewRoutingTable(&Options{
		BucketSize: 20,
	})
	peers := genRandomNode(500)

	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			n := rand.Intn(len(peers))
			tab.Update(peers[n])
		}
		done <- struct{}{}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			n := rand.Intn(len(peers))
			tab.Update(peers[n])
		}
		done <- struct{}{}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			n := rand.Intn(len(peers))
			tab.Find(peers[n].ID)
		}
		done <- struct{}{}
	}()
	<-done
	<-done
	<-done
}

func BenchmarkUpdates(b *testing.B) {
	b.StopTimer()

	tab := NewRoutingTable(&Options{
		BucketSize: 20,
	})

	num := b.N

	peers := genRandomNode(num)

	b.StartTimer()
	for i := 0; i < num; i++ {
		tab.Update(peers[i])
	}
}

func BenchmarkFinds(b *testing.B) {
	b.StopTimer()
	num := b.N

	tab := NewRoutingTable(&Options{
		BucketSize: 20,
	})

	peers := genRandomNode(num)
	for i := 0; i < num; i++ {
		tab.Update(peers[i])
	}

	b.StartTimer()
	for i := 0; i < num; i++ {
		tab.Find(peers[i].ID)
	}
}
