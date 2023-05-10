package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now().Format(time.Layout)
	fmt.Println(t)
}
