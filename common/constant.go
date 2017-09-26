package common

import "errors"

const (
	RELOAD = iota
)

var (
	ErrNoFrontend  = errors.New("No frontend configuration detected")
	ErrNoBackend   = errors.New("No backend configuration detected")
	ErrTimeout     = errors.New("Timeout")
	ErrPortExists  = errors.New("Port already in use")
	ErrRouteExists = errors.New("Route already in use")
)
