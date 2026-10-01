package lsm

import (
	"slices"
	"strings"
)

// entry is a single versioned record. A tombstone marks a deletion so that
// older values in lower SSTables are shadowed until compaction drops them.
type entry struct {
	key       string
	value     []byte
	tombstone bool
}

// Memtable is the in-memory write buffer. Keys must be kept in sorted order
// so that flushing to an SSTable is a simple sequential scan.
type Memtable struct {
	size    int // approximate bytes held, used to decide when to flush
	entries []entry
}

func NewMemtable() *Memtable {
	return &Memtable{
		entries: make([]entry, 0),
	}
}

// Put inserts or overwrites key.
func (m *Memtable) Put(key string, value []byte) {
	idx, found := slices.BinarySearchFunc(m.entries, key, func(e entry, k string) int {
		return strings.Compare(e.key, k)
	})
	if found {
		m.size += (len(value) - len(m.entries[idx].value))
		m.entries[idx].tombstone = false
		m.entries[idx].value = value
		return
	}
	m.size += len(value)
	m.entries = slices.Insert(m.entries, idx, entry{key: key, value: value, tombstone: false})
}

// Delete records a tombstone for key.
func (m *Memtable) Delete(key string) {
	idx, found := slices.BinarySearchFunc(m.entries, key, func(e entry, k string) int {
		return strings.Compare(e.key, k)
	})
	if found {
		m.size -= len(m.entries[idx].value)
		m.entries[idx].value = nil
		m.entries[idx].tombstone = true
		return
	}
	m.entries = slices.Insert(m.entries, idx, entry{key: key, value: nil, tombstone: true})
}

// Get returns (value, found, deleted). deleted=true means a tombstone was hit
// and the caller must NOT consult older SSTables.
func (m *Memtable) Get(key string) (value []byte, found bool, deleted bool) {
	idx, found := slices.BinarySearchFunc(m.entries, key, func(e entry, k string) int {
		return strings.Compare(e.key, k)
	})
	if found {
		return m.entries[idx].value, true, m.entries[idx].tombstone
	}
	return nil, false, false
}

// Size returns approximate memory usage in bytes.
func (m *Memtable) Size() int { return m.size }
