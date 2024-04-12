package main

import "os"

func main() {
	config, err := NewConfig(os.Args)

	if err != nil {
		panic(err)
	}

	app := NewApp(config)
	app.run()
}
