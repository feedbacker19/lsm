package lsm

import "testing"

func TestPutGet(t *testing.T) {
	m := NewMemtable()
	for _, k := range []string{"mango", "apple", "zebra", "cherry"} {
		m.Put(k, []byte("v-"+k))
	}
	for _, k := range []string{"mango", "apple", "zebra", "cherry"} {
		v, found, deleted := m.Get(k)
		if !found || deleted || string(v) != "v-"+k {
			t.Errorf("Get(%q) = %q, found=%v, deleted=%v", k, v, found, deleted)
		}
	}
	if _, found, _ := m.Get("nope"); found {
		t.Error("Get(nope) should not be found")
	}
}

func TestOverwriteSize(t *testing.T) {
	m := NewMemtable()
	m.Put("a", []byte("12345"))
	m.Put("a", []byte("12"))
	if m.Size() != 2 {
		t.Errorf("size = %d, want 2", m.Size())
	}
}

func TestDelete(t *testing.T) {
	m := NewMemtable()
	m.Put("a", []byte("1"))
	m.Put("b", []byte("2"))
	m.Put("c", []byte("3"))
	m.Delete("b")
	if _, found, deleted := m.Get("b"); !found || !deleted {
		t.Errorf("Get(b) after delete: found=%v deleted=%v, want true,true", found, deleted)
	}
	// Tombstone for a key never written here must still shadow older data.
	m.Delete("zzz")
	if _, found, deleted := m.Get("zzz"); !found || !deleted {
		t.Errorf("Get(zzz): found=%v deleted=%v, want true,true", found, deleted)
	}
	// Put after Delete resurrects the key.
	m.Put("b", []byte("new"))
	if v, _, deleted := m.Get("b"); deleted || string(v) != "new" {
		t.Errorf("Get(b) after re-Put = %q deleted=%v", v, deleted)
	}
}
