package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	file, err := os.Open("testdata/latest.log")

	if err != nil {
		fmt.Println("Error opening log:", err)

		return
	}

	defer file.Close()

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		fmt.Println("Error seeking file:", err)

		return
	}

	reader := bufio.NewReader(file)

	fmt.Println("Watching Minecraft logs...")

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)

			continue
		}

		if err != nil {
			fmt.Println("Error reading log:", err)

			return
		}

		event, ok := parseLine(line)

		if ok {
			fmt.Printf("[%s] %s\n", event.Type, event.Player)
		}
	}
}
