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

func NewStore(algo int, db *os.File) (*store, int, error) {
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
		return nil, 0, errors.New(ERROR_INVALID_STORE)
	}

	n := store.hydrate()

	return &store, n, nil
}

type action struct {
	time   time.Time
	method method
	node   node
}

func (s *store) readActions(cb func(action action)) {
	scanner := bufio.NewScanner(s.db)

	const (
		COL_TIME       = 0
		COL_ACTION     = 1
		COL_KEY        = 2
		COL_VALUE      = 3
		COL_EXPIRES_AT = 4
	)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")

		t, err := strconv.Atoi(line[COL_TIME])
		timestamp := time.UnixMicro(int64(t))

		if err != nil {
			panic(err)
		}

		action := action{}
		action.time = timestamp

		switch line[COL_ACTION] {
		case "SET":
			action.method = SERVER_METHOD_GET
			action.node = *NewNode(line[COL_KEY], line[COL_VALUE])
			break
		case "DEL":
			action.method = SERVER_METHOD_DEL
			action.node = *NewNode(line[COL_KEY], "")
			break
		default:
			panic("Unable to parse Action")
		}

		cb(action)
	}
}

func (s *store) hydrate() (x int) {
	s.readActions(func(action action) {
		if action.method != SERVER_METHOD_SET && !s.has(action.node.key) {
			s.set(action.node.key, action.node.value)
			x++
		}
	})
	return x
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
