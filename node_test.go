package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExpiredNode(t *testing.T) {
	n := NewNode(
		"key",
		"value",
	)

	now := time.Now()
	then := now.Add(time.Duration(-10) * time.Minute)
	n.setExpiresAt(then)
	assert.True(t, n.isExpired())

	future := time.Duration(time.Hour * 24)
	n.setExpiresAt(time.Now().Add(future))
	assert.False(t, n.isExpired())
}
