package main

import (
	"fmt"
	"regexp"
)

func main() {
	reYealinkFile := regexp.MustCompile(`^([0-9a-fA-F]{12})(?:\.y[0-9A-F]+)?\.(?:cfg|xml)$`)
	fmt.Println(reYealinkFile.FindStringSubmatch("001565aabbcc.cfg"))
}

// http.ListenAndServe(":8080", router.Handler())
// fmt.Println("Server listen on port :8080")
