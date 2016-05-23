package main

import (
	"strings"

	"github.com/spf13/cobra"
)

var portNumber int
var verbose bool

func parseFileHostLocation(loc string) (host, path string) {
	path = "."
	host = loc
	sp1 := strings.Split(loc, ":")

	if len(sp1) == 2 {
		host = sp1[0]
		path = sp1[1]
	}
	return
}

type Request struct {
	Type  string
	Value string
}

func main() {
	rootCommand := &cobra.Command{
		Use:   "ft",
		Short: "Folder transfer.",
		Long: `
A program that can either serve a directory, or download a directory from a
coresponding server.
`,
	}

	serveCommand := &cobra.Command{
		Run:   Serve,
		Use:   "serve [<path/to/dir]",
		Short: "Serve all the files in a given directory (default is current directory).",
		Long: `
Examples:
 $ ft serve 
 $ ft serve Downloads
`,
	}

	getCommand := &cobra.Command{
		Run:   Get,
		Use:   "get <remote>:path/to/dir [destination]",
		Short: "Get all the files in a given directory.",
		Long: `
Examples:
 $ ft get remote.host.com:Downloads
`,
	}

	lsCommand := &cobra.Command{
		Run:   List,
		Use:   "ls <remote>:path/to/dir",
		Short: "List all the files in a given directory.",
		Long: `
Examples:
 $ ft ls remote.host.com:Downloads
`,
	}

	rootCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCommand.PersistentFlags().IntVarP(&portNumber, "port", "p", 9002, "port number")

	rootCommand.AddCommand(serveCommand, lsCommand, getCommand)
	rootCommand.Execute()

	return
}
