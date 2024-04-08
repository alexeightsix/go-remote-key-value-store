package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	STORE_TYPE_ARRAY = iota
	STORE_TYPE_MAP
)

const (
	ERROR_KEY_NOT_FOUND = "ERROR_KEY_NOT_FOUND"
	ERROR_INVALID_STORE = "ERROR_INVALID_STORE_DRIVER"
)

type store struct {
	mu     sync.Mutex
	driver StoreInterface
	db     *os.File
}

type StoreInterface interface {
	init()
	get(key string) (node, error)
	has(key string) bool
	set(key, value string)
	delete(key string) (bool, error)
}

func NewStore(algo int, db *os.File) (*store, error) {
	store := store{}
	store.mu = sync.Mutex{}
	store.db = db

	switch algo {
	case STORE_TYPE_ARRAY:
		store.driver = &ArrayStore{}
	case STORE_TYPE_MAP:
		store.driver = &MapStore{}
		store.driver.init()
	default:
		return nil, errors.New(ERROR_INVALID_STORE)
	}

	return &store, nil
}

type Action struct {
	pos    time.Duration
	method method
	node   node
}

func (s *store) hydrate() (int, error) {
	scanner := bufio.NewScanner(s.db)

	lp := 0

	for scanner.Scan() {
		z := strings.Split(scanner.Text(), "\t")

		if z[1] != string(SERVER_METHOD_SET) {
			continue
		}

		node := node{}
		node.key = z[2]
		node.value = z[3]

		if s.driver.has(node.key) {
			continue
		}

		s.driver.set(node.key, node.value)

		lp++
	}

	return lp, nil
}

func (s *store) commitSet(key string, value string) (int, error) {
	pos := time.Now().UnixNano()
	str := strconv.FormatInt(pos, 10) + "\t" + string(SERVER_METHOD_SET) + "\t" + key + "\t" + value
	return s.db.Write([]byte(str))
}

func (s *store) commitDel(key string) (int, error) {
	pos := time.Now().UnixNano()
	str := strconv.FormatInt(pos, 10) + "\t" + string(SERVER_METHOD_DEL) + "\t" + key + "\n"
	return s.db.Write([]byte(str))
}

func (s *store) has(key string) bool {
	return s.driver.has(key)
}

func (s *store) set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.driver.set(key, value)
	_, err := s.commitSet(key, value)
	return err
}

func (s *store) get(key string) (node, error) {
	return s.driver.get(key)
}

func (s *store) delete(key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, err := s.driver.delete(key)
	s.commitDel(key)
	return n, err
}
