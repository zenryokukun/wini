package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type File map[string]*Section

//`fpath` is the config file path.
//It reads up the file, and link each line as linked-list.
//Returns the linked-list as a File map.
func Load(fpath string) File {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	head := newLNode(scanner.Text())
	tail := head

	for scanner.Scan() {
		line := scanner.Text()
		node := newLNode(line)
		tail.insert(node)
		tail = node
	}

	classifyComments(head)
	classifyEmptyLines(tail)

	return newFile(head)
}

// internal. Called from Load.
func newFile(l *lnode) File {
	var file File = File{}
	//scroll to head
	h := head(l)
	for h != nil {
		if h.ntype == SEC || h.ntype == SECCOM {
			h = addSecInfo(file, h)
		} else {
			h = h.next
		}
	}
	return file
}

//Returns empty File. Used to generate ini file from scratch.
func NewFile() File {
	return File{}
}

//Merges File maps after the caller.
//They will be merged in order they are passed.
func (f File) Merge(fs ...File) {
	_, prevTail := f.Range()
	for _, nf := range fs {
		for k, v := range nf {
			f[k] = v
		}
		h, _ := nf.Range()
		if prevTail.ntype != EMPTY && h.ntype != EMPTY {
			empty := &lnode{ntype: EMPTY, identifier: prevTail.identifier}
			prevTail.insert(empty)
			prevTail = empty
		}
		prevTail.insertBlock(nf)
		_, prevTail = f.Range()
	}
}

func (f File) ChangeSectionName(name, newName string) File {
	checkSecSym(newName)
	sec := f[name]
	sec.changeName(newName)
	//change map key.
	f[newName] = sec
	delete(f, name)
	return f
}

//Adds section to file.
func (f File) AddSec(ns ...*Section) File {
	//File has no keys.
	if len(f) == 0 {
		fsec := ns[0]
		f[fsec.name] = ns[0]
		for _, s := range ns[1:] {
			insertBlock(f, s)
			f[s.name] = s
		}
		return f
	}
	//Has some keys..
	for _, s := range ns {
		insertBlock(f, s)
		f[s.name] = s
	}
	return f
}

//Swaps sections.`s1` and `s2` are section keys without "[" and "]".
func (f File) Swap(s1, s2 string) error {
	k1, ok := f[s1]
	if !ok {
		return fmt.Errorf("section name not found:%v", s1)
	}
	k2, ok := f[s2]
	if !ok {
		return fmt.Errorf("section name not found:%v", s2)
	}
	swap(k1, k2)
	return nil
}

//Pops `section` from `file`.
//`name` is the section name to pop.
func (f File) Pop(name string) {
	sec, ok := f[name]
	if ok {
		pop(sec)
		delete(f, name) //delete from map.
	}
}

//Pops all comments from `file`.
func (f File) PopAllCom() {
	//pop all comment nodes.
	h, _ := f.Range()
	for h != nil {
		if h.ntype == KEYCOM || h.ntype == SECCOM {
			//You need to set the popping node to another variable
			//,and set h to h.next before popping,because `pop`` will set .next and .prev to <nil>.
			//Otherwise, h = h.next will be <nil> all the time and breaks
			//before searching 'till the end of the linked-list.
			nodeToPop := h
			h = h.next
			pop(nodeToPop)
		} else {
			h = h.next
		}
	}
	//clear comments list
	for _, sec := range f {
		sec.comments = Comments{}
		for _, kv := range sec.data {
			kv.comments = Comments{}
		}
	}
}

//Pops all empty lines from `file`.
func (f File) PopEmptyLines() {
	h, _ := f.Range()
	for h != nil {
		if h.ntype == EMPTY {
			nodeToPop := h
			h = h.next
			pop(nodeToPop)
		} else {
			h = h.next
		}
	}
}

// Removes left indents.
// " " and tabs are considered as indents.
func (f File) RemoveIndent() {
	h, _ := f.Range()
	for {
		if h == nil {
			break
		}
		txt := strings.TrimLeft(h.text, " ")
		txt = strings.TrimLeft(txt, "\t")
		h.setText(txt)
		h = h.next
	}
}

//Returns the head and the tail of `File`.
//This will be the full linked-list.
func (f File) Range() (*lnode, *lnode) {
	var ptr *lnode
	//Just iterate once, to get the ptr to search from.
	/*
		Todo:
			Change the File to struct and hold `head` and `tail` field.
	*/
	for _, sec := range f {
		ptr = sec.ptr
		break
	}
	h := head(ptr)
	t := tail(ptr)
	return h, t
}

//**********
// Saves ini file.
// It will scroll to the head of the linked-list that Pointer belongs to,
// and writes out the whole text.
// New bkupfile will be created,when there is none.
func (f File) Save(fpath string) {
	head, _ := f.Range()
	str := asString(head)
	dir := filepath.Dir(fpath)
	fname := filepath.Base(fpath)
	bkpath := filepath.Join(dir, "winiBK_"+fname)
	_, err := os.Stat(fpath)
	_, errBK := os.Stat(bkpath)
	if !os.IsNotExist(err) && os.IsNotExist(errBK) {
		//create backup file,when fpath and bkup file do not exist.
		backupFile(fpath, bkpath)
	}
	save(str, fpath)
}

// Saves ini file.
// Removes all existing
// new lines,and formats output strings
// according to the parameters.
// secLines -> number of `extra` empty lines between Sections.
//             if kvLines > 0, then kvLines + secLines will be the
//             number of empty lines before sections.
// kvLines  -> number of empty Lines between keyvals.
// indent   -> number of indentation of keyvals.
func (f File) Savef(fpath string, secLines, kvLines, indent int) {
	// Removes all empty lines and indents.
	// Call this before calling Range(),because
	// the head could be an empty line.
	f.PopEmptyLines()
	f.RemoveIndent()
	head, _ := f.Range()
	str := asStringf(head, secLines, kvLines, indent)
	dir := filepath.Dir(fpath)
	fname := filepath.Base(fpath)
	bkpath := filepath.Join(dir, "winiBK_"+fname)

	_, err := os.Stat(fpath)
	_, errBK := os.Stat(bkpath)
	if !os.IsNotExist(err) && os.IsNotExist(errBK) {
		//create backup file,when fpath and bkup file do not exist.
		backupFile(fpath, bkpath)
	}
	save(str, fpath)
}

// internal. Called from Save.
func save(text, fpath string) {
	f, err := os.Create(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write([]byte(text))
	if err != nil {
		fmt.Println(err)
	}
}

//internal. Called from Save.
func backupFile(fpath, bkpath string) {
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Println(fpath)
		log.Fatal(err)
	}
	defer f.Close()

	bf, err := os.Create(bkpath)
	if err != nil {
		fmt.Println(bf)
		log.Fatal(err)
	}
	defer bf.Close()
	io.Copy(bf, f)
}

//internal. Called from Save()
func asString(n *lnode) string {
	str := ""
	for {
		str += n.text
		if n.next == nil {
			break
		}
		str += "\n"

		//Adding extra empty-line before [section] or its commets,
		//if there is none.
		if n.next.ntype == SEC {
			if n.ntype != SECCOM && n.ntype != EMPTY {
				str += "\n"
			}
		} else if n.next.ntype == SECCOM {
			if n.ntype != SECCOM && n.ntype != EMPTY {
				str += "\n"
			}
		}

		n = n.next
	}
	return str
}

// internal. Called from Savef
func asStringf(n *lnode, secLines, kvLines, indent int) string {
	str := ""
	for n != nil {
		txt := n.text
		if n.ntype == KEYVAL || n.ntype == KEYCOM {
			txt = getStr(" ", indent) + txt
		}
		if n.ntype == SEC || n.ntype == KEYVAL {
			// new lines after section or keyval
			txt += getStr("\n", kvLines)
		}

		if n.next != nil && (n.ntype != SEC && n.ntype != SECCOM) && (n.next.ntype == SECCOM || n.next.ntype == SEC) {
			// new lines before section

			txt += getStr("\n", secLines)
		}

		str += txt + "\n"
		n = n.next
	}

	return str
	/*
		for {
			//str += n.text
			t := n.text
			hit := false

			if n.ntype == KEYVAL || n.ntype == KEYCOM {
				str += getStr(" ", indent) + t
			}

			if n.next == nil {
				str += t
				break
			}

			if n.ntype == SEC || n.ntype == KEYVAL {
				// new lines after section or keyval
				str += t + getStr("\n", kvLines)
				hit = true
			}

			if (n.ntype != SEC && n.ntype != SECCOM) && (n.next.ntype == SECCOM || n.next.ntype == SEC) {
				// new lines before section
				if hit == false {
					str += t + getStr("\n", secLines)
				}
				hit = true
			}

			if hit == false {
				str += t
			}

			str += "\n"

			n = n.next
		}
		return str
	*/
}

//Called from newFile.
func addSecInfo(f File, l *lnode) *lnode {
	for l != nil {
		sec := &Section{}

		if l.ntype == SECCOM || l.ntype == SEC {
			l = _addSecInfo(sec, l)
			l = addKeyValInfo(sec, l)
		} else {
			l = l.next
		}

		if len(sec.name) > 0 {
			f[sec.name] = sec
		}
	}
	return l
}

//Adds section name and ptr.
func _addSecInfo(s *Section, l *lnode) *lnode {
	id := l.identifier
	for l.identifier == id {
		if l.ntype == SECCOM {
			s.comments = append(s.comments, newCommentFromNode(l.text, l))
		} else if l.ntype == SEC {
			s.name = l.identifier
			s.ptr = l
		}
		l = l.next
	}
	return l
}

//Adds KeyVals to Section.
func addKeyValInfo(s *Section, l *lnode) *lnode {
	//var kvs KeyVals = KeyVals{}
	for l != nil {
		if l.ntype == SEC || l.ntype == SECCOM {
			break
		}
		if l.ntype == KEYCOM || l.ntype == KEYVAL {
			l = _addKeyValInfo(s, l)
		} else {
			l = l.next
		}
	}
	//s.data = kvs
	return l
}

func _addKeyValInfo(s *Section, l *lnode) *lnode {
	id := l.identifier
	var kv *keyval = &keyval{}
	for l != nil && l.identifier == id {
		// if l == nil {
		// 	break
		// }
		if l.ntype == KEYCOM {
			kv.comments = append(kv.comments, newCommentFromNode(l.text, l))
		} else if l.ntype == KEYVAL {
			kv.key, kv.val = sepSplit(l.text)
			kv.ptr = l
		}
		l = l.next
	}
	if kv.ptr != nil {
		s.addKeyVals(kv)
	}
	return l
}

func getStr(chr string, cnt int) string {
	str := ""
	for i := 0; i < cnt; i++ {
		str += chr
	}
	return str
}
