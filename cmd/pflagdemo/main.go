package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
)

type host struct {
	value string
}

func (h *host) String() string {
	return h.value
}

func (h *host) Set(v string) error {
	h.value = v
	return nil
}

func (h *host) Type() string {
	return "host"
}

func main() {
	flagset := pflag.NewFlagSet("test", pflag.ExitOnError)

	var ip = flagset.IntP("ip", "i", 1234, "help message for ip")

	var boolVar bool
	flagset.BoolVarP(&boolVar, "boolVar", "b", true, "help message for boolVar")

	var h host
	flagset.VarP(&h, "host", "H", "help message for host")

	flagset.SortFlags = false

	err := flagset.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		return
	}

	fmt.Printf("ip: %d\n", *ip)
	fmt.Printf("boolVar: %t\n", boolVar)
	fmt.Printf("host: %+v\n", h)

	i, err := flagset.GetInt("ip")
	getBool, err := flagset.GetBool("boolVar")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("i: %d,%v err: %v\n", i, getBool, err)
}
