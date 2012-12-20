package ps

import "testing"
import "sort"

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

    // both maps have the right keys
    if keys := world.Keys(); len(keys) != 1 || keys[0] != "hello" {
        t.Errorf("world has the wrong keys: %#v", keys)
    }
    if keys := kids.Keys(); len(keys) != 1 || keys[0] != "hello" {
        t.Errorf("kids has the wrong keys: %#v", keys)
    }

    // test deletion
    empty := kids.Delete("hello")
    if size := empty.Size(); size != 0 {
        t.Errorf("empty size is not 1 : %d", size)
    }
    if keys := empty.Keys(); len(keys) != 0 {
        t.Errorf("empty has the wrong keys: %#v", keys)
    }
}

func TestMapMultipleKeys(t *testing.T) {
    // map with multiple keys each with pointer values
    one := 1
    two := 2
    three := 3
    m := NewMap().Set("one", &one).Set("two", &two).Set("three", &three)

    // do we have the right number of keys?
    keys := m.Keys()
    if len(keys) != 3 {
        t.Logf("wrong size keys: %d", len(keys))
        t.FailNow()
    }

    // do we have the right keys?
    sort.Strings(keys)
    if keys[0] != "one" {
        t.Errorf("unexpected key: %s", keys[0])
    }
    if keys[1] != "three" {
        t.Errorf("unexpected key: %s", keys[1])
    }
    if keys[2] != "two" {
        t.Errorf("unexpected key: %s", keys[2])
    }


    // do we have the right values?
    vp, ok := m.Lookup("one");
    if !ok {
        t.Logf("missing value for one")
        t.FailNow()
    }
    if v := vp.(*int); *v != 1 {
        t.Errorf("wrong value: %d\n", *v)
    }
    vp, ok = m.Lookup("two");
    if !ok {
        t.Logf("missing value for two")
        t.FailNow()
    }
    if v := vp.(*int); *v != 2 {
        t.Errorf("wrong value: %d\n", *v)
    }
    vp, ok = m.Lookup("three");
    if !ok {
        t.Logf("missing value for three")
        t.FailNow()
    }
    if v := vp.(*int); *v != 3 {
        t.Errorf("wrong value: %d\n", *v)
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
