package kademlia

import (
	"container/list"
	"math/big"
	"sort"
)

// A helper struct to sort peers by their distance to the local node
type peerDistance struct {
	p        Node
	distance []byte
}

// peerDistanceSorter implements sort.Interface to sort peers by xor distance
type peerDistanceSorter struct {
	peers  []peerDistance
	target []byte
}

func (pds *peerDistanceSorter) Len() int { return len(pds.peers) }
func (pds *peerDistanceSorter) Swap(a, b int) {
	pds.peers[a], pds.peers[b] = pds.peers[b], pds.peers[a]
}
func (pds *peerDistanceSorter) Less(a, b int) bool {
	first := &big.Int{}
	second := &big.Int{}

	first.SetBytes(pds.peers[a].distance)
	second.SetBytes(pds.peers[b].distance)

	return first.Cmp(second) == -1
}

// Append the peer.ID to the sorter's slice. It may no longer be sorted.
func (pds *peerDistanceSorter) appendPeer(p Node) {
	pds.peers = append(pds.peers, peerDistance{
		p:        p,
		distance: XOR(pds.target, p.HashedID),
	})
}

// Append the peer.ID values in the list to the sorter's slice. It may no longer be sorted.
func (pds *peerDistanceSorter) appendPeersFromList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		pds.appendPeer(e.Value.(Node))
	}
}

func (pds *peerDistanceSorter) sort() {
	sort.Sort(pds)
}

// Sort the given peers by their ascending distance from the target. A new slice is returned.
func SortClosestPeers(peers []Node, target []byte) []Node {
	sorter := peerDistanceSorter{
		peers:  make([]peerDistance, 0, len(peers)),
		target: target,
	}
	for _, p := range peers {
		sorter.appendPeer(p)
	}
	sorter.sort()
	out := make([]Node, 0, sorter.Len())
	for _, p := range sorter.peers {
		out = append(out, p.p)
	}
	return out
}
