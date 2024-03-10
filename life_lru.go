package life_lru

import (
	"container/heap"
	"sync"
	"time"
)

type nodeHeap[K, V comparable] []*lruNode[K, V]

func (h nodeHeap[K, V]) Len() int { return len(h) }
func (h nodeHeap[K, V]) Less(i, j int) bool {
	return h[i].timeout.Before(h[j].timeout)
}
func (h nodeHeap[K, V]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIdx = i
	h[j].heapIdx = j
}

func (h *nodeHeap[K, V]) Push(n any) {
	*h = append(*h, n.(*lruNode[K, V]))
}
func (h *nodeHeap[K, V]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type lruNode[K, V comparable] struct {
	k       K
	v       V
	timeout time.Time
	next    *lruNode[K, V]
	prev    *lruNode[K, V]
	heapIdx int
}

type LRUCache[K, V comparable] struct {
	mu        sync.Mutex
	data      map[K]*lruNode[K, V]
	capacity  uint32
	empty     V
	headNode  *lruNode[K, V]
	tailNode  *lruNode[K, V]
	timerHeap nodeHeap[K, V]
}

func NewLRUCache[K, V comparable](capacity uint32) *LRUCache[K, V] {
	cache := &LRUCache[K, V]{
		data:      make(map[K]*lruNode[K, V]),
		capacity:  capacity,
		timerHeap: make(nodeHeap[K, V], 0, capacity),
		headNode:  &lruNode[K, V]{},
		tailNode:  &lruNode[K, V]{},
	}
	cache.headNode.next = cache.tailNode
	cache.tailNode.prev = cache.headNode
	return cache
}

func (lc *LRUCache[K, V]) removeNode(node *lruNode[K, V]) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (lc *LRUCache[K, V]) addNode(node *lruNode[K, V]) {
	node.prev = lc.headNode
	node.next = lc.headNode.next
	lc.headNode.next.prev = node
	lc.headNode.next = node
}

func (lc *LRUCache[K, V]) Set(key K, val V, duration time.Duration) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	node, ok := lc.data[key]
	if ok {
		node.v = val
		lc.removeNode(node)
		lc.addNode(node)
		if !node.timeout.Equal(time.Now().Add(duration)) {
			node.timeout = time.Now().Add(duration)
			heap.Fix(&lc.timerHeap, node.heapIdx)
		}
		return
	}
	if len(lc.data) >= int(lc.capacity) {
		var del *lruNode[K, V]
		//If there is an expired key, remove the expired key;
		//otherwise, remove the least recently used key.
		if time.Now().After(lc.timerHeap[0].timeout) {
			del = heap.Pop(&lc.timerHeap).(*lruNode[K, V])
		} else {
			del = lc.tailNode.prev
			lc.removeNode(del)
			heap.Remove(&lc.timerHeap, del.heapIdx)
		}
		lc.removeNode(del)
		delete(lc.data, del.k)
	}

	//not exist, push to cache
	node = &lruNode[K, V]{
		k:       key,
		v:       val,
		timeout: time.Now().Add(duration),
	}
	lc.addNode(node)
	heap.Push(&lc.timerHeap, node)
	lc.data[key] = node
}

func (lc *LRUCache[K, V]) Get(key K) (V, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	node, ok := lc.data[key]
	if !ok {
		return lc.empty, false
	}
	lc.removeNode(node)
	if time.Now().After(node.timeout) {
		heap.Remove(&lc.timerHeap, node.heapIdx)
		delete(lc.data, node.k)
		return lc.empty, false
	}
	lc.addNode(node)
	return node.v, true
}
