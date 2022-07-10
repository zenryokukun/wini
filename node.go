package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	EMPTY     = iota //new-line code
	SEC              //Section
	KEYVAL           //KeyVal
	SECCOM           //Section Comment
	KEYCOM           //KeyVal Comment
	UNDEFINED        //other line
)

var (
	commentSymbol = []string{"#", ";"}
	sectionSymbol = []string{"[", "]"}
	sepSymbol     = "="
)

type lnode struct {
	ntype      int
	identifier string
	text       string
	next       *lnode
	prev       *lnode
}

func ChangeSepSym(ch string) {
	sepSymbol = ch
}

func ChangeSectionSym(left, right string) {
	sectionSymbol[0], sectionSymbol[1] = left, right
}

func AddCommentSym(ch string) {
	commentSymbol = append(commentSymbol, ch)
}

func newLNode(l string) *lnode {
	tp := which(l)
	id := ""
	if tp == SEC {
		id = getSectionName(l)
	} else if tp == KEYVAL {
		id = getKeyName(l)
	}
	return &lnode{
		ntype:      tp,
		identifier: id,
		text:       l,
	}
}

func (l *lnode) Range() (*lnode, *lnode) {
	return l, l
}

//Inserts single node.
func (l *lnode) insert(n *lnode) {
	nxt := l.next //backup
	//link l <-> n
	l.next = n
	n.prev = l
	//link l <-> nxt
	n.next = nxt
	if nxt != nil {
		nxt.prev = l
	}
}

//Inserts single node before the caller.
func (l *lnode) insertBefore(n *lnode) {
	prev := l.prev
	//link n <-> l
	l.prev = n
	n.next = l
	//link prev <-> n
	n.prev = prev
	if prev != nil {
		prev.next = n
	}
}

//Inserts block of nodes after the caller.
func (l *lnode) insertBlock(rg Ranger) {
	h, t := rg.Range()

	nxt := l.next
	//link l <-> n
	l.next = h
	h.prev = l
	//link t <-> nxt
	t.next = nxt
	if nxt != nil {
		nxt.prev = t
	}
}

func (l *lnode) setIdentifier(id string) {
	l.identifier = id
}

func (l *lnode) setType(ntype int) {
	l.ntype = ntype
}

func (l *lnode) setText(text string) {
	l.text = text
}

//*****************************************************
// helper functions
//*****************************************************

func tail(l *lnode) *lnode {
	if l.next == nil {
		return l
	}
	return tail(l.next)
}

func head(l *lnode) *lnode {
	if l.prev == nil {
		return l
	}
	return head(l.prev)
}

func trimSpaces(line string) string {
	tm := strings.Trim(line, " ")
	tm = strings.Trim(tm, "\t")
	return tm
}

func isComment(line string) bool {
	tm := trimSpaces(line)
	if len(tm) == 0 {
		return false
	}
	ch := tm[:1]
	for _, v := range commentSymbol {
		if v == ch {
			return true
		}
	}
	return false
}

func isSection(line string) bool {
	tm := trimSpaces(line)
	if len(tm) == 0 {
		return false
	}
	left := sectionSymbol[0]
	right := sectionSymbol[1]
	return strings.HasPrefix(tm, left) && strings.HasSuffix(tm, right)
}

func which(line string) int {
	tm := trimSpaces(line)
	if len(tm) == 0 {
		return EMPTY
	}
	if len(tm) == 1 && tm == "\n" {
		return EMPTY
	}
	if isComment(tm) {
		return UNDEFINED
	}
	if isSection(tm) {
		return SEC
	}
	return KEYVAL
}

func getSectionName(line string) string {
	left, right := sectionSymbol[0], sectionSymbol[1]
	tm := trimSpaces(line)
	tm = strings.TrimLeft(tm, left)
	tm = strings.TrimRight(tm, right)
	return tm
}

func getKeyName(line string) string {
	k, _ := sepSplit(line)
	return k
}

func genKeyValText(key, val string) string {
	text := key
	if len(val) > 0 {
		text += sepSymbol + val
	}
	return text
}

//Splits keyVal text by sepSymbol.
func sepSplit(line string) (string, string) {
	list := strings.Split(line, sepSymbol)
	var k, v string

	if len(list) == 1 {
		k = trimSpaces(list[0])
		v = ""

	} else if len(list) == 2 {
		k, v = trimSpaces(list[0]), trimSpaces(list[1])

	} else if len(list) > 2 {
		k = trimSpaces(list[0])
		for i := 1; i < len(list); i++ {
			v += trimSpaces(list[i])
			if i < len(list)-1 {
				v += sepSymbol
			}
		}
	}
	return k, v
}

func classifyComments(l *lnode) {
	cms := []*lnode{}
	for l != nil {
		if l.ntype == UNDEFINED {
			cms = append(cms, l)
		} else if l.ntype == SEC || l.ntype == KEYVAL {
			_classifyComments(cms, l)
			cms = []*lnode{}
		}
		l = l.next
	}
}

func _classifyComments(cms []*lnode, l *lnode) {
	var tp int
	if l.ntype == SEC {
		tp = SECCOM
	} else if l.ntype == KEYVAL {
		tp = KEYCOM
	}
	for _, c := range cms {
		c.setType(tp)
		c.setIdentifier(l.identifier)
	}
}

//`l` should be the tail of the linked-list
func classifyEmptyLines(l *lnode) {
	emp := []*lnode{}
	for l != nil {
		if l.ntype == EMPTY {
			emp = append(emp, l)
		} else {
			for _, v := range emp {
				v.setIdentifier(l.identifier)
			}
			emp = []*lnode{}
		}
		l = l.prev
	}
}

func checkComSym(text string) {
	chk := string(text[0])
	for _, ch := range commentSymbol {
		if ch == chk {
			return
		}
	}
	log.Fatal(fmt.Sprintf("lacking comment symbol:%v", text))
}

func checkSecSym(text string) {
	left := string(text[0])
	right := string(text[len(text)-1])
	if left != sectionSymbol[0] || right != sectionSymbol[1] {
		log.Fatal(fmt.Sprintf("lacking section symbol:%v", text))
	}
}

// Used when changing section names or keyval key name.
func updateIdentifier(n *lnode, id string) {
	for {
		if n == nil {
			break
		}
		n.setIdentifier(id)
		if n.prev != nil && n.prev.identifier != n.identifier {
			break
		}
		n = n.prev
	}
}
