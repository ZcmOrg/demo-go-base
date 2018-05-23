package mime_test

import (
	"fmt"
	"mime"
	"path"
	"testing"
)

//函数将扩展名和mimetype建立偶联；扩展名应以点号开始，例如".html"
func TestAddExtensionType(t *testing.T) {
	a := mime.AddExtensionType(".svg", "image/svg+xml")
	b := mime.AddExtensionType(".m3u8", "application/x-mpegURL")
	c := mime.AddExtensionType(".ts", "video/MP2T")
	println(a)
	fmt.Println(b)
	fmt.Println(c)
}

//函数返回与扩展名偶联的MIME类型。扩展名应以点号开始，如".html"。如果扩展名未偶联类型，函数会返回""。
// 内建的偶联表很小，但在unix系统会从本地系统的一或多个mime.types文件（参加下表）进行增补。
//这里使用的path包的Ext方法，获取路径字符串中的文件扩展名。
func TestTypeByExtension(t *testing.T) {
	fmt.Println("================")
	filepath := "./1.png"
	mimetype := mime.TypeByExtension(path.Ext(filepath))
	fmt.Println(mimetype)

	filepath = "./2.txt"
	mimetype = mime.TypeByExtension(path.Ext(filepath))
	fmt.Println(mimetype)

	filepath = "./3.html"
	mimetype = mime.TypeByExtension(path.Ext(filepath))
	fmt.Println(mimetype)
}

/**
函数根据RFC 1521解析一个媒体类型值以及可能的参数。
媒体类型值一般应为Content-Type和Conten-Disposition头域的值（参见RFC 2183）。
成功的调用会返回小写字母、去空格的媒体类型和一个非空的map。
返回的map映射小写字母的属性和对应的属性值。
*/
func TestParseMediaType(t *testing.T) {
	a := "1111"
	mType, parameters, err := mime.ParseMediaType(a)
	if err != nil {
		fmt.Println("err ", err)
	}
	fmt.Println("Media type : ", mType)
	for param := range parameters {
		fmt.Printf("%v = %v\n\n", param, parameters[param])
	}
	/**
	函数根据RFC 2045和 RFC 2616的规定将媒体类型t和参数param连接为一个mime媒体类型，类型和参数都采用小写字母。
	任一个参数不合法都会返回空字符串。
	*/
	mt := mime.FormatMediaType(mType, parameters)
	fmt.Println("Media : ", mt)
}
