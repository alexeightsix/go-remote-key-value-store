package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParserIsEnd(t *testing.T) {
	assert.True(t, isEnd([]byte("\000\n")))
	assert.False(t, isEnd([]byte("\001\n")))
}

func TestParserSet(t *testing.T) {
	payload := "SET\nkey\nvalue\n\x00\n"
	bytes := []byte(payload)

	parser := parser{}
	parser.parse(bytes).method().subject().value()
	assert.Nil(t, parser.error)
}

func TestParserDel(t *testing.T) {
	payload := "DEL\nkey\n\x00\n"
	bytes := []byte(payload)

	parser := parser{}
	parser.parse(bytes).method().subject().value()
	assert.Nil(t, parser.error)
}

func TestParserGet(t *testing.T) {
	payload := "DEL\nkey\n\x00\n"
	bytes := []byte(payload)

	parser := parser{}
	parser.parse(bytes).method().subject().value()
	assert.Nil(t, parser.error)
}

func TestParserInvalidMethod(t *testing.T) {
	payload := "FOO\nkey\n\x00\n"
	bytes := []byte(payload)

	parser := parser{}
	parser.parse(bytes).method().subject().value()
	assert.ErrorContains(t, parser.error, "INVALID_METHOD")
}

func TestParserInvalidValue(t *testing.T) {
	payload := "DEL\nkey\n\x00"
	bytes := []byte(payload)

	parser := parser{}
	parser.parse(bytes).method().subject().value()
	assert.ErrorContains(t, parser.error, "INVALID_VALUE")
}
