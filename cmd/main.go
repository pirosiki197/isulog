package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pirosiki197/isulog/internal"
)

const usage = `Usage:
	isulog [-f filename]`

func printUsage() {
	fmt.Println(usage)
}

func main() {
	var (
		filename string
		help     bool
	)
	flag.StringVar(&filename, "f", "isulog.out", "the file name record log")
	flag.BoolVar(&help, "help", false, "print help")
	flag.Parse()

	if help {
		printUsage()
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	parser := internal.NewParser(f)
	parser.Parse()
	parser.Print("sum")
}
