//comperss主要是文件、数据压缩，archive主要是打包

package luom_test

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBzip2(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "e.txt.bz2"))
	reader := bzip2.NewReader(f) //解压缩
	b, err := ioutil.ReadAll(reader)
	fmt.Println(string(b), err)
}

func TestFlate(t *testing.T) {
	buffer := bytes.NewBufferString("")
	w, err := flate.NewWriterDict(buffer, 2, []byte("end la")) //buffer不能有初始值,dict相当于压缩密码
	fmt.Println(err)
	w.Write([]byte("end end end"))
	w.Flush() //会先把writer的内容放进buffer里
	fmt.Println(buffer.Bytes())
	w.Close()
	fmt.Println(buffer.Bytes())
	r := flate.NewReaderDict(buffer, []byte("end la")) //reader、writer同时使用dict，否则可能发生错误
	b, err := ioutil.ReadAll(r)
	fmt.Println(string(b), err)
	r.Close()
}

func TestGzip(t *testing.T) { //gzip 只能压缩单个文件，可以配合tar使用进行打包目录与文件
	buffer := bytes.NewBufferString("") //操作基本跟上面的flate相似
	w := gzip.NewWriter(buffer)
	w.Write([]byte("end end end la"))
	w.Close()
	fmt.Println(buffer.Bytes())

	r, err := gzip.NewReader(buffer)
	fmt.Println(err)
	b, err := ioutil.ReadAll(r)
	fmt.Println(string(b), err)
	r.Close()

	//test ungzip file
	fmt.Println("====test file=====")
	f, err := os.Open("testdata/test.txt.gz")
	fmt.Println(err)
	fr, err := gzip.NewReader(f)
	fmt.Println(err)
	b, err = ioutil.ReadAll(fr)
	fmt.Println(string(b), err)
	fr.Close()
	fmt.Println(fr.Header)
	//test gzip file
	wf, _ := os.Create("testdata/gzipf.txt.gz")
	fw := gzip.NewWriter(wf)
	fw.Write([]byte("hello! golang gzip~"))
	fw.Header = gzip.Header{
		Comment: "",
		Extra:   []byte(""),
		ModTime: time.Now(),
		Name:    "gzipf.txt",
		OS:      3,
	}
	fw.Close()
	wf.Close()
}

func TestLzw(t *testing.T) {
	buffer := bytes.NewBufferString("")
	w := lzw.NewWriter(buffer, lzw.MSB, 8) //读写的order跟litWidth必须一致，否则可能报错或者读取的数据错误
	w.Write([]byte("hello!golang lzw"))
	w.Close()

	r := lzw.NewReader(buffer, lzw.MSB, 8)
	b, err := ioutil.ReadAll(r)
	fmt.Println(string(b), err)
}

func TestZlib(t *testing.T) {
	buffer := bytes.NewBufferString("")
	w, err := zlib.NewWriterLevelDict(buffer, 2, []byte("end la")) //buffer不能有初始值,dict相当于压缩密码
	fmt.Println(err)
	w.Write([]byte("end end end"))
	w.Flush() //会先把writer的内容放进buffer里
	fmt.Println(buffer.Bytes())
	w.Close()
	fmt.Println(buffer.Bytes())
	r, _ := zlib.NewReaderDict(buffer, []byte("end la")) //reader、writer同时使用dict，否则可能发生错误
	b, err := ioutil.ReadAll(r)
	fmt.Println(string(b), err)
	r.Close()
}
