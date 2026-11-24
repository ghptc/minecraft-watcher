package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	if webhookURL == "" {
		fmt.Println("Missing DISCORD_WEBHOOK_URL")

		return
	}

	logPath := "testdata/latest.log"
	file, err := os.Open(logPath)

	if err != nil {
		fmt.Println("Error opening log:", err)

		return
	}

	defer func() {
		file.Close()
	}()

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		fmt.Println("Error seeking file:", err)

		return
	}

	reader := bufio.NewReader(file)

	fmt.Println("Watching Minecraft logs...")

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			rotated, err := fileRotated(file, logPath)

			if err == nil && rotated {
				newFile, err := os.Open(logPath)

				if err == nil {
					file.Close()
					file = newFile
					reader.Reset(file)
				}
			}

			time.Sleep(500 * time.Millisecond)

			continue
		}

		if err != nil {
			fmt.Println("Error reading log:", err)

			return
		}

		event, ok := parseLine(line)

		if !ok {
			continue
		}

		fmt.Println(event.Message)

		if err := sendDiscord(webhookURL, event.Message); err != nil {
			fmt.Println("Notification error:", err)
		}
	}

}
