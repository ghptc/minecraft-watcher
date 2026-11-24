package main

import "testing"

func TestParsePlayerJoin(t *testing.T) {
	line := "[20:00:01 INFO]: System chat: Pufe joined the game"

	event, ok := parseLine(line)

	if !ok {
		t.Fatal("expected a valid event")
	}

	if event.Type != "player_join" {
		t.Errorf("unexpected event type: %s", event.Type)
	}

	if event.Message != "Pufe joined the game" {
		t.Errorf("unexpected message: %s", event.Message)
	}
}

func TestParsePlayerLeave(t *testing.T) {
	line := "[20:00:05 INFO]: System chat: Pufe left the game"

	event, ok := parseLine(line)

	if !ok || event.Type != "player_leave" {
		t.Fatal("expected player_leave event")
	}

	if event.Message != "Pufe left the game" {
		t.Errorf("unexpected message: %s", event.Message)
	}
}

func TestParseUnauthorizedAccess(t *testing.T) {
	line := "[20:00:10 INFO]: Disconnecting UnknownPlayer: You are not white-listed on this server!"

	event, ok := parseLine(line)

	if !ok || event.Type != "unauthorized_access" {
		t.Fatal("expected unauthorized_access event")
	}

	expected := "🚨 Disconnecting UnknownPlayer: You are not white-listed on this server!"

	if event.Message != expected {
		t.Errorf("unexpected message: %s", event.Message)
	}
}

func TestParseIgnoredLog(t *testing.T) {
	line := "[20:00:15 INFO]: Saving worlds"

	_, ok := parseLine(line)

	if ok {
		t.Fatal("expected log to be ignored")
	}
}

func TestParseInvalidUsername(t *testing.T) {
	line := "[20:00:01 INFO]: System chat: @everyone joined the game"

	_, ok := parseLine(line)

	if ok {
		t.Fatal("invalid username should be rejected")
	}
}

func TestParseUnauthorizedAccessWithIP(t *testing.T) {
	line := "[20:00:10 INFO]: Disconnecting UnknownPlayer[/192.0.2.10:38478]: You are not white-listed on this server!"

	event, ok := parseLine(line)

	if !ok {
		t.Fatal("expected unauthorized access event")
	}

	expected := "🚨 Disconnecting UnknownPlayer: You are not white-listed on this server!"

	if event.Message != expected {
		t.Errorf("unexpected message: %s", event.Message)
	}
}
