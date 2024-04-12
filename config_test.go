package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestValidConfig(t *testing.T) {
	store, _ := StoreFactory(0)
	argv := []string{"/something/foo", "--database=" + store.db.Name()}
	_, err := NewConfig(argv)
	assert.Nil(t, err)
}

func TestConfigParametersMissing(t *testing.T) {
	_, err := NewConfig([]string{""})
	assert.ErrorContains(t, err, "Config Parameters Missing")
}

func TestConfigDatabaseParamMissing(t *testing.T) {
	argv := []string{"/something/foo", "-something=\"foo\""}
	_, err := NewConfig(argv)
	assert.ErrorContains(t, err, "--database parameter is missing")
}
