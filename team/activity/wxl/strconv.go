package main

import (
	"fmt"
	"strconv"
)

func main() {

}

func Parse() {

	b := []byte("bool:")
	b = strconv.AppendBool(b, true)
	fmt.Println(string(b))

	b10 := []byte("int:")
	b10 = strconv.AppendInt(b10, -42, 10)
	fmt.Println(string(b10))

	const v = true
	s := strconv.FormatBool(v)
	fmt.Printf("bool to string :%s\n", s)

	i64 := int64(-42)

	s10 := strconv.FormatInt(i64, 10)
	fmt.Printf("int to string: %s", s10)

	const intString = "12"
	res, err := strconv.Atoi(intString)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("toint: %d\n", res)

	const i = 10
	st := strconv.Itoa(i)
	fmt.Printf("int to string: %s", st)

	const hex = "2f"
	res64, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("toint: %d\n", res64)

	const floatString = "12.3"
	resF, err := strconv.ParseFloat(floatString, 32)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("tofloat: %.5f\n", resF)

	const bstring = "true"
	sbool, err := strconv.ParseBool(bstring)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("tobool: %v\n", sbool)
}
