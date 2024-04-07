package main

import (
	"errors"
)

type ArrayStore struct {
	nodes []node
}

func (s *ArrayStore) init() {
	panic("Not implemented")
}

func (s *ArrayStore) _has(key string) (int, bool) {
	for i := 0; i < len(s.nodes); i++ {
		if s.nodes[i].key == key {
			return i, true
		}
	}
	return -1, false
}

func (s *ArrayStore) has(key string) bool {
	_, has := s._has(key)
	return has
}

func (s *ArrayStore) set(key, value string) error {
	idx, has := s._has(key)

	if has {
		s.nodes[idx] = node{key, value}
	} else {
		s.nodes = append(s.nodes, node{key, value})
	}

	return nil
}

func (s *ArrayStore) get(key string) (node, error) {
	idx, has := s._has(key)

	if has == false {
		return node{}, errors.New("Key not found")
	}

	return s.nodes[idx], nil
}

func (s *ArrayStore) delete(key string) (bool, error) {
	idx, has := s._has(key)

	if has == false {
		return false, errors.New("Key not found")
	}

	tmp := []node{}
	tmp = append(tmp, s.nodes[:idx]...)
	s.nodes = append(tmp, s.nodes[idx+1:]...)

	return true, nil
}
