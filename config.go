package main

import (
	"errors"
	"os"
	"regexp"
)

type config struct {
	pwd string
	db  *os.File
}

func NewConfig(args []string) (config, error) {
	config := config{}
	config.pwd = args[0]

	if len(args) < 2 {
		return config, errors.New("Config Parameters Missing")
	}

	input := args[1:][0]

	re := regexp.MustCompile(`^--(?P<key>database)=(?P<value>.*)$`)
	matches := re.FindSubmatch([]byte(input))

	if len(matches) == 0 {
		return config, errors.New("--database parameter is missing")
	}

	db_file, err := os.OpenFile(string(matches[2]), os.O_RDWR, 0666)

	if err == nil {
		config.db = db_file
	}

	return config, err
}
