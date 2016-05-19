package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

type Request struct {
	Type  string
	Value string
}

func handle(c net.Conn) {
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
			if info.IsDir() == false {
				results = append(results, listResult{Path: path, Size: info.Size()})
			}
			return err
		}
		filepath.Walk(req.Value, walker)
		genc.Encode(results)
		c.Close()
		return
	case "Get":
		ifile, ferr := os.Open(req.Value)
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

func main() {
	var port int
	flag.IntVar(&port, "port", 9002, "listen port")
	flag.Parse()

	listenAddr, rerr := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
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
		go handle(conn)
	}
}
