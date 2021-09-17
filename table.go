package kademlia

import (
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
)

var ErrPeerRejectedNoCapacity = errors.New("peer rejected; insufficient capacity")

const (
	_defaultBucketSize = 20
	_defaultPort       = "50051"
)

// RoutingTable defines the routing table.
type RoutingTable struct {
	// ID of the local peer
	selfID       string
	hashedSelfID []byte

	selfIP   string
	selfPort string

	// Blanket lock, refine later for better performance
	tabLock sync.RWMutex

	// kBuckets define all the fingers to other nodes.
	Buckets    []*Bucket
	bucketsize int

	PeerRemoved func(string)
	PeerAdded   func(string)
}

//Options for initialize routing table
type Options struct {
	BucketSize  int
	ID          string
	IP          string
	Port        string
	PeerRemoved func(string)
	PeerAdded   func(string)
}

// NewRoutingTable creates a new routing table with a given bucketsize, local ID, and latency tolerance.
func NewRoutingTable(options *Options) *RoutingTable {
	if options.IP == "" {
		myIP, err := GetMyIP()
		if err != nil {
			log.Fatal(err)
		}
		options.IP = myIP
	}

	if options.Port == "" {
		options.Port = _defaultPort
	}

	if options.BucketSize == 0 {
		options.BucketSize = _defaultBucketSize
	}

	if options.ID == "" {
		options.ID = NewUUIDv4()
	}

	rt := &RoutingTable{
		Buckets:      []*Bucket{newBucket()},
		bucketsize:   options.BucketSize,
		selfID:       options.ID,
		selfIP:       options.IP,
		selfPort:     options.Port,
		hashedSelfID: ConvertPeerID(options.ID),

		PeerRemoved: func(string) {},
		PeerAdded:   func(string) {},
	}

	return rt
}

// Update adds or moves the given peer to the front of its respective bucket
func (rt *RoutingTable) Update(n Node) (err error) {
	cpl := CommonPrefixLen(n.HashedID, rt.hashedSelfID)

	rt.tabLock.Lock()
	defer rt.tabLock.Unlock()
	bucketID := cpl
	if bucketID >= len(rt.Buckets) {
		bucketID = len(rt.Buckets) - 1
	}

	bucket := rt.Buckets[bucketID]
	if bucket.Has(n) {
		// If the peer is already in the table, move it to the front.
		// This signifies that it it "more active" and the less active nodes
		// Will as a result tend towards the back of the list
		bucket.MoveToFront(n)
		return nil
	}

	// We have enough space in the bucket (whether spawned or grouped).
	if bucket.Len() < rt.bucketsize {
		n.makeConnection()
		bucket.PushFront(n)
		rt.PeerAdded(n.ID)
		return nil
	}

	if bucketID == len(rt.Buckets)-1 {
		// if the bucket is too large and this is the last bucket (i.e. wildcard), unfold it.
		rt.nextBucket()
		// the structure of the table has changed, so let's recheck if the peer now has a dedicated bucket.
		bucketID = cpl
		if bucketID >= len(rt.Buckets) {
			bucketID = len(rt.Buckets) - 1
		}
		bucket = rt.Buckets[bucketID]
		if bucket.Len() >= rt.bucketsize {
			// if after all the unfolding, we're unable to find room for this peer, scrap it.
			return ErrPeerRejectedNoCapacity
		}
		n.makeConnection()
		bucket.PushFront(n)
		rt.PeerAdded(n.ID)
		return nil
	}

	return ErrPeerRejectedNoCapacity
}

// Remove deletes a peer from the routing table. This is to be used
// when we are sure a node has disconnected completely.
func (rt *RoutingTable) Remove(n Node) {
	cpl := CommonPrefixLen(n.HashedID, rt.hashedSelfID)

	rt.tabLock.Lock()
	defer rt.tabLock.Unlock()

	bucketID := cpl
	if bucketID >= len(rt.Buckets) {
		bucketID = len(rt.Buckets) - 1
	}

	bucket := rt.Buckets[bucketID]
	if n.Conn != nil {
		n.Conn.Close()
	}
	if bucket.Remove(n) {
		rt.PeerRemoved(n.ID)
	}
}

func (rt *RoutingTable) RemoveDeadNodes() {
	for _, bucket := range rt.Buckets {
		bucket.RemoveDeadNodes()
	}
}

func (rt *RoutingTable) nextBucket() {
	// This is the last bucket, which allegedly is a mixed bag containing peers not belonging in dedicated (unfolded) buckets.
	// _allegedly_ is used here to denote that *all* peers in the last bucket might feasibly belong to another bucket.
	// This could happen if e.g. we've unfolded 4 buckets, and all peers in folded bucket 5 really belong in bucket 8.
	bucket := rt.Buckets[len(rt.Buckets)-1]
	newBucket := bucket.Split(len(rt.Buckets)-1, rt.hashedSelfID)
	rt.Buckets = append(rt.Buckets, newBucket)

	// The newly formed bucket still contains too many peers. We probably just unfolded a empty bucket.
	if newBucket.Len() >= rt.bucketsize {
		// Keep unfolding the table until the last bucket is not overflowing.
		rt.nextBucket()
	}
}

// Find a specific peer by ID or return nil
func (rt *RoutingTable) Find(id string) Node {
	srch := rt.NearestPeers(ConvertPeerID(id), 1)
	if len(srch) == 0 || srch[0].ID != id {
		return Node{}
	}
	return srch[0]
}

// NearestPeer returns a single peer that is nearest to the given ID
func (rt *RoutingTable) NearestPeer(hashedID []byte) Node {
	peers := rt.NearestPeers(hashedID, 1)
	if len(peers) > 0 {
		return peers[0]
	}

	return Node{}
}

// NearestPeers returns a list of the 'count' closest peers to the given ID
func (rt *RoutingTable) NearestPeers(hashedID []byte, count int) []Node {
	cpl := CommonPrefixLen(hashedID, rt.hashedSelfID)

	// It's assumed that this also protects the buckets.
	rt.tabLock.RLock()

	// Get bucket at cpl index or last bucket
	var bucket *Bucket
	if cpl >= len(rt.Buckets) {
		cpl = len(rt.Buckets) - 1
	}
	bucket = rt.Buckets[cpl]

	pds := peerDistanceSorter{
		peers:  make([]peerDistance, 0, 3*rt.bucketsize),
		target: hashedID,
	}
	pds.appendPeersFromList(bucket.list)
	if pds.Len() < count {
		// In the case of an unusual split, one bucket may be short or empty.
		// if this happens, search both surrounding buckets for nearby peers
		if cpl > 0 {
			pds.appendPeersFromList(rt.Buckets[cpl-1].list)
		}
		if cpl < len(rt.Buckets)-1 {
			pds.appendPeersFromList(rt.Buckets[cpl+1].list)
		}
	}
	rt.tabLock.RUnlock()

	// Sort by distance to local peer
	pds.sort()

	if count < pds.Len() {
		pds.peers = pds.peers[:count]
	}

	out := make([]Node, 0, pds.Len())
	for _, p := range pds.peers {
		out = append(out, p.p)
	}

	return out
}

// Size returns the total number of peers in the routing table
func (rt *RoutingTable) Size() int {
	var tot int
	rt.tabLock.RLock()
	for _, buck := range rt.Buckets {
		tot += buck.Len()
	}
	rt.tabLock.RUnlock()
	return tot
}

// ListPeers takes a RoutingTable and returns a list of all peers from all buckets in the table.
func (rt *RoutingTable) ListPeers() []Node {
	var peers []Node
	rt.tabLock.RLock()
	for _, buck := range rt.Buckets {
		peers = append(peers, buck.Nodes()...)
	}
	rt.tabLock.RUnlock()
	return peers
}
