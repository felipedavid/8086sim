package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fatalError("Usage: ./%s <binary_image>", os.Args[0])
	}

	stream, err := os.ReadFile(os.Args[1])
	if err != nil {
		fatalError("Unable to read file.")
	}

	disassemble(stream)
}

func fatalError(fmtStr string, msg ...any) {
	if len(msg) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, fmtStr, msg...)
		os.Exit(-1)
	}
	fmt.Printf(fmtStr)
	os.Exit(-1)
}
