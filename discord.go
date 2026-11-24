package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

	response, err := client.Post(
		url,
		"application/json",
		bytes.NewReader(payload),
	)

	if err != nil {
		return fmt.Errorf("discord request failed")
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("discord status: %d", response.StatusCode)
	}

	return nil
}
