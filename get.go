package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/spf13/cobra"
)

func Get(cmd *cobra.Command, args []string) {
	dest := "."
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}
	host, path := parseFileHostLocation(args[0])

	if len(args) > 1 {
		dest = args[1]
	}

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

	// Setup bars
	totalBarIdx := len(results)
	var totalSize int
	bars := make([]*pb.ProgressBar, totalBarIdx+1)
	for idx, res := range results {
		totalSize += int(res.Size)
		bars[idx] = pb.New(int(res.Size)).Prefix(res.Path)
		bars[idx].SetUnits(pb.U_BYTES)
		bars[idx].ShowSpeed = true
	}
	bars[totalBarIdx] = pb.New(totalSize).Prefix("Total")
	bars[totalBarIdx].SetUnits(pb.U_BYTES)
	bars[totalBarIdx].ShowSpeed = true
	pool, _ := pb.StartPool(bars...)

	var wg sync.WaitGroup

	for i, res := range results {
		wg.Add(1)
		go func(path string, size int64, idx int) {
			fulldest := filepath.Join(dest, path)
			fulldir := filepath.Dir(fulldest)
			os.MkdirAll(fulldir, 0777)
			ofile, oerr := os.Create(filepath.Join(dest, path))
			if oerr != nil {
				return
			}

			fserver, fderr := net.Dial("tcp", remoteAddr)
			if fderr != nil {
				panic(fderr)
			}

			ggenc := gob.NewEncoder(fserver)
			ggenc.Encode(Request{Type: "Get", Value: path})

			w := io.MultiWriter(ofile, bars[idx], bars[totalBarIdx])

			io.CopyN(w, fserver, size)
			fserver.Close()
			wg.Done()
		}(res.Path, res.Size, i)
	}
	wg.Wait()

	pool.Stop()
}
