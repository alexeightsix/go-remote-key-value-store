package main

import "time"

type node struct {
	key        string
	value      string
	expires_at time.Time
}

func (n node) isExpired() bool {
	if n.expires_at.IsZero() {
		return false
	}

	return time.Now().UnixMicro() < n.expires_at.UnixMicro()
}

func (n *node) setExpiresAt(time time.Time) *node {
	n.expires_at = time
	return n
}

func NewNode(key string, value string) *node {
	n := node{}
	n.key = key
	n.value = value

	return &n
}
