package main

import (
	"fmt"
)

type Section struct {
	block
	data KeyVals
	name string
}

func NewSection(text string) *Section {
	checkSecSym(text)
	l := &lnode{}
	l.setType(SEC)
	l.setIdentifier(getSectionName(text))
	l.setText(text)
	sec := &Section{name: getSectionName(text)}
	sec.ptr = l
	return sec
}

func (sec *Section) Ptr() *lnode {
	return sec.ptr
}

//Pops `keyval` from `section`.
//`key` is the key to pop.
func (s *Section) Pop(key string) {
	kv := s.Key(key)
	if kv != nil {
		pop(kv)            //pop nodes from the linked-list.
		s.popDataSlice(kv) //make sure to pop s.data.
	}
}

//Called from section.Pop()
//Pops kv from section.data
func (s *Section) popDataSlice(kv *keyval) {
	at := -1
	for i, v := range s.data {
		if kv == v {
			at = i
			break
		}
	}
	if at >= 0 {
		s.data = append(s.data[:at], s.data[at+1:]...)
	}
}

//Returns the head and the tail of `Section`.
func (sec *Section) Range() (*lnode, *lnode) {
	h, t := sec.head(), sec.tail()
	return h, t
}

//Searches key-val that matches `key`, and returns *keyval.
func (s *Section) Key(key string) *keyval {
	for _, kv := range s.data {
		if kv.key == key {
			return kv
		}
	}
	return nil //not found.
}

//Returns all KeyVals under section as a "key"-"val" map.
//Returns <nil> when no keyval is set.
//It extracts only key-val data, discluding key-val comments.
func (s *Section) Data() map[string]string {
	if len(s.data) == 0 {
		return nil
	}
	m := map[string]string{}
	for _, kv := range s.data {
		m[kv.key] = kv.val
	}
	return m
}

//Adds keyval to section.
func (s *Section) AddKeyVal(kvs ...*keyval) *Section {
	var lastkv *keyval
	for _, kv := range kvs {
		insertBlock(s, kv)
		s.data = append(s.data, kv)
		lastkv = kv
	}
	//call this on last keyval.
	if lastkv != nil {
		adjustEmptyLine(lastkv)
	}
	return s

}

func (s *Section) AddCom(texts ...string) *Section {
	s.block.addCom(s.ptr, SECCOM, s.name, texts...)
	return s
}

//Swaps keyvals.`k1` and `k2` are keys of keyvals.
func (s *Section) Swap(k1, k2 string) error {
	var keyval1, keyval2 *keyval
	for _, kv := range s.data {
		if kv.key == k1 {
			keyval1 = kv
		} else if kv.key == k2 {
			keyval2 = kv
		}
	}
	if keyval1 == nil {
		return fmt.Errorf("key not found:%v", k1)
	}
	if keyval2 == nil {
		return fmt.Errorf("key not found:%v", k2)
	}
	swap(keyval1, keyval2)
	return nil
}

//************************************************
// internal functions and methods
//************************************************

func (s *Section) changeName(name string) {
	// Call this before changing name. Range() woul not work.
	updateIdentifier(s.ptr, name)
	s.name = name
	left, right := sectionSymbol[0], sectionSymbol[1]
	nameWithSym := left + name + right
	//update the underlying node.
	s.ptr.setText(nameWithSym)
}

func (s *Section) addKeyVals(kvs ...*keyval) {
	for _, kv := range kvs {
		s.data = append(s.data, kv)
	}
}

//Called from sec.Range()
func (sec *Section) head() *lnode {
	return headBlock(sec)
}

//Called from sec.Range()
func (sec *Section) tail() *lnode {
	node := sec.ptr
	for {
		if node.next == nil {
			break
		}
		np := node.next.ntype
		if np == SEC || np == SECCOM {
			break
		}
		node = node.next
	}
	return node
}

//Swaps empty-line and keyval.
//Last node of section can be a empty line,
//so added keyval could be inserted after it.
//It looks better swapped...
func adjustEmptyLine(kv *keyval) {
	h, t := kv.Range()
	if h == t {
		if t.ntype == KEYVAL && t.next != nil && t.prev != nil {
			if (t.next.ntype == SEC || t.next.ntype == SECCOM) && t.prev.ntype == EMPTY {
				pop(t.prev)
				t.insert(&lnode{ntype: EMPTY, identifier: kv.key})
				//t.insertBlock(&lnode{ntype: EMPTY, identifier: kv.key})

			}
		}
	} else {
		if h.ntype == KEYCOM && t.next != nil && h.prev != nil {
			if (t.next.ntype == SEC || t.next.ntype == SECCOM) && h.prev.ntype == EMPTY {
				pop(h.prev)
				t.insert(&lnode{ntype: EMPTY, identifier: kv.key})
			}
		}
	}
	//printall(h)
}
