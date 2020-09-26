package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	input string
)

func main() {
	flag.StringVar(&input, "c", "", "op could be get/set/del")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("my-cli> ")

		line, _, _ := reader.ReadLine()

		flag.CommandLine.Parse([]string{"-c", string(line)})

		if input == "quit" {
			fmt.Println("my-cli quit!")
			return
		}

		if strings.TrimSpace(input) == "" {
			continue
		}

		fmt.Printf(" cmd: %s \n", input)
	}
}
