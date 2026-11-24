package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendDiscordSuccess(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	defer server.Close()

	err := sendDiscord(server.URL, "Pufe joined the game")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendDiscordRetry(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++

			if attempts == 1 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	defer server.Close()

	if err := sendDiscord(server.URL, "Test message"); err != nil {
		t.Fatal(err)
	}

	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestSendDiscordRateLimit(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++

			if attempts == 1 {
				w.Header().Set("Retry-After", "0.1")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	defer server.Close()

	if err := sendDiscord(server.URL, "Test message"); err != nil {
		t.Fatal(err)
	}

	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestSendDiscordUnauthorized(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusUnauthorized)
		}),
	)

	defer server.Close()

	if err := sendDiscord(server.URL, "Test message"); err == nil {
		t.Fatal("expected unauthorized error")
	}

	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}
