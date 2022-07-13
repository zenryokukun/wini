package main

type (
	Pointer interface {
		Ptr() *lnode
	}
	Ranger interface {
		Range() (*lnode, *lnode)
	}
	Printer interface {
		Print() string
	}
)

//Call this to check how your File,Section,Keyvals look like as string.
func Check(rg Ranger) string {
	h, t := rg.Range()
	str := ""
	for {
		str += h.text
		if h == t {
			break
		}
		str += "\n"
		h = h.next
	}
	return str
}

//Inserts rg2 after rg1
func insertBlock(rg1, rg2 Ranger) {
	_, t1 := rg1.Range()
	h2, t2 := rg2.Range()
	nxt := t1.next
	//link t1<->h2
	t1.next = h2
	h2.prev = t1
	//link t2 <-> nxt
	t2.next = nxt
	if nxt != nil {
		nxt.prev = t2
	}
}

//Swaps rg1 and rg2
func swap(rg1, rg2 Ranger) {
	h1, t1 := rg1.Range()
	h2, t2 := rg2.Range()

	if isAdjacent(h1, t1, h2, t2) {
		swapAdjacent(h1, t1, h2, t2)
	} else {
		swapRange(h1, t1, h2, t2)
	}
}

func swapAdjacent(h1, t1, h2, t2 *lnode) {
	var left_h, left_t, right_h, right_t *lnode
	if t1.next == h2 {
		left_h, left_t, right_h, right_t = h1, t1, h2, t2
	} else if t2.next == h1 {
		left_h, left_t, right_h, right_t = h2, t2, h1, t1
	}

	prev := left_h.prev
	next := right_t.next

	//link t2 <-> h1
	right_t.next = left_h
	left_h.prev = right_t
	//link t1 <-> next
	left_t.next = next
	if next != nil {
		next.prev = left_t
	}
	//link h2 <-> prev
	right_h.prev = prev
	if prev != nil {
		prev.next = right_h
	}

	//Add empty line between section after swapping.
	if right_t.ntype != EMPTY {
		right_t.insert(&lnode{ntype: EMPTY, identifier: right_t.identifier, text: ""})
	}

}

func swapRange(h1, t1, h2, t2 *lnode) {
	tmpHeadPrev := h1.prev
	tmpTailNext := t1.next
	h1.prev = h2.prev
	if h2.prev != nil {
		h2.prev.next = h1
	}
	t1.next = t2.next
	if t2.next != nil {
		t2.next.prev = t1
	}
	h2.prev = tmpHeadPrev
	if tmpHeadPrev != nil {
		tmpHeadPrev.next = h2
	}
	t2.next = tmpTailNext
	if tmpTailNext != nil {
		tmpTailNext.prev = t2
	}

	//Add empty line between section after swapping.
	if t1.ntype != EMPTY {
		t1.insert(&lnode{ntype: EMPTY, identifier: t1.identifier, text: ""})
	}
	if t2.ntype != EMPTY {
		t2.insert(&lnode{ntype: EMPTY, identifier: t2.identifier, text: ""})
	}
}

func isAdjacent(h1, t1, h2, t2 *lnode) bool {
	return (t1.next == h2 && h2.prev == t1) || (t2.next == h1 && h1.prev == t2)
}

//Pops range of nodes.
//Ranger will be unlikned,which means that
//head.prev,tail.prev of rg.Range() will be <nil>.
func pop(rg Ranger) {
	h, t := rg.Range()
	prev := h.prev
	next := t.next
	if prev != nil {
		prev.next = next
	}
	if next != nil {
		next.prev = prev
	}
	//unlink range
	h.prev = nil
	t.next = nil
}

func headBlock(ptr Pointer) *lnode {
	node := ptr.Ptr()
	for {
		if node.prev == nil {
			break
		}
		if node.identifier != node.prev.identifier {
			break
		}
		node = node.prev
	}
	return node
}

func tailBlock(ptr Pointer) *lnode {
	node := ptr.Ptr()
	for {
		if node.next == nil {
			break
		}
		if node.identifier != node.next.identifier {
			break
		}
		node = node.next
	}
	return node
}
