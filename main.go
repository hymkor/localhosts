package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var output = flag.String("o", "", "output filename [default:Stdout]")

func makeIpList() (map[string][]net.IP, error) {
	ints, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ipDict := make(map[string][]net.IP)
	for _, int1 := range ints {
		addrs, err := int1.Addrs()
		if err == nil &&
			(int1.Flags&net.FlagUp) != 0 &&
			(int1.Flags&net.FlagLoopback) == 0 {

			ipList := make([]net.IP, 0)
			for _, addr1 := range addrs {
				ip, _, err := net.ParseCIDR(addr1.String())
				if err == nil && ip.To4() != nil {
					ipList = append(ipList, ip)
				}
			}
			if len(ipList) >= 1 {
				ipDict[int1.Name] = ipList
			}
		}
	}
	return ipDict, nil
}

func main1(patterns []string, out io.Writer) error {
	ipDict, err := makeIpList()
	if err != nil {
		return err
	}
	if len(patterns) < 1 {
		for key, vals := range ipDict {
			fmt.Fprint(out, key)
			for _, val1 := range vals {
				fmt.Fprintf(out, ";%s", val1.String())
			}
			fmt.Fprintln(out)
		}
	} else {
		for _, name := range patterns {
			for key, vals := range ipDict {
				if strings.Contains(key, name) {
					for _, val1 := range vals {
						fmt.Fprintln(out, val1.String())
					}
				}
			}
		}
	}
	return nil
}

func main() {
	flag.Parse()
	var w io.Writer = os.Stdout
	if output != nil && *output != "" {
		fd, err := os.Create(*output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
		defer fd.Close()
		w = fd
	}
	if err := main1(flag.Args(), w); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
