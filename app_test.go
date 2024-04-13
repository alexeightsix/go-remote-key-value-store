package main

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func ClientFactory() (net.Conn, error) {
	conn, err := net.Dial("tcp", ":1337")
	return conn, err
}

func AppFactory() app {
	store, _ := StoreFactory(0)
	argv := []string{"/something/foo", "--database=" + store.db.Name()}
	config, _ := NewConfig(argv)

	app := NewApp(config)

	go func() {
		app.run()
	}()

	for {
		conn, err := net.Dial("tcp", ":1337")
		if err == nil {
			conn.Close()
			break
		}
	}
	return app
}

func WriteReadAndClose(req []byte) string {
	conn, _ := ClientFactory()
	defer conn.Close()

	conn.Write(req)

	buf := make([]byte, 1014)
	var n int
	for {
		n, _ = conn.Read(buf)
		break
	}
	return string(buf[0:n])
}

func TestAppErr(t *testing.T) {
	app := AppFactory()
	req := []byte("FOO\nkey\nvalue\n\000\n")
	res := WriteReadAndClose(req)
	assert.Equal(t, "INVALID_METHOD", res)
	app.shutdown()
}

func TestAppSet(t *testing.T) {
	app := AppFactory()
	req := []byte("SET\nkey\nvalue\n\000\n")
	res := WriteReadAndClose(req)
	assert.Equal(t, "OK", res)
	app.shutdown()
}

func TestAppGet(t *testing.T) {
	app := AppFactory()
	app.store.set("hello", "world")
	req := []byte("GET\nhello\n\000\n")
	res := WriteReadAndClose(req)
	assert.Equal(t, "world", res)
	app.shutdown()
}

func TestAppDel(t *testing.T) {
	app := AppFactory()
	app.store.set("hello", "world")
	req := []byte("DEL\nhello\n\000\n")
	res := WriteReadAndClose(req)
	assert.Equal(t, "OK", res)
	app.shutdown()
}
