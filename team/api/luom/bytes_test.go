package luom_test

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"unicode"
	"unsafe"
)

var compareTests = []struct {
	a, b []byte
	i    int
}{
	{[]byte(""), []byte(""), 0},
	{[]byte("a"), []byte(""), 1},
	{[]byte(""), []byte("a"), -1},
	{[]byte("abc"), []byte("abc"), 0},
	{[]byte("abd"), []byte("abc"), 1},
	{[]byte("abc"), []byte("abd"), -1},
	{[]byte("ab"), []byte("abc"), -1},
	{[]byte("abc"), []byte("ab"), 1},
	{[]byte("x"), []byte("ab"), 1},
	{[]byte("ab"), []byte("x"), -1},
	{[]byte("x"), []byte("a"), 1},
	{[]byte("b"), []byte("x"), -1},
	// test runtime·memeq's chunked implementation
	{[]byte("abcdefgh"), []byte("abcdefgh"), 0},
	{[]byte("abcdefghi"), []byte("abcdefghi"), 0},
	{[]byte("abcdefghi"), []byte("abcdefghj"), -1},
	{[]byte("abcdefghj"), []byte("abcdefghi"), 1},
	// nil tests
	{nil, nil, 0},
	{[]byte(""), nil, 0},
	{nil, []byte(""), 0},
	{[]byte("a"), nil, 1},
	{nil, []byte("a"), -1},
	//大小写测试
	{[]byte("A"), []byte("a"), 1},
	{[]byte("!"), []byte("！"), 1},
}

func TestCompare(t *testing.T) { //notice:nil==empty slice
	for _, tt := range compareTests {
		cmp := bytes.Compare(tt.a, tt.b)
		if cmp != tt.i {
			t.Errorf(`Compare(%q, %q) = %v`, tt.a, tt.b, cmp)
		}
	}
}

func TestEqual(t *testing.T) {
	for _, tt := range compareTests {
		cmp := bytes.Equal(tt.a, tt.b)
		if cmp && tt.i != 0 {
			t.Errorf(`equal(%q, %q) = %v`, tt.a, tt.b, cmp)
		}
	}
}

func TestEqualFold(t *testing.T) {
	for _, tt := range compareTests {
		cmp := bytes.EqualFold(tt.a, tt.b)
		if false && tt.i == 0 {
			t.Errorf(`equal(%q, %q) = %v`, tt.a, tt.b, cmp)
		}
	}
}

var runes = []struct {
	a []byte
	b []rune
	i int
}{
	{[]byte("abcd"), []rune("abcd"), 0},
	{[]byte("Abcd"), []rune("abcd"), 1},
	{[]byte("~abcd!_"), []rune("~abcd!_"), 0},
	{[]byte("中文测试"), []rune("中文测试"), 0},
	{[]byte("复杂文字测试②のⅡ🔤"), []rune("复杂文字测试②のⅡ🔤"), 0},
	{[]byte("	 abba"), []rune("	 abba"), 0},
}

func TestRunes(t *testing.T) {
	for _, tt := range runes {
		c := bytes.Runes(tt.a)
		if len(tt.b) != len(c) && tt.i == 0 {
			t.Errorf(`equal(%q, %q) = %v`, tt.a, tt.b, c)
		}
		for i := 0; i < len(c); i++ {
			if tt.b[i] != c[i] && tt.i == 0 {
				t.Errorf(`equal(%q, %q) = %v`, tt.a, tt.b, c)
			}
		}
	}
}

var testcount = []struct {
	a []byte
	b []byte
	i int
}{
	{[]byte("abcdefg"), []byte("abc"), 1},
	{[]byte("复杂文字测试②のⅡ🔤哈哈哈"), []byte("复杂文字测试②のⅡ🔤"), 1},
	{[]byte("ababaaba"), []byte("aba"), 2},
}

func TestCount(t *testing.T) {
	for _, tt := range testcount {
		i := bytes.Count(tt.a, tt.b)
		if i != tt.i {
			t.Errorf(`equal(%q, %q) = %v`, tt.a, tt.b, i)
		}
	}
}

func TestIndex(t *testing.T) {
	i := bytes.IndexAny([]byte("a复杂文字测试②のⅡ🔤a"), "🔤")
	if i != 28 {
		t.Errorf("test IndexAny:%d", i)
	}
	if i != 28 {
		t.Errorf("test IndexAny:%d", i)
	}
	i = bytes.IndexFunc([]byte("a复杂文字测试②のⅡ🔤a"), func(r rune) bool {
		if r != rune('🔤') {
			return false
		}
		return true
	})
	if i != 28 {
		t.Errorf("test IndexAny:%d", i)
	}
}

func TestTitle(t *testing.T) {
	println(string(bytes.Title([]byte("abcA!！测试"))))
	println(string(bytes.ToTitle([]byte("abcA!！测试"))))
}

func TestToLowerSpecial(t *testing.T) {
	println("tolower:", string(bytes.ToLowerSpecial(unicode.AzeriCase, []byte("ABCDEFG"))))
	println("totitle:", string(bytes.ToTitleSpecial(unicode.AzeriCase, []byte("abcdefg"))))
	println("toupper:", string(bytes.ToUpperSpecial(unicode.AzeriCase, []byte("abcdefg"))))
}

func TestRepeatAndReplace(t *testing.T) {
	println("repeat:", string(bytes.Repeat([]byte("a!~"), 5)))
	println("replace:", string(bytes.Replace([]byte("a!~aba~!!@cab"), []byte("ab"), []byte("ef"), -1)))
}

func a(x int) int {
	if i := rand.Intn(x); i > 50 {
		return i
	}
	return a(x)
}
func TestMap(t *testing.T) {
	println("Map:", string(bytes.Map(func(r rune) rune {
		return rune(a(int(r)))
	}, []byte("abcedef"))))
}

func TestTrim(t *testing.T) {
	s := []byte("abcdefg测试试")
	a := bytes.Trim(s, "试")
	a[0] = 'b'
	println(string(a), string(s))           //会重复除去直到遇到不同的字符串为止
	println(len(a), cap(a), len(s), cap(s)) //s,a使用同一个切片

	//
	s = []byte("abcdefg测")
	a = bytes.Trim(s, string([]byte{byte(230), byte(181), byte(139)}))
	a[0] = 'b'
	println(string(a))                      //bytes里Trim开头的是转换成rune来去掉字符的
	println(len(a), cap(a), len(s), cap(s)) //s,a使用同一个切片
}

func TestFields(t *testing.T) {
	c := []string{"speace", "tab", "enter", "全角空格", "end"}
	s := bytes.Fields([]byte(`speace tab	enter
		全角空格　end`))
	for i, b := range s {
		println("string:", string(b), "\tbyte:", b, "\tequal:", c[i] == string(b))
	}
}

func TestSplit(t *testing.T) {
	s := bytes.Split([]byte("测试split测试分割测试"), []byte("测试"))
	for i, b := range s {
		println(i, string(b))
	}
	//
	println("========")
	s = bytes.SplitAfter([]byte("测试split测试分割测试"), []byte("测试"))
	for i, b := range s {
		println(i, string(b))
	}
}

func TestReader(t *testing.T) {
	reader := bytes.NewReader([]byte("New reader"))
	fmt.Println("len:", reader.Len())                                           //未读字节长度
	fmt.Println(reader.ReadByte())                                              //读取一个字节
	fmt.Println(reader.Len(), "unreadByte:", reader.UnreadByte(), reader.Len()) //未读长度+1,reader未读时调用会报错

	//
	fmt.Println("=======")
	reader = bytes.NewReader([]byte("测试New reader"))
	fmt.Println("len:", reader.Len())                                           //未读字节长度
	fmt.Println(reader.ReadRune())                                              //读取一个unicode字符
	fmt.Println(reader.Len(), "unreadByte:", reader.UnreadRune(), reader.Len()) //未读读取一个unicode字符+1,reader未读时调用会报错

	//
	fmt.Println("===seek====")
	fmt.Println(reader.Seek(3, io.SeekCurrent)) //args1:设置偏移量,args12:相对值开头、当前、结尾,返回绝对偏移量
	fmt.Println("len:", reader.Len())
	//
	fmt.Println("===readeat===")
	tmp := make([]byte, 20)
	fmt.Println(reader.ReadAt(tmp, 3)) //从args2位置读取最多len(tmp)字节到tmp中
	fmt.Println("len:", reader.Len(), tmp)

	fmt.Println("====writeto===")
	fmt.Println(reader.WriteTo(os.Stdout)) //读取reader到writer中
}

func TestBuffer(t *testing.T) {
	//buffer实现了reader，使用方式同上
	buffer := new(bytes.Buffer) //零值已可用
	buffer.WriteString("测试buffer")
	fmt.Println("len:", buffer.Len(), "cap:", buffer.Cap())
	buffer.WriteString("中文比较占字节，哈哈哈哈哈哈哈哈哈a") //刚好64byte
	buffer.Grow(1)                           //当buffer没耗尽时，调用不会增加cap
	fmt.Println("len:", buffer.Len(), "cap:", buffer.Cap())
	//
	fmt.Println("====Next====")
	a := buffer.Next(60)
	a[0], a[1], a[2] = 51, 51, 51
	c := (*Buffer)(unsafe.Pointer(buffer))
	fmt.Println(string(c.buf), string(a)) //Next返回的数组是公用buffer的数组
	//
	fmt.Println("====ReadFrom=====")
	reader := bytes.NewBuffer([]byte("reader"))
	fmt.Println(buffer.ReadFrom(reader)) //在buffer后添加,将已读的buf删掉
	fmt.Println(buffer.String(), string(c.buf))
}

type Buffer struct {
	buf      []byte // contents are the bytes buf[off : len(buf)]
	off      int    // read at &buf[off], write at &buf[len(buf)]
	lastRead int    // last read operation, so that Unread* can work correctly.
	// FIXME: lastRead can fit in a single byte

	// memory to hold first slice; helps small buffers avoid allocation.
	// FIXME: it would be advisable to align Buffer to cachelines to avoid false
	// sharing.
	bootstrap [64]byte
}
