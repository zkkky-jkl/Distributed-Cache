package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps byte to uint32
type Hash func(data []byte) uint32

// Map contains all hashed keys
type Map struct {
	hash     Hash
	replicas int   // 添加一个节点 增加的虚拟节点倍数
	keys     []int // sorted
	hashMap  map[int]string
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	// cal hashVal of key
	hash := int(m.hash([]byte(key)))
	// binary search for the first virtual node
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// idx = len(m.keys) if hash > all keys
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
