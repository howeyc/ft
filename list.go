package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

func List(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}
	host, path := parseFileHostLocation(args[0])

	remoteAddr := fmt.Sprintf("%s:%d", host, portNumber)

	server, derr := net.Dial("tcp", remoteAddr)
	if derr != nil {
		panic(derr)
	}

	genc := gob.NewEncoder(server)
	genc.Encode(Request{Type: "List", Value: path})

	type listResult struct {
		Path string
		Size int64
	}
	var results []listResult
	gdec := gob.NewDecoder(server)
	gdec.Decode(&results)
	server.Close()

	for _, res := range results {
		fmt.Println(res)
	}

}
