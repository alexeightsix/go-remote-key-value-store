package main

import (
	"os"
	"strconv"
	"time"
)

func main() {
	Log("Initalizing...")

	config, err := NewConfig(os.Args)

	if err != nil {
		panic(err)
	}

	Log("Found " + config.db.Name())

	store, err := NewStore(STORE_TYPE_MAP, config.db)

	n, err := store.hydrate()

	Log(strconv.FormatInt(int64(n), 10) + " record(s) imported")

	if err != nil {
		panic(err)
	}

	var port int64 = 1337

	read_timeout := time.Second * 5000

	server := NewServer("127.0.0.1", port, read_timeout)

	err = server.serve()

	if err != nil {
		panic(err)
	}

	defer server.ln.Close()

	Log("Listening for new connections on 127.0.0.1:" + strconv.FormatInt(int64(port), 10))

	for {
		c, err := server.ln.Accept()

		server.total_connections++

		if err == nil {
			go server.handleConnection(c, store)
		}
	}
}
