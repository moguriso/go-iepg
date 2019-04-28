package main

import (
	"go-iepg/iepg"
	"go-iepg/param"
)

func main() {
	d_conf := param.LoadDynamicParam("reserve.json")
	for _, v := range d_conf {
		iepg.PrintReserve(v) /* Linux doesn't support PLUMAGE */
	}
}
