package kademlia

import (
	"math/rand"
	"testing"
)

func TestBucket(t *testing.T) {
	b := newBucket()

	peers := genRandomNode(100)
	for i := 0; i < 100; i++ {
		b.PushFront(peers[i])
	}

	localID := selfID
	hashedLocalID := ConvertPeerID(localID)

	i := rand.Intn(len(peers))
	if !b.Has(peers[i]) {
		t.Errorf("Failed to find peer: %v", peers[i])
	}

	spl := b.Split(0, ConvertPeerID(localID))
	llist := b.list
	for e := llist.Front(); e != nil; e = e.Next() {
		p := ConvertPeerID(e.Value.(Node).ID)
		cpl := CommonPrefixLen(p, hashedLocalID)
		if cpl > 0 {
			t.Fatalf("Split failed. found id with cpl > 0 in 0 bucket")
		}
	}

	rlist := spl.list
	for e := rlist.Front(); e != nil; e = e.Next() {
		p := ConvertPeerID(e.Value.(Node).ID)
		cpl := CommonPrefixLen(p, hashedLocalID)
		if cpl == 0 {
			t.Fatalf("Split failed. found id with cpl == 0 in non 0 bucket")
		}
	}
}

func TestTableCallbacks(t *testing.T) {

	localID := selfID
	rt := NewRoutingTable(10, localID, selfIP, selfPort)

	peers := genRandomNode(100)

	pset := make(map[string]struct{})
	rt.PeerAdded = func(p string) {
		pset[p] = struct{}{}
	}
	rt.PeerRemoved = func(p string) {
		delete(pset, p)
	}

	rt.Update(peers[0])
	if _, ok := pset[peers[0].ID]; !ok {
		t.Fatal("should have this peer")
	}

	rt.Remove(peers[0])
	if _, ok := pset[peers[0].ID]; ok {
		t.Fatal("should not have this peer")
	}

	for _, p := range peers {
		rt.Update(p)
	}

	out := rt.ListPeers()
	for _, outp := range out {
		if _, ok := pset[outp.ID]; !ok {
			t.Fatal("should have peer in the peerset")
		}
		delete(pset, outp.ID)
	}

	if len(pset) > 0 {
		t.Fatal("have peers in peerset that were not in the table", len(pset))
	}
}

// Right now, this just makes sure that it doesnt hang or crash
func TestTableUpdate(t *testing.T) {
	localID := selfID
	rt := NewRoutingTable(10, localID, selfIP, selfPort)

	peers := genRandomNode(100)
	// Testing Update
	for i := 0; i < 10000; i++ {
		rt.Update(peers[rand.Intn(len(peers))])
	}

	for i := 0; i < 100; i++ {
		id := ConvertPeerID(genRandomID())
		ret := rt.NearestPeers(id, 5)
		if len(ret) == 0 {
			t.Fatal("Failed to find node near ID.")
		}
	}
}

func TestTableFind(t *testing.T) {
	localID := selfID

	rt := NewRoutingTable(10, localID, selfIP, selfPort)

	peers := genRandomNode(100)
	for i := 0; i < 5; i++ {
		rt.Update(peers[i])
	}

	t.Logf("Searching for peer: '%s'", peers[2].ID)
	found := rt.NearestPeer(peers[2].HashedID)
	if !(found.ID == peers[2].ID) {
		t.Fatalf("Failed to lookup known node...")
	}
}

func TestTableFindMultiple(t *testing.T) {
	localID := selfID

	rt := NewRoutingTable(20, localID, selfIP, selfPort)

	peers := genRandomNode(100)
	for i := 0; i < 25; i++ {
		rt.Update(peers[i])
	}

	t.Logf("Searching for peer: '%s'", peers[2].ID)
	found := rt.NearestPeers(peers[2].HashedID, 15)

	if len(found) != 15 {
		t.Fatalf("Got back different number of peers than we expected.")
	}
}

// Looks for race conditions in table operations. For a more 'certain'
// test, increase the loop counter from 1000 to a much higher number
// and set GOMAXPROCS above 1
func TestTableMultithreaded(t *testing.T) {
	localID := selfID

	tab := NewRoutingTable(20, localID, selfIP, selfPort)
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
	localID := selfID

	tab := NewRoutingTable(20, localID, selfIP, selfPort)

	num := b.N

	peers := genRandomNode(num)

	b.StartTimer()
	for i := 0; i < num; i++ {
		tab.Update(peers[i])
	}
}

func BenchmarkFinds(b *testing.B) {
	b.StopTimer()
	localID := selfID

	num := b.N

	tab := NewRoutingTable(20, localID, selfIP, selfPort)

	peers := genRandomNode(num)
	for i := 0; i < num; i++ {
		tab.Update(peers[i])
	}

	b.StartTimer()
	for i := 0; i < num; i++ {
		tab.Find(peers[i].ID)
	}
}
