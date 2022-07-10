package main

type (
	keyval struct {
		block
		key string
		val string
	}

	KeyVals []*keyval
)

func NewKeyVal(key, val string) *keyval {
	key = trimSpaces(key)
	val = trimSpaces(val)
	text := genKeyValText(key, val)
	l := &lnode{}
	l.setType(KEYVAL)
	l.setIdentifier(key)
	l.setText(text)
	kv := &keyval{key: key, val: val}
	kv.ptr = l
	return kv
}

func (kv *keyval) Ptr() *lnode {
	return kv.ptr
}

//Returns the head and the tail of `KeyVal`
func (kv *keyval) Range() (*lnode, *lnode) {
	h, t := headBlock(kv), tailBlock(kv)
	return h, t
}

func (kv *keyval) ChangeKey(key string) *keyval {
	key = trimSpaces(key)
	//update underlying node.
	kv.update(key, kv.val)
	return kv
}

func (kv *keyval) ChangeVal(val string) *keyval {
	val = trimSpaces(val)
	kv.update(kv.key, val)
	return kv
}

func (kv *keyval) ChangeKeyVal(key, val string) *keyval {
	key = trimSpaces(key)
	val = trimSpaces(val)
	kv.update(key, val)
	return kv
}

//Returns the `val` field.
func (kv *keyval) Val() string {
	return kv.val
}

func (kv *keyval) AddCom(texts ...string) *keyval {
	kv.block.addCom(kv.ptr, KEYCOM, kv.key, texts...)
	return kv
}

func (kv *keyval) update(key, val string) {
	kv.key = key
	kv.val = val
	newName := genKeyValText(key, val)
	kv.ptr.setText(newName)
}
