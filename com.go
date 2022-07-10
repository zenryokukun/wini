//Comment and block structs.

package main

type (
	Comment struct {
		text string
		ptr  *lnode
	}

	Comments []*Comment

	//Base struct for Section and KeyVal
	block struct {
		comments Comments
		ptr      *lnode
	}
)

//Creates *Comment from text.
func NewComment(ntype int, id, text string) *Comment {
	checkComSym(text)
	l := &lnode{}
	l.setType(ntype)
	l.setIdentifier(id)
	l.setText(text)
	return &Comment{text: text, ptr: l}
}

//Creates Comment from lnode.Used when instantiating `file`.
func newCommentFromNode(text string, ptr *lnode) *Comment {
	return &Comment{text: text, ptr: ptr}
}

func (com *Comment) Ptr() *lnode {
	return com.ptr
}

//range of `comments.`
func (coms Comments) Range() (*lnode, *lnode) {
	//find a comment that is not <nil>.
	var tg *Comment
	for _, c := range coms {
		if c != nil {
			tg = c
			break
		}
	}
	if tg == nil {
		return nil, nil
	}
	h := headBlock(tg)
	//find t.
	t := h
	for {
		if t.next == nil {
			break
		}
		if t.next.ntype == SEC || t.next.ntype == KEYVAL {
			break
		}
		t = t.next
	}
	//h, t := c[0].ptr, c[len(c)-1].ptr
	return h, t
}

func (c *Comment) Range() (*lnode, *lnode) {
	return c.ptr.Range()
}

//Changes Comment text
func (c *Comment) Change(text string) *Comment {
	checkComSym(text)
	c.text = text
	c.ptr.setText(text)
	return c
}

//Pops section `comment` by index.
func (b *block) PopCom(index int) {
	if index < 0 || index > len(b.comments)-1 {
		return
	}
	com := b.comments[index]
	pop(com)
	b.popComSlice(com)
}

//pops all comments of `keyval`.
func (b *block) PopAllCom() {
	pop(b.comments)
	b.comments = Comments{} //init comments.
}

func (b *block) popComSlice(c *Comment) {
	at := -1
	for i, v := range b.comments {
		if v == c {
			at = i
			break
		}
	}
	if at >= 0 {
		b.comments = append(b.comments[:at], b.comments[at+1:]...)
	}
}

//Gets comment.
func (bl *block) Com(i int) *Comment {
	if bl.comments == nil || (i > len(bl.comments)-1) || i < 0 {
		return nil
	}
	return bl.comments[i]
}

//Adds Comment(s).
//`tnode`` is either section.ptr or kv.ptr.
//All comments should be inserted before tnode.
func (bl *block) addCom(tnode *lnode, ntype int, id string, texts ...string) {
	for _, text := range texts {
		com := NewComment(ntype, id, text)
		tnode.insertBefore(com.ptr)
		bl.comments = append(bl.comments, com)
	}
}
