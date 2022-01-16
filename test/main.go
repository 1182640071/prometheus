package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "  1asdf1 "
	c := strings.Trim(s, " ")
	fmt.Println(c)


}