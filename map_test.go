package ps

import "testing"

func TestMapImmutable(t *testing.T) {
    // build a couple small maps
    world := NewMap().Set("hello", "world")
    kids := world.Set("hello", "kids")

    // both maps should still retain their data
    if v, _ := world.Lookup("hello"); v != "world" {
        t.Errorf("Set() modified the receiving map")
    }
    if size := world.Size(); size != 1 {
        t.Errorf("world size is not 1 : %d", size)
    }
    if v, _ := kids.Lookup("hello"); v != "kids" {
        t.Errorf("Set() did not modify the resulting map")
    }
    if size := kids.Size(); size != 1 {
        t.Errorf("kids size is not 1 : %d", size)
    }

    empty := kids.Delete("hello")
    if size := empty.Size(); size != 0 {
        t.Errorf("empty size is not 1 : %d", size)
    }
}

func BenchmarkMapSet(b *testing.B) {
    m := NewMap()
    for i := 0; i < b.N; i++ {
        m = m.Set("foo", i)
    }
}

func BenchmarkMapDelete(b *testing.B) {
    m := NewMap().Set("key", "value")
    for i := 0; i < b.N; i++ {
        m.Delete("key")
    }
}
