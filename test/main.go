package main

import (
	"fmt"
	"time"
)

func main() {
	timeStr := time.Now().Format("2006-01-02 15:04:05") //转化所需模板
	fmt.Println(timeStr)


}