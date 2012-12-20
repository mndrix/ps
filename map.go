// Immutable (i.e. persistent) data structures
package ps

import "bytes"
import "encoding/gob"

type Any interface{}

type Map map[string]Any

// NewMap allocates a new, persistent map from strings to any value
func NewMap() Map {
    m := make(map[string]Any)
    return m
}

// Set returns a new map similar to this one but with key and value
// associated.  If the key didn't exist, it's created; otherwise, the
// associated value is changed.
func (prev Map) Set(key string, value Any) Map {
    next := prev.Clone()
    next[key] = value
    return next
}

// Delete returns a new map with the association for key, if any, removed
func (prev Map) Delete(key string) Map {
    next := prev.Clone()
    delete(next, key)
    return next
}

// Lookup returns a pair of values.  The first is the value, the second is
// true if the value exists.
func (m Map) Lookup(key string) (Any, bool) {
    v, ok := m[key]
    return v, ok
}

// Size returns the number of key value pairs in the map
func (m Map) Size() int {
    return len(m)
}

// Clone deep copies `src` into `dst`
func (prev Map) Clone() Map {
    next := NewMap()
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    dec := gob.NewDecoder(&buf)
    enc.Encode(prev)
    dec.Decode(&next)
    return next
}
