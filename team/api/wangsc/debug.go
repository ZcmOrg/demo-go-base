/**
 * debug 调试信息相关
 */

package debug

import (
	"debug/dwarf"
	"debug/elf"
	"io"
	"reflect"
	"testing"
)

//dwarf	提供了对从可执行文件加载的DWARF调试信息的访问，这个包对于实现Go语言的调试器非常有价值

//elf	实现了对ELF对象文件的访问。 ELF是一种常见的二进制可执行文件和共享库的文件格式。 Linux采用了ELF格式
//gosym	访问Go语言二进制程序中的调试信息。对于可视化调试很有价值
//macho	实现了对http://developer.apple.com/mac/library/documentation/DeveloperTools/Conceptual/MachORuntime/Reference/reference.html 所定义的Mach-O对象文件的访问
//pe	实现了对PE（ Microsoft Windows Portable Executable）文件的访问

func newElf(name string) *dwarf.Data {
	f, err := elf.Open(name)
	if err != nil {
		t.Fatal(err)
	}

	d, err := f.DWARF()
	if err != nil {
		t.Fatal(err)
	}
	return d
}

type wantRange struct {
	pc     uint64
	ranges [][2]uint64
}

func testRanges() {
	d := newElf("")
	want := []wantRange{
		{0x400500, [][2]uint64{{0x400500, 0x400549}, {0x400400, 0x400408}}},
		{0x400400, [][2]uint64{{0x400500, 0x400549}, {0x400400, 0x400408}}},
		{0x400548, [][2]uint64{{0x400500, 0x400549}, {0x400400, 0x400408}}},
		{0x400407, [][2]uint64{{0x400500, 0x400549}, {0x400400, 0x400408}}},
		{0x400408, nil},
		{0x400449, nil},
		{0x4003ff, nil},
	}
	r := d.Reader()
	for _, w := range want {
		entry, err := r.SeekPC(w.pc)
		if err != nil {
			if w.ranges != nil {
				t.Errorf("%s: missing Entry for %#x", name, w.pc)
			}
			if err != ErrUnknownPC {
				t.Errorf("%s: expected ErrUnknownPC for %#x, got %v", name, w.pc, err)
			}
			continue
		}

		ranges, err := d.Ranges(entry)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if !reflect.DeepEqual(ranges, w.ranges) {
			t.Errorf("%s: for %#x got %x, expected %x", name, w.pc, ranges, w.ranges)
		}
	}
}

func TestSplit() {
	d := newElf("")
	r := d.Reader()
	e, err := r.Next()
	if err != nil {
		t.Fatal(err)
	}
	const AttrGNUAddrBase Attr = 0x2133
	f := e.AttrField(AttrGNUAddrBase)
	if _, ok := f.Val.(int64); !ok {
		t.Fatalf("bad attribute value type: have %T, want int64", f.Val)
	}
	if f.Class != ClassUnknown {
		t.Fatalf("bad class: have %s, want %s", f.Class, ClassUnknown)
	}
}

type LineEntry struct {
	Address       uint64
	OpIndex       int
	Line          int
	Column        int
	IsStmt        bool
	BasicBlock    bool
	PrologueEnd   bool
	EpilogueBegin bool
	ISA           int
	Discriminator int
	EndSequence   bool
}

func testLineTable() {
	d := newElf("")
	var got []LineEntry
	dr := d.Reader()
	for {
		ent, err := dr.Next()
		if err != nil {
			t.Fatal("dr.Next:", err)
		} else if ent == nil {
			break
		}
		lr, _ := d.LineReader(ent)
		for {
			var line LineEntry
			lr.Next(&line)
			got = append(got, line)
		}
	}
}

func TestLineSeek(t *testing.T) {
	d := newElf("")

	cu, _ := d.Reader().Next()
	lr, _ := d.LineReader(cu)

	var line LineEntry
	var posTable []LineReaderPos
	var table []LineEntry
	for {
		posTable = append(posTable, lr.Tell())

		lr.Next(&line)
		table = append(table, line)
	}
	lr.Reset()
	for i := len(posTable) - 1; i >= 0; i-- {
		lr.Seek(posTable[i])
		err := lr.Next(&line)
		if i == len(posTable)-1 {
			if err != io.EOF {
				t.Fatal("expected io.EOF after seek to end, got", err)
			}
		} else if err != nil {
			t.Fatal("lr.Next after seek to", posTable[i], "failed:", err)
		} else if line != table[i] {
			t.Fatal("lr.Next after seek to", posTable[i], "returned", line, "instead of", table[i])
		}
	}
}
