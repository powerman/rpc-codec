package jsonrpc2

import "log"

func logIfFail(f func() error) {
	err := f()
	if err != nil {
		log.Print(err)
	}
}
