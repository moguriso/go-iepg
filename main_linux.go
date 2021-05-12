package main

import (
	"flag"
	"fmt"
	"go-iepg/iepg"
	"go-iepg/param"
	"os"
)

var (
	reserveFile = flag.String("r", "", "reserve file name (fullpath)")
)

func main() {

	flag.Parse()
	if _, err := os.Stat(*reserveFile); err != nil {
		*reserveFile = "./reserve.json"
	}

	fmt.Println("file = ", *reserveFile)
	d_conf := param.LoadDynamicParam("reserve.json")
	for _, v := range d_conf {
		iepg.Reserve(v, true) /* Linux doesn't support PLUMAGE */
	}
}
