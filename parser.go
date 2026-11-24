package main

import (
	"regexp"
	"strings"
)

type Event struct {
	Type    string
	Message string
}

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)

func validUsername(username string) bool {
	return usernamePattern.MatchString(username)
}

func parseLine(line string) (Event, bool) {
	_, message, found := strings.Cut(line, "]: ")

	if !found {
		return Event{}, false
	}

	message = strings.TrimSpace(message)
	message = strings.TrimPrefix(message, "System chat: ")

	switch {
	case strings.HasSuffix(message, " joined the game"):
		player := strings.TrimSuffix(message, " joined the game")

		if !validUsername(player) {
			return Event{}, false
		}

		return Event{"player_join", message}, true

	case strings.HasSuffix(message, " left the game"):
		player := strings.TrimSuffix(message, " left the game")

		if !validUsername(player) {
			return Event{}, false
		}

		return Event{"player_leave", message}, true

	case strings.Contains(message, " has made the advancement "):
		player, _, _ := strings.Cut(message, " has made the advancement ")

		if !validUsername(player) {
			return Event{}, false
		}

		return Event{"advancement", message}, true

	case strings.HasPrefix(message, "Disconnecting ") &&
		(strings.HasSuffix(message, ": You are not white-listed on this server!") ||
			strings.HasSuffix(message, ": You are not whitelisted on this server!")):

		player := strings.TrimPrefix(message, "Disconnecting ")
		player, _, _ = strings.Cut(player, "[")
		player, _, _ = strings.Cut(player, ":")
		player = strings.TrimSpace(player)

		if !validUsername(player) {
			return Event{}, false
		}

		return Event{
			"unauthorized_access",
			"🚨 Disconnecting " + player +
				": You are not white-listed on this server!",
		}, true

	default:
		return Event{}, false
	}
}
