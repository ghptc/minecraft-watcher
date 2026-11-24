package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	maxAttempts = 3
	retryDelay  = 2 * time.Second
)

type DiscordMessage struct {
	Content         string `json:"content"`
	AllowedMentions struct {
		Parse []string `json:"parse"`
	} `json:"allowed_mentions"`
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func sendDiscord(url, content string) error {
	message := DiscordMessage{
		Content: content,
	}

	message.AllowedMentions.Parse = []string{}

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		response, err := client.Post(
			url,
			"application/json",
			bytes.NewReader(payload),
		)

		delay := retryDelay

		if err == nil {
			status := response.StatusCode
			delay = retryAfter(response)

			response.Body.Close()

			switch {
			case status >= 200 && status < 300:
				return nil

			case status == 429:
				fmt.Println("Discord rate limit reached")

			case status >= 500:
				fmt.Println("Discord temporarily unavailable")

			default:
				return fmt.Errorf("discord status: %d", status)
			}
		} else {
			fmt.Println("Discord connection failed")
		}

		if attempt < maxAttempts {
			fmt.Printf("Retrying Discord request (%d/%d)\n",
				attempt+1, maxAttempts)

			time.Sleep(delay)
		}
	}

	return fmt.Errorf("discord delivery failed after %d attempts",
		maxAttempts)
}

func retryAfter(response *http.Response) time.Duration {
	value := response.Header.Get("Retry-After")

	seconds, err := strconv.ParseFloat(value, 64)

	if err != nil || seconds <= 0 {
		return retryDelay
	}

	return time.Duration(seconds * float64(time.Second))
}
