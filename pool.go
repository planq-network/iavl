package iavl

import (
	"math"
	"sync"
)

type NodePool struct {
	syncPool      *sync.Pool
	hashBytesPool *sync.Pool
	keyBytesPool  *sync.Pool

	free     chan int
	nodes    []Node
	poolSize uint64

	poolId uint64
}

func NewNodePool() *NodePool {
	np := &NodePool{
		syncPool: &sync.Pool{
			New: func() interface{} {
				return &Node{}
			},
		},
		hashBytesPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 32)
			},
		},
		keyBytesPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0)
			},
		},
		free: make(chan int, 1000),
	}
	return np
}

func (np *NodePool) Get() *Node {
	if np.poolId == math.MaxUint64 {
		np.poolId = 1
	} else {
		np.poolId++
	}
	n := np.syncPool.Get().(*Node)
	n.poolId = np.poolId
	return n
}

func (np *NodePool) Put(node *Node) {
	np.resetNode(node)
	node.poolId = 0
	np.syncPool.Put(node)
}

func (np *NodePool) resetNode(node *Node) {
	node.leftNodeKey = emptyNodeKey
	node.rightNodeKey = emptyNodeKey
	node.rightNode = nil
	node.leftNode = nil
	node.nodeKey = emptyNodeKey
	node.hash = nil
	node.key = nil
	node.value = nil
	node.subtreeHeight = 0
	node.size = 0
	node.dirty = false
}
