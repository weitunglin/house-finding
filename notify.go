package main

import (
	"fmt"

	"github.com/utahta/go-linenotify"
)

const token string = "LGPKiwIbPmVYEvUww5cZsZqUvR5hh5lzpX1MBCavVHw"

var ln *linenotify.Client

func init() {
	ln = linenotify.New()
}

func notify(msg string) (err error) {
	_, err = ln.Notify(token, msg, "", "", nil)
	if err != nil {
		fmt.Printf("err %v", err.Error())
		return err
	}

	return nil
}
