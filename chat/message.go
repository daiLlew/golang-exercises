package main

import "time"

// Represents a single message.
type message struct {
	Name    string
	Message string
	When    time.Time
	WhenStr string
}
