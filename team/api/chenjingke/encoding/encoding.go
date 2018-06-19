package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

/*
1.定义了供其它包使用的可以将数据在字节水平和文本表示之间转换的接口
2.encoding/gob、encoding/json、encoding/xml三个包都会检查使用这些接口。因此，只要实现了这些接口一次，就可以在多个包里使用
3.标准包内建类型time.Time和net.IP都实现了这些接口。接口是成对的，分别产生和还原编码后的数据
*/

// 默认情况下，time.Time 类型按 RFC 3339 格式提供。也就是说，它会是一个字符串，类似于 2016-12-07T17:47:35.099008045-05:00。
type Dog struct {
	ID     int       `json:"id"`
	Name   string    `json:"name"`
	Breed  string    `json:"breed"`
	BornAt time.Time `json:"-"` //"-" 来提醒 JSON 编码器应该忽略这个字段，即使它被导出
}

type JSONDog struct {
	Dog
	BornAt int64 `json:"born_at"`
}

func NewJSONDog(dog Dog) JSONDog {
	return JSONDog{
		dog,
		dog.BornAt.Unix(),
	}
}

func (jd JSONDog) ToDog() Dog {
	return Dog{
		jd.Dog.ID,
		jd.Dog.Name,
		jd.Dog.Breed,
		time.Unix(jd.BornAt, 0),
	}
}

func (jd JSONDog) Dog1() Dog {
	return Dog{
		jd.ID,
		jd.Name,
		jd.Breed,
		time.Unix(jd.BornAt, 0),
	}
}

// 实现 Marshaler 和 Unmarshaler 接口
func (d Dog) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONDog(d))
}

func (d *Dog) UnmarshalJSON(data []byte) error {
	var jd JSONDog
	if err := json.Unmarshal(data, &jd); err != nil {
		return err
	}
	*d = jd.Dog1()
	return nil
}

func main() {

	Json1()
	Json2()
	Json3()
	Xml1()
	Xml2()
	Gob1()
	Gob2()
	Base64()
	Hex1()
	Hex2()
	Hex3()
	Hex4()
	Binary1()
	Binary2()
}

// XML更适合标记文档，JSON更适合数据交互：
// XML是一个完整的标记语言，而JSON不是.
// XML利用标记语言的特性提供了绝佳的延展性（如XPath），在数据存储，扩展及高级检索方面优势明显;而JSON则由于比XML更加小巧，以及浏览器的内建快速解析支持，使得其更适用于网络数据传输领域。
// 就可读性而言，两者都具备很好的可读性，但XML文档的可读性更高。
// 就数据表示和传输性能而言，JSON明显比XML简洁，格式简单，占用带宽少。
/*
各自的例子：
{"employees":[
    { "firstName":"John", "lastName":"Doe" },
    { "firstName":"Anna", "lastName":"Smith" },
    { "firstName":"Peter", "lastName":"Jones" }
]}

<employees>
    <employee>
        <firstName>John</firstName> <lastName>Doe</lastName>
    </employee>
    <employee>
        <firstName>Anna</firstName> <lastName>Smith</lastName>
    </employee>
    <employee>
        <firstName>Peter</firstName> <lastName>Jones</lastName>
    </employee>
</employees>
*/

/*============================= json ===========================*/
/*
		把对象转换为JSON:
　　　　• 布尔型转换为 JSON 后仍是布尔型　， 如true -> true

　　　　• 浮点型和整数型转换后为JSON里面的常规数字，如 1.23 -> 1.23

　　　　• 字符串将以UTF-8编码转化输出为Unicode字符集的字符串，特殊字符比如<将会被转义为\u003c

　　　　• 数组和切片被转换为JSON 里面的数组，[]byte类会被转换为base64编码后的字符串，slice的零值被转换为null

　　　　• 结构体会转化为JSON对象，并且只有结构体里边以大写字母开头的可被导出的字段才会被转化输出，而这些可导出的字段会作为JSON对象的字符串索引

　　　　• 转化一个map 类型的数据结构时，该数据的类型必须是 map[string]T（T 可以是encoding/json 包支持的任意数据类型）
*/
// Marshal 和 Unmarshal
func Json1() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

	var jsonBlob = []byte(`[
        {"Name": "Platypus", "Order": "Monotremata"},
        {"Name": "Quoll",    "Order": "Dasyuromorphia"}
    ]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err = json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
}

// Encoders and Decoders
func Json2() {
	const jsonStream = `
        {"Name": "Ed", "Text": "Knock knock."}
        {"Name": "Sam", "Text": "Who's there?"}
        {"Name": "Ed", "Text": "Go fmt."}
        {"Name": "Sam", "Text": "Go fmt who?"}
        {"Name": "Ed", "Text": "Go fmt yourself!"}
    `
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))
	for {
		var m Message
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", m.Name, m.Text)
	}
}

// Encode和Marshal的区别:Encode是编码器上的一种方法，它将JSON编码的Go类型写入输出流（func NewEncoder需要一个io.Writer并返回一个*编码器）。Marshal是一个返回Go类型的JSON编码的函数。
type Response1 struct {
	Page   int
	Fruits []string
}

type Response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func Json3() {

	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))

	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))

	fltB, _ := json.Marshal(2.34)
	fmt.Println(string(fltB))

	strB, _ := json.Marshal("gopher")
	fmt.Println(string(strB))

	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]int{"apple": 5, "lettuce": 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	res1D := &Response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	res2D := &Response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))

	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)

	var dat map[string]interface{}

	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

	num := dat["num"].(float64)
	fmt.Println(num)

	strs := dat["strs"].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)

	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := Response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])

	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)
}

/*============================= xml ===========================*/
/*
	在进行封装时， XML 元素的名字由一系列规则决定， 这些规则的优先级从高到低依次为：
	1.如果给定的数据是一个结构， 那么使用 XMLName 字段的标签作为元素名
	2.使用类型为 Name 的 XMLName 字段的值为元素名
	3.将用于获取数据的结构字段的标签用作元素名
	4.将用于获取数据的结构字段的名字用作元素名
	5.将被封装类型的名字用作元素名
*/
/*
	结构中的每个已导出字段都会被封装为相应的元素并包含在 XML 里面， 但以下规则中提到的内容除外：
	1.XMLName 字段，因为前面提到的原因，会被忽略
	2.带有 “-” 标签的字段会被忽略
	3.带有 “name,attr” 标签的字段会成为 XML 元素的属性， 其中属性的名字为这里给定的 name
	4.带有 ”,attr” 标签的字段会成为 XML 元素的属性， 其中属性的名字为字段的名字
	5.带有 ”,chardata” 标签的字段将会被封装为字符数据而不是 XML 元素。
	6.带有 ”,cdata” 标签的字段将会被封装为字符数据而不是 XML 元素， 并且这些数据还会被一个或多个
*/
// Marshal 函数 func Marshal(v interface{}) ([]byte, error)
func Xml1() {
	type Address struct {
		City, State string
	}
	type Person struct {
		XMLName   xml.Name `xml:"person"`
		Id        int      `xml:"id,attr"`
		FirstName string   `xml:"name>first"`
		LastName  string   `xml:"name>last"`
		Age       int      `xml:"age"`
		Height    float32  `xml:"height,omitempty"`
		Married   bool
		Address
		Comment string `xml:",comment"`
	}

	v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	v.Comment = " Need more details. "
	v.Address = Address{"Hanga Roa", "Easter Island"}

	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)
	fmt.Println("\n")
}

// Unmarshal 函数 func Unmarshal(data []byte, v interface{}) error
func Xml2() {
	type Email struct {
		Where string `xml:"where,attr"`
		Addr  string
	}
	type Address struct {
		City, State string
	}
	type Result struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"FullName"`
		Phone   string
		Email   []Email
		Groups  []string `xml:"Group>Value"`
		Address
	}
	v := Result{Name: "none", Phone: "none"}

	data := `
        <Person>
            <FullName>Grace R. Emlin</FullName>
            <Company>Example Inc.</Company>
            <Email where="home">
                <Addr>gre@example.com</Addr>
            </Email>
            <Email where='work'>
                <Addr>gre@work.com</Addr>
            </Email>
            <Group>
                <Value>Friends</Value>
                <Value>Squash</Value>
            </Group>
            <City>Hanga Roa</City>
            <State>Easter Island</State>
        </Person>
    `
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Printf("XMLName: %#v\n", v.XMLName)
	fmt.Printf("Name: %q\n", v.Name)
	fmt.Printf("Phone: %q\n", v.Phone)
	fmt.Printf("Email: %v\n", v.Email)
	fmt.Printf("Groups: %v\n", v.Groups)
	fmt.Printf("Address: %v\n", v.Address)
}

// func NewDecoder(r io.Reader) *Decoder:创建一个新的读取 r 的 XML 语法分析器。 如果 r 没有实现 io.ByteReader ， 那么函数将使用它自有的缓冲机制。

// func (d *Decoder) Decode(v interface{}) error:执行与 Unmarshal 一样的解码工作， 唯一的不同在于这个方法会通过读取解码器流来查找起始元素。

// func NewEncoder(w io.Writer) *Encoder:返回一个能够对 w 进行写入的编码器。

// func (enc *Encoder) Encode(v interface{}) error:将 XML 编码的 v 写入到流里面。

/*============================= gob ===========================*/

// gob是Golang包自带的一个数据结构序列化的编码/解码工具。
// 编码使用Encoder，解码使用Decoder。一种典型的应用场景就是RPC(remote procedure calls)。

// func Register(value interface{}):在其内部类型名称下记录一种类型，由该类型的值标识。

// func NewDecoder(r io.Reader) *Decoder:NewDecoder返回一个从io.Reader读取的新解码器。

// func (dec *Decoder) Decode(e interface{}) error:解码从输入流中读取下一个值并将其存储在由空接口值表示的数据中。

// func NewEncoder(w io.Writer) *Encoder:NewEncoder返回一个将在io.Writer上传输的新编码器。

// func (enc *Encoder) Encode(e interface{}) error:编码传输由空接口值表示的数据项，保证所有必需的类型信息先传送。

// gob的优势就是：发送方的结构和接受方的结构并不需要完全一致，例如定义一个结构体：
type MyFace interface {
	A()
}

type Cat struct{}
type Dog2 struct{}

func (c *Cat) A() {
	fmt.Println("Meow")
}

func (d *Dog2) A() {
	fmt.Println("Woof")
}

func init() {
	gob.Register(&Cat{})
	gob.Register(&Dog{})
}

func Gob1() {
	network := new(bytes.Buffer)
	enc := gob.NewEncoder(network)

	var inter MyFace
	inter = new(Cat)

	// Note: pointer to the interface
	err := enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	inter = new(Dog2)
	err = enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	// Now lets get them back out
	dec := gob.NewDecoder(network)

	var get MyFace
	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}

	// Should meow
	get.A()

	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}

	// Should woof
	get.A()
}

// gob包实现的序列化struct对象保存到本地，利用gob反序列化本地的struct对象
const file = "./test.gob"

type User struct {
	Name, Pass string
}

func Gob2() {
	var datato = &User{"Donald", "DuckPass"}
	var datafrom = new(User)

	err := Save(file, datato)
	Check(err)
	err = Load(file, datafrom)
	Check(err)
	fmt.Println(datafrom)
}

func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func Check(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		os.Exit(1)
	}
}

/*============================= base64 ===========================*/

func Base64() {
	data := "abc123!?$*&()'-=@~"

	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)

	sDec, _ := base64.StdEncoding.DecodeString(sEnc)
	fmt.Println(string(sDec))

	uEnc := base64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println(uEnc)

	uDec, _ := base64.URLEncoding.DecodeString(uEnc)
	fmt.Println(string(uDec))

	s := []byte("http://golang.org/pkg/encoding/base64/#variables")
	sDec1 := base64.StdEncoding.EncodeToString(s)
	fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(s))

	uDec1, _ := base64.URLEncoding.DecodeString(sDec1)
	fmt.Println(string(uDec1))
}

/*============================= hex ===========================*/

func Hex1() {
	src := []byte("48656c6c6f20476f7068657221")

	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", dst[:n])
}

func Hex2() {
	const s = "48656c6c6f20476f7068657221"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", decoded)
}

func Hex3() {
	content := []byte("Go is an open source programming language.")

	fmt.Printf("%s", hex.Dump(content))
}

func Hex4() {
	lines := []string{
		"Go is an open source programming language.",
		"\n",
		"We encourage all Go users to subscribe to golang-announce.",
	}

	stdoutDumper := hex.Dumper(os.Stdout)

	defer stdoutDumper.Close()

	for _, line := range lines {
		stdoutDumper.Write([]byte(line))
	}
}

/*============================= binary ===========================*/
/*
	二进制协议:高效地在底层处理数据通信:字节序决定字节输出的顺序、通过可变长度编码压缩数据存储空间.

	基于文本类型的协议（比如 JSON）和二进制协议都是字节通信，他们不同点在于他们使用哪种类型的字节和如何组织这些字节。

	文本协议只适用于 ASCII 或 Unicode 编码可打印的字符通信。例如 "26" 使用 "2" 和 "6" 的 utf 编码的字符串表示，这种方式方便我们读，但对于计算机效率较低。

	在二进制协议中，同样数字 "26" 可使用一个字节 0x1A 十六进制表示，减少了一半的存储空间且原始的字节格式能够被计算机直接识别而不需解析。当一个数字足够大的时候，性能优势就会明显体现。
*/
/*
	计算机字节序和网络字节序:
	字节序:多字节数据类型 (int, float 等)在内存中的存储顺序:
	1.大端序:低地址端存放高位字节:譬如网络传输和文件存储
	2.小端序:低地址端存放低位字节:广泛应用于现代性 CPU 内部存储数据
*/

// 内置的读写固定长度值的流
// Read 通过指定类型的字节序把字节解码 (decode) 到 data 变量中。解码布尔类型时，0 字节 (也就是 []byte{0x00}) 为 false, 其他都为 true
func Binary1() {
	var (
		piVar   float64
		boolVar bool
	)

	piByte := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	boolByte := []byte{0x00}

	piBuffer := bytes.NewReader(piByte)
	boolBuffer := bytes.NewReader(boolByte)

	binary.Read(piBuffer, binary.LittleEndian, &piVar)
	binary.Read(boolBuffer, binary.LittleEndian, &boolByte)

	fmt.Println("pi", piVar)     // pi 3.141592653589793
	fmt.Println("bool", boolVar) // bool false

}

// Write 是 Read 的逆过程
func Binary2() {
	buf := new(bytes.Buffer)
	var pi float64 = math.Pi
	err := binary.Write(buf, binary.LittleEndian, pi)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf.Bytes()) // 18 2d 44 54 fb 21 09 40
}

// 在实际编码中，面对复杂的数据结构，可考虑使用更标准化高效的协议，比如 Protocol Buffer
