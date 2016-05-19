package main

import (
	"encoding/gob"
	"flag"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/cheggaaa/pb"
)

type Request struct {
	Type  string
	Value string
}

func main() {
	var folder string
	var dest string
	var remoteAddr string
	flag.StringVar(&folder, "folder", ".", "folder to download")
	flag.StringVar(&dest, "dest", ".", "folder to save files to")
	flag.StringVar(&remoteAddr, "remote", "", "server:port to connect to")
	flag.Parse()

	server, derr := net.Dial("tcp", remoteAddr)
	if derr != nil {
		panic(derr)
	}

	genc := gob.NewEncoder(server)
	genc.Encode(Request{Type: "List", Value: folder})

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
			os.MkdirAll(dest, 0777)
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
