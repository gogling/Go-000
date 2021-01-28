package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println("error dialing", err.Error())
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("type msg to send server : ")
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\r\n")

		_, err = conn.Write([]byte( input))
	}
}
