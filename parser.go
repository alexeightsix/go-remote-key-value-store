package main

import (
	"errors"
)

type parser struct {
	error   error
	payload payload
	buf     []byte
}

const (
	NULL_BYTE = 0
	NEW_LINE  = 10
)

func isEnd(buf []byte) bool {
	if buf[len(buf)-1] == NULL_BYTE {
		return true
	}

	if buf[len(buf)-2] == NULL_BYTE && buf[len(buf)-1] == NEW_LINE {
		return true
	}

	return false
}

func (p *parser) parse(buf []byte) *parser {
	p.buf = buf
	return p
}

func (p *parser) method() *parser {
	if p.payload.method != "" {
		return p
	}

	method := p.buf[:4]

	switch string(method) {
	case "SET\n":
		p.payload.method = SERVER_METHOD_SET
		break
	case "GET\n":
		p.payload.method = SERVER_METHOD_GET
		break
	case "DEL\n":
		p.payload.method = SERVER_METHOD_DEL
		break
	default:
		p.error = errors.New("INVALID_METHOD")
		break
	}

	p.buf = p.buf[4:]

	return p
}

func (p *parser) subject() *parser {
	if p.error != nil {
		return p
	}

	var end int

	for i := 0; i < len(p.buf); i++ {
		if i > 25 {
			p.error = errors.New("INVALID_SUBJECT")
			return p
		}

		if p.buf[i] == NEW_LINE {
			end = i
			break
		}
	}

	p.payload.subject = string(p.buf[:end])

	p.buf = p.buf[end+1:]

	return p
}

func (p *parser) value() *parser {
	if p.error != nil {
		return p
	}

	if p.buf[len(p.buf)-1] != NEW_LINE {
		p.error = errors.New("invalid value, value must end with")
		return p
	}

	p.payload.value = string(p.buf[:len(p.buf)-2])
	p.buf = []byte{}

	return p
}
