package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"mnemosyne/internals"
	"os"
	"strconv"
)

const VERSION = 1

func main() {

	fmt.Println("\n🅜 🅝 🅔 🅜 🅞 🅢 🅨 🅝 🅔 v." + strconv.Itoa(2) + "\n")

	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename as an argument.")
		return
	}

	kafkaFlag := flag.Bool("kafka", false, "Enable Kafka integration")
	filename := flag.String("file", "", "Preprocessing file")

	flag.Parse()

	//fmt.Println("KAFKA: " + strconv.FormatBool(*kafkaFlag))

	if code, err := ioutil.ReadFile(*filename); err == nil {

		internals.Start(string(code), *filename, *kafkaFlag)
	} else {
		fmt.Println("Cannot read the input file: ", err)
	}

	return

}
