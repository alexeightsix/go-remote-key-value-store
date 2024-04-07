package main

import (
	"errors"
)

type MapStore struct {
	nodes map[string]node
}

func (s *MapStore) init() {
	s.nodes = make(map[string]node)
}

func (s *MapStore) has(key string) bool {
	_, exists := s.nodes[key]
	return exists
}

func (s *MapStore) set(key, value string) error {
	s.nodes[key] = node{key, value}
	return nil
}

func (s *MapStore) get(key string) (node, error) {
	if s.has(key) {
		return s.nodes[key], nil
	}

	return node{}, errors.New("Key does not exists")
}

func (s *MapStore) delete(key string) (bool, error) {
	if s.has(key) {
		delete(s.nodes, key)
		return true, nil
	}
	return false, errors.New("Key does not exists")
}
