package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var err error
	errorLogger := log.New(os.Stderr, "[ERROR] ", 0)

	flag.Usage = func() {
		fmt.Printf("%s <volume-name> <rows> <columns>\n", os.Args[0])
	}

	// Parse commmand
	var volumeName string
	var rows int
	var columns int

	flag.Parse()
	if flag.NArg() != 3 {
		flag.Usage()
		return
	}

	volumeName = flag.Arg(0)
	rows, err = strconv.Atoi(flag.Arg(1))
	if err != nil {
		flag.Usage()
		return
	}
	columns, err = strconv.Atoi(flag.Arg(2))
	if err != nil {
		flag.Usage()
		return
	}

	fileName := "newdisk.txt"
	outFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		errorLogger.Fatalln(err)
	}
	defer outFile.Close()

	// Create new disk
	_, err = fmt.Fprintf(outFile, "XX: ")
	if err != nil {
		errorLogger.Fatalln(err)
	}

	for i := 1; i < columns; i++ {
		if i%0x10 == 0 {
			_, err = fmt.Fprintf(outFile, "%X", i/0x10)
		} else {
			_, err = fmt.Fprintf(outFile, "%s", " ")
		}
		if err != nil {
			errorLogger.Fatalln(err)
		}
	}

	_, err = fmt.Fprintf(outFile, "%s", "\nXX:")
	if err != nil {
		errorLogger.Fatalln(err)
	}

	for i := 0; i < columns; i++ {
		_, err = fmt.Fprintf(outFile, "%X", i%0x10)
		if err != nil {
			errorLogger.Fatalln(err)
		}
	}

	_, err = fmt.Fprintf(outFile, "%s", "\n")
	if err != nil {
		errorLogger.Fatalln(err)
	}

	encodedVolumeName := hex.EncodeToString([]byte(volumeName))

	_, err = fmt.Fprintf(outFile, "%02X:0010000%s%s\n", 0x0, encodedVolumeName, strings.Repeat("0", columns-len(encodedVolumeName)-7))
	if err != nil {
		errorLogger.Fatalln(err)
	}

	for i := 1; i < rows-1; i++ {
		_, err = fmt.Fprintf(outFile, "%02X:1%02X%s\n", i, i+1, strings.Repeat("0", columns-3))
		if err != nil {
			errorLogger.Fatalln(err)
		}
	}

	_, err = fmt.Fprintf(outFile, "%02X:100%s\n", rows-1, strings.Repeat("0", columns-3))
	if err != nil {
		errorLogger.Fatalln(err)
	}
}
