package main

import "strings"

type Event struct {
	Type   string
	Player string
}

func parseLine(line string) (Event, bool) {
	_, message, found := strings.Cut(line, "]: System chat: ")

	if !found {
		return Event{}, false
	}

	message = strings.TrimSpace(message)

	if player, ok := strings.CutSuffix(message, " joined the game"); ok {
		return Event{Type: "player_join", Player: player}, true
	}

	if player, ok := strings.CutSuffix(message, " left the game"); ok {
		return Event{Type: "player_leave", Player: player}, true
	}

	return Event{}, false
}
