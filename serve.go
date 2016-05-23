package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Serve(cmd *cobra.Command, args []string) {
	root := ""
	if len(args) > 0 {
		root = args[0]
	}

	listenAddr, rerr := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", portNumber))
	if rerr != nil {
		panic(rerr)
	}

	listener, lerr := net.ListenTCP("tcp", listenAddr)
	if lerr != nil {
		panic(lerr)
	}

	for {
		conn, aerr := listener.Accept()
		if aerr != nil {
			fmt.Println(aerr)
		}
		go handle(root, conn)
	}
}

func handle(root string, c net.Conn) {
	gdec := gob.NewDecoder(c)
	var req Request
	derr := gdec.Decode(&req)
	if derr != nil {
		fmt.Println(derr)
		fmt.Println("Closed conn", c.RemoteAddr())
		c.Close()
		return
	}

	genc := gob.NewEncoder(c)

	type listResult struct {
		Path string
		Size int64
	}

	switch req.Type {
	case "List":
		var results []listResult
		walker := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false {
				relpath, rerr := filepath.Rel(root, path)
				if rerr != nil {
					return rerr
				}
				results = append(results, listResult{Path: relpath, Size: info.Size()})
			}
			return err
		}
		werr := filepath.Walk(filepath.Join(root, req.Value), walker)
		if werr != nil {
			fmt.Println(werr)
		}
		genc.Encode(results)
		c.Close()
		return
	case "Get":
		ifile, ferr := os.Open(filepath.Join(root, req.Value))
		if ferr != nil {
			fmt.Println("Error opening file", ferr)
			fmt.Println("Closed conn", c.RemoteAddr())
			c.Close()
			return
		}
		io.Copy(c, ifile)
		c.Close()
		return
	}

	fmt.Println("unknown request", req.Type)
	fmt.Println("Closed conn", c.RemoteAddr())
	c.Close()
}
