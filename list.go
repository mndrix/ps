// Immutable (i.e. persistent) list
package ps

type List struct {
    depth   int     // the number of nodes after, and including, this one
    value   Any
    tail    *List
}

// NewList returns a new, empty list
func NewList() *List {
    var list List
    return &list
}

// IsNil returns true if the list is empty
func (self *List) IsNil() bool {
    return self.depth == 0
}

// Size returns the list's length
func (self *List) Size() int {
    return self.depth
}

// Cons returns a new list with val added onto the head
func (tail *List) Cons(val Any) *List {
    var list List
    list.depth = tail.depth + 1
    list.value = val
    list.tail = tail
    return &list
}

// Head returns the first element in the list or panics if the list is empty
func (self *List) Head() Any {
    if self.IsNil() {
        panic("Called Head() on an empty list")
    }

    return self.value
}

// Tail returns the tail of this list or panics if the list is empty
func (self *List) Tail() *List {
    if self.IsNil() {
        panic("Called Tail() on an empty list")
    }

    return self.tail
}

// ForEach executes a callback for each value in the list
func (self *List) ForEach(f func(Any)) {
    if self.IsNil() {
        return
    }
    f(self.Head())
    self.Tail().ForEach(f)
}

// Reverse returns a list with elements in opposite order as this list
func (self *List) Reverse() *List {
    reversed := NewList()
    self.ForEach( func (v Any) { reversed = reversed.Cons(v) })
    return reversed
}
