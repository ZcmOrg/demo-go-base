package main

import (
	"fmt"
	"regexp"
	"strings"
)

func StringsCon() {
	refString := "Hello world My darling"

	lookFor := "world"
	contain := strings.Contains(refString, lookFor)
	fmt.Printf("The \"%s\" contains \"%s\": %t \n", refString, lookFor, contain)

	lookFor = "wolf"
	contain = strings.Contains(refString, lookFor)
	fmt.Printf("The \"%s\" contains \"%s\": %t \n", refString, lookFor, contain)

	startsWith := "Hello"
	starts := strings.HasPrefix(refString, startsWith)
	fmt.Printf("The \"%s\" starts with \"%s\": %t \n", refString, startsWith, starts)

	endWith := "darling"
	ends := strings.HasSuffix(refString, endWith)
	fmt.Printf("The \"%s\" ends with \"%s\": %t \n", refString, endWith, ends)
}

func StringS() {
	refString := "darling_has a little_sad"
	words := strings.Split(refString, "_")
	for idx, word := range words {
		fmt.Printf("Word %d is: %s\n", idx, word)
	}

	refString = "darling*has,a,%little_sad"
	words = regexp.MustCompile("[*,%_]{1}").Split(refString, -1)
	for idx, word := range words {
		fmt.Printf("Word %d is: %s\n", idx, word)
	}

	splitFunc := func(r rune) bool {
		return strings.ContainsRune("*%,_", r)
	}

	words = strings.FieldsFunc(refString, splitFunc)
	for idx, word := range words {
		fmt.Printf("Word %d is: %s\n", idx, word)
	}

	refString = "darling has a little sad"
	words = strings.Fields(refString)
	for idx, word := range words {
		fmt.Printf("Word %d is: %s\n", idx, word)
	}
}

func JoinSql() {

	selectBase := "SELECT * FROM users WHERE %s"

	var refStringSlice = []string{
		" nick_name = 'Jack' ",
		" account = 333444555 "}
	sentence := strings.Join(refStringSlice, "AND")
	fmt.Printf(selectBase+"\n", sentence)
}

func regexpReplace() {
	refString := "darling has a little sad"
	regex := regexp.MustCompile("l[a-z]+")
	out := regex.ReplaceAllString(refString, "big")
	fmt.Println(out)
}

func stringReplace() {
	refString := "darling has a little sad"
	out := strings.Replace(refString, "little", "big", -1)
	fmt.Println(out)
}

func stringReplacer() {
	refString := "darling has a little sad"
	replacer := strings.NewReplacer("little", "big", "darling", "BatMan")
	out := replacer.Replace(refString)
	fmt.Println(out)
}

func toCamelCase(input string) string {
	titleSpace := strings.Title(strings.Replace(input, "_", " ", -1))
	camel := strings.Replace(titleSpace, " ", "", -1)
	return strings.ToLower(camel[:1]) + camel[1:]
}

func Snake() {
	email := "ExamPle@domain.com"
	emailToCompare := strings.ToLower(email)
	fmt.Printf("Email lower: %s\n", emailToCompare)

	name := "exl"
	upc := strings.ToUpper(name)
	upt := strings.ToTitle(name)
	fmt.Println("UPPER : " + upc)
	fmt.Println("TITLE : " + upt)
}
