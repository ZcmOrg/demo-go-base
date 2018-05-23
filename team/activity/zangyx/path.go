package main

import (
	"fmt"
	"path"
)

func main() {
	//IsAbs返回路径是否是一个绝对路径。
	fmt.Println(path.IsAbs("d:/mygo/a.go"))   //false
	fmt.Println(path.IsAbs("d:\\mygo\\a.go")) //true

	// Split函数将路径从最后一个斜杠后面位置分隔为两个部分（dir和file）并返回
	fmt.Println(path.Split("d:/mygo/a.go")) //d:/mygo/ a.go

	//Join函数可以将任意数量的路径元素放入一个单一路径里，会根据需要添加斜杠。结果是经过简化的，所有的空字符串元素会被忽略。
	fmt.Println(path.Join("d:", "aa", "bb", "cc.go")) //c:/aa/bb/cc.txt

	// Dir返回路径除去最后一个路径元素的部分，即该路径最后一个元素所在的目录。
	fmt.Println(path.Dir("d:/mygo/a.go")) //d:/mygo

	// Base函数返回路径的最后一个元素。
	fmt.Println(path.Base("d:/mygo/a.go")) //a.go

	// Ext函数返回path文件扩展名。
	fmt.Println(path.Ext("d:/mygo/a.go")) //.go

	// Clean函数通过单纯的词法操作返回和path代表同一地址的最短路径。
	fmt.Println(path.Ext("d:\\mygo/../a.go")) //.go

	// 如果name匹配shell文件名模式匹配字符串，Match函数返回真。
	m, _ := path.Match("/usr/*", "/usr/local")
	fmt.Println(m) //true
}
