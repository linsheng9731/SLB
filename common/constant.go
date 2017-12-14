package common

import (
	"errors"
)

const (
	RELOAD = iota
)

const (
	CONFIG_FILENAME = "config.json"
)

var (
	ErrNoFrontend = errors.New("No frontend configuration detected")
	ErrNoBackend  = errors.New("No backend configuration detected")
	ErrTimeout    = errors.New("Timeout")
	ErrPortExists = errors.New("Port already in use")
)
