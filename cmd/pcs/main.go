package main

import (
	"fmt"
	"regexp"
)

func main() {

	ipPattern := regexp.MustCompile(`^((25[0-5]|(2[0-4]|[01]\d|[0-9]|)\d)\.?\b){4}$`)

	fmt.Println(ipPattern.FindString("10.0.0.01"))
	fmt.Println(ipPattern.FindString("001.0.0.1"))
	fmt.Println(ipPattern.FindString("10.0.255.255"))
	fmt.Println(ipPattern.FindString("255.255.255.255"))
	fmt.Println(ipPattern.FindString("255255255255"))

	macPattern := regexp.MustCompile(`^(([a-fA-F0-9]{2}[-:. ]){5}[a-fA-F0-9]{2}|[a-fA-F0-9]{12})$`)

	fmt.Println(macPattern.FindString("aa-ff-45-5c-c7-81"))
	fmt.Println(macPattern.FindString("aa:ff:45:5c:c7:81"))
	fmt.Println(macPattern.FindString("aa.ff.45.5c.c7.81"))
	fmt.Println(macPattern.FindString("aaff455cc781"))
	fmt.Println(macPattern.FindString("aa-ff-45-5c-c7-81"))
	fmt.Println(macPattern.FindString("ag:f33:45:5c:c7:81"))
	fmt.Println(macPattern.FindString("at.334.45.5c.c7.81"))
	fmt.Println(macPattern.FindString("aa ff 45 35 cc 78"))

}
