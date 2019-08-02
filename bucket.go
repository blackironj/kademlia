package kademlia

import (
	"bytes"
	"container/list"
	"sync"
)

// Bucket holds a list of peers.
type Bucket struct {
	lk   sync.RWMutex
	list *list.List
}

func newBucket() *Bucket {
	b := new(Bucket)
	b.list = list.New()
	return b
}

func (b *Bucket) Nodes() []Node {
	b.lk.RLock()
	defer b.lk.RUnlock()
	ps := make([]Node, 0, b.list.Len())
	for e := b.list.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		ps = append(ps, node)
	}
	return ps
}

func (b *Bucket) Has(n Node) bool {
	b.lk.RLock()
	defer b.lk.RUnlock()
	for e := b.list.Front(); e != nil; e = e.Next() {
		if bytes.Equal(e.Value.(Node).HashedID, n.HashedID) {
			return true
		}
	}
	return false
}

func (b *Bucket) Remove(n Node) bool {
	b.lk.Lock()
	defer b.lk.Unlock()
	for e := b.list.Front(); e != nil; e = e.Next() {
		if bytes.Equal(e.Value.(Node).HashedID, n.HashedID) {
			b.list.Remove(e)
			return true
		}
	}
	return false
}

func (b *Bucket) RemoveDeadNodes() {
	var next *list.Element
	for e := b.list.Front(); e != nil; e = next {
		node := e.Value.(Node)
		if !node.IsAlive() {
			next = e.Next()
			node.Conn.Close()
			b.list.Remove(e)
		} else {
			next = e.Next()
		}
	}
}

func (b *Bucket) MoveToFront(n Node) {
	b.lk.Lock()
	defer b.lk.Unlock()
	for e := b.list.Front(); e != nil; e = e.Next() {
		if bytes.Equal(e.Value.(Node).HashedID, n.HashedID) {
			b.list.MoveToFront(e)
		}
	}
}

func (b *Bucket) PushFront(n Node) {
	b.lk.Lock()
	b.list.PushFront(n)
	b.lk.Unlock()
}

func (b *Bucket) PopBack() Node {
	b.lk.Lock()
	defer b.lk.Unlock()
	last := b.list.Back()
	b.list.Remove(last)
	return last.Value.(Node)
}

func (b *Bucket) Len() int {
	b.lk.RLock()
	defer b.lk.RUnlock()
	return b.list.Len()
}

// Split splits a buckets peers into two buckets, the methods receiver will have
// peers with CPL equal to cpl, the returned bucket will have peers with CPL
// greater than cpl (returned bucket has closer peers)
func (b *Bucket) Split(cpl int, target []byte) *Bucket {
	b.lk.Lock()
	defer b.lk.Unlock()

	out := list.New()
	newbuck := newBucket()
	newbuck.list = out
	e := b.list.Front()
	for e != nil {
		peerID := e.Value.(Node).HashedID
		peerCPL := CommonPrefixLen(peerID, target)
		if peerCPL > cpl {
			cur := e
			out.PushBack(e.Value)
			e = e.Next()
			b.list.Remove(cur)
			continue
		}
		e = e.Next()
	}
	return newbuck
}
