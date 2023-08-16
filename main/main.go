package main

import (
	"fmt"
	"io/ioutil"
	"mnemosyne/internals"
	"os"
	"strconv"
)

const VERSION = 1

func main() {

	fmt.Println("Mnemosyne v." + strconv.Itoa(1) + "\n")

	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename as an argument.")
		return
	}

	if code, err := ioutil.ReadFile(os.Args[1]); err == nil {
		filename := os.Args[1]
		internals.Start(string(code), filename)
	} else {
		fmt.Println("Cannot read the input file: ", err)
	}

	return

}
