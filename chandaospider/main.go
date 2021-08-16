package main

import (
	"chandaospider/cmd"
	_ "chandaospider/log"
	"fmt"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("%v", err)
	}

}
