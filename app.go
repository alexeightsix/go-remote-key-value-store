package main

import (
	"strconv"
	"time"
)

type app struct {
	store  *store
	log    Log
	config config
	server *server
}

func NewApp(config config) app {
	app := app{}
	app.config = config
	return app
}

func (app *app) shutdown() {
	app.server.ln.Close()
}

func (app *app) run() {
	app.log.Notice("Starting Application...")

	config := app.config

	app.log.Notice("Found " + config.db.Name())

	store, n, err := NewStore(STORE_TYPE_MAP, config.db)

	app.store = store

	if err != nil {
		panic(err)
	}

	app.log.Notice(strconv.FormatInt(int64(n), 10) + " record(s) imported")

	var port int64 = 1337
	hostname := "127.0.0.1"
	read_timeout := time.Second * 5000

	server, err := NewServer(hostname, port, read_timeout)

	app.server = server

	defer server.ln.Close()

	if err != nil {
		panic(err)
	}

	app.log.Notice("Listening for new connections on " + server.addr + ":" + strconv.FormatInt(int64(port), 10))

	for {
		c, err := server.ln.Accept()

		server.total_connections++

		if err != nil {
			app.log.Error(err.Error())
			continue
		}

		go server.handleConnection(c, store)
	}
}
