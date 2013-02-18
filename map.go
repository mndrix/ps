// Fully persistent data structures. A persistent data structure is a data
// structure that always preserves the previous version of itself when
// it is modified. Such data structures are effectively immutable,
// as their operations do not update the structure in-place, but instead
// always yield a new structure.
//
// Persistent
// data structures typically share structure among themselves.  This allows
// operations to avoid copying the entire data structure.
package ps

import . "fmt"

import "bytes"
import "hash/fnv"

// Any is a shorthand for Go's verbose interface{} type.
type Any interface{}

// A Map associates unique keys (type string) with values (type Any).
type Map interface {
    // IsNil returns true if the Map is empty
    IsNil() bool

    // Set returns a new map in which key and value are associated.
    // If the key didn't exist before, it's created; otherwise, the
    // associated value is changed.
    // This operation is O(log N) in the number of keys.
    Set(key string, value Any) Map

    // Delete returns a new map with the association for key, if any, removed.
    // This operation is O(log N) in the number of keys.
    Delete(key string) Map

    // Lookup returns the value associated with a key, if any.  If the key
    // exists, the second return value is true; otherwise, false.
    // This operation is O(log N) in the number of keys.
    Lookup(key string) (Any, bool)

    // Size returns the number of key value pairs in the map.
    // This takes O(1) time.
    Size() int

    // ForEach executes a callback on each key value pair in the map.
    ForEach(f func(key string, val Any))

    // Keys returns a slice with all keys in this map.
    // This operation is O(N) in the number of keys.
    Keys() []string

    String() string
}

// Immutable (i.e. persistent) associative array
type tree struct {
    count   int
    hash    uint64  // hash of the key (used for tree balancing)
    key     string
    value   Any
    left    *tree
    right   *tree
}
var nilMap = &tree{}

// Recursively set nilMap's subtrees to point at itself.
// This eliminates all nil pointers in the map structure.
// All map nodes are created by cloning this structure so
// they avoid the problem too.
func init () {
    nilMap.left = nilMap
    nilMap.right = nilMap
}

// NewMap allocates a new, persistent map from strings to values of
// any type.
// This is currently implemented as a path-copying binary tree.
func NewMap() Map {
    return nilMap
}

func (self *tree) IsNil() bool {
    return self == nilMap
}

// clone returns an exact duplicate of a tree node
func (self *tree) clone() *tree {
    var m tree
    m.count = self.count
    m.hash  = self.hash
    m.key   = self.key
    m.value = self.value
    m.left  = self.left
    m.right = self.right
    return &m
}

// hashKey returns a hash code for a given string
func hashKey(key string) uint64 {
    hasher := fnv.New64()
    Fprint(hasher, key)
    return hasher.Sum64()
}

// Set returns a new map similar to this one but with key and value
// associated.  If the key didn't exist, it's created; otherwise, the
// associated value is changed.
func (self *tree) Set(key string, value Any) Map {
    hash := hashKey(key)
    return setLowLevel(self, hash, key, value)
}

func setLowLevel(self *tree, hash uint64, key string, value Any) *tree {
    if self.IsNil() { // an empty tree is easy
        m := self.clone()
        m.count = 1
        m.hash  = hash
        m.key   = key
        m.value = value
        return m
    }

    if hash < self.hash { // insert into left subtree
        m := self.clone()
        m.left = setLowLevel(self.left, hash, key, value)
        recalculateCount(m)
        return m
    }
    if hash > self.hash { // insert into right subtree
        m := self.clone()
        m.right = setLowLevel(self.right, hash, key, value)
        recalculateCount(m)
        return m
    }

    // replacing a key's previous value
    m := self.clone()
    m.value = value
    return m
}

// modifies a map by recalculating its key count based on the counts
// of its subtrees
func recalculateCount(m *tree) {
    count := 0
    if !m.left.IsNil() {
        count += m.left.Size()
    }
    if !m.right.IsNil() {
        count += m.right.Size()
    }
    m.count = count + 1 // add one to count ourself
}

func (m *tree) Delete(key string) Map {
    hash := hashKey(key)
    newMap, _ := deleteLowLevel(m, hash)
    return newMap
}

func deleteLowLevel(self *tree, hash uint64) (*tree, bool) {
    // empty trees are easy
    if self.IsNil() {
        return self, false
    }

    if hash < self.hash { // look in the left subtree
        newLeft, found := deleteLowLevel(self.left, hash)
        if !found {
            return self, false
        }
        newMap := self.clone()
        newMap.left = newLeft
        recalculateCount(newMap)
    }
    if hash > self.hash { // look in the right subtree
        newRight, found := deleteLowLevel(self.right, hash)
        if !found {
            return self, false
        }
        newMap := self.clone()
        newMap.right = newRight
        recalculateCount(newMap)
    }

    // we must delete our own node
    if self.isLeaf() {  // we have no children
        return nilMap, true
    }
    if self.subtreeCount() == 1 { // only one subtree
        if self.hasLeft() {  // it's the left one
            return self.left, true
        }
        return self.right, true  // it's the right one
    }

    // find a node to replace us
    if self.left.Size() > self.right.Size() {  // make left side smaller
        replacement, newLeft := self.left.deleteRightmost()
        newMap := replacement.clone()
        newMap.left = newLeft
        newMap.right = self.right
        recalculateCount(newMap)
        return newMap, true
    }

    // make right side smaller
    replacement, newRight := self.right.deleteLeftmost()
    newMap := replacement.clone()
    newMap.right = newRight
    newMap.left = self.left
    recalculateCount(newMap)
    return newMap, true
}

// delete the left or rightmost node in a tree returning the node that
// was deleted and the tree left over after its deletion
func (m *tree) deleteRightmost() (*tree, *tree) {
    if m.isLeaf() {
        return m, nilMap
    }
    if m.hasRight() {
        deleted, newRight := m.right.deleteRightmost()
        newMap := m.clone()
        newMap.right = newRight
        recalculateCount(newMap)
        return deleted, newMap
    }

    deleted := m.clone()
    deleted.left = nilMap
    return deleted, m.left
}
func (m *tree) deleteLeftmost() (*tree, *tree) {
    if m.isLeaf() {
        return m, nilMap
    }
    if m.hasLeft() {
        deleted, newLeft := m.left.deleteLeftmost()
        newMap := m.clone()
        newMap.left = newLeft
        recalculateCount(newMap)
        return deleted, newMap
    }

    deleted := m.clone()
    deleted.count = 1
    deleted.right = nilMap
    return deleted, m.right
}

// hasLeft and hasRight return true if this tree has a left or right subtree
func (m *tree) hasLeft() bool {
    return !m.left.IsNil()
}
func (m *tree) hasRight() bool {
    return !m.right.IsNil()
}

// isLeaf returns true if this is a leaf node
func (m *tree) isLeaf() bool {
    return m.Size() == 1
}

// returns the number of child subtrees we have
func (m *tree) subtreeCount() int {
    count := 0
    if m.hasLeft() {
        count++
    }
    if m.hasRight() {
        count++
    }
    return count
}

func (m *tree) Lookup(key string) (Any, bool) {
    hash := hashKey(key)
    return lookupLowLevel(m, hash)
}

func lookupLowLevel(self *tree, hash uint64) (Any, bool) {
    if self.IsNil() { // an empty tree is easy
        return nil, false
    }

    if hash < self.hash { // look in the left subtree
        return lookupLowLevel(self.left, hash)
    }
    if hash > self.hash { // look in the right subtree
        return lookupLowLevel(self.right, hash)
    }

    // we found it
    return self.value, true
}

func (m *tree) Size() int {
    return m.count
}

func (m *tree) ForEach(f func(key string, val Any)) {
    if m.IsNil() {
        return
    }

    // left branch
    if !m.left.IsNil() {
        m.left.ForEach(f)
    }

    // ourself
    f(m.key, m.value)

    // right branch
    if !m.right.IsNil() {
        m.right.ForEach(f)
    }
}

func (m *tree) Keys() []string {
    keys := make([]string, m.Size())
    i := 0
    m.ForEach( func (k string, v Any) {
        keys[i] = k
        i++
    })
    return keys
}

// make it easier to display maps for debugging
func (m *tree) String() string {
    keys := m.Keys()
    buf := bytes.NewBufferString("{")
    for _, key := range keys {
        val, _ := m.Lookup(key)
        Fprintf(buf, "%s: %s, ", key, val)
    }
    Fprintf(buf, "}\n")
    return buf.String()
}
