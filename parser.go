package main

import "strings"

type Event struct {
	Type    string
	Message string
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
		return Event{"player_join", message}, true

	case strings.HasSuffix(message, " left the game"):
		return Event{"player_leave", message}, true

	case strings.Contains(message, " has made the advancement "):
		return Event{"advancement", message}, true

	case strings.HasPrefix(message, "Disconnecting ") &&
		(strings.HasSuffix(message, ": You are not white-listed on this server!") ||
			strings.HasSuffix(message, ": You are not whitelisted on this server!")):

		player, _, _ := strings.Cut(
			strings.TrimPrefix(message, "Disconnecting "), ":",
		)

		player, _, _ = strings.Cut(player, "[")

		return Event{
			"unauthorized_access",
			"🚨 Disconnecting " + player +
				": You are not white-listed on this server!",
		}, true

	default:
		return Event{}, false
	}
}
