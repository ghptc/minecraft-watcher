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

	logPath := os.Getenv("LOG_PATH")

	if logPath == "" {
    	logPath = "/logs/latest.log"
	}

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

	var pending string

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			pending += line

			rotated, rotationErr := fileRotated(file, logPath)

			if rotationErr == nil && rotated {
				newFile, openErr := os.Open(logPath)

				if openErr == nil {
					file.Close()
					file = newFile
					reader.Reset(file)
					pending = ""
				}
			} else {
				truncated, truncateErr := fileTruncated(file)

				if truncateErr == nil && truncated {
					if _, err := file.Seek(0, io.SeekStart); err != nil {
						fmt.Println("Seek error:", err)
						return
					}

					reader.Reset(file)
					pending = ""
				}
			}

			time.Sleep(500 * time.Millisecond)
			continue
		}

		if err != nil {
			fmt.Println("Error reading log:", err)
			return
		}

		line = pending + line
		pending = ""

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
