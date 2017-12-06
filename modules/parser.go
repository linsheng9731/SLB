package modules

import (
	"strings"
	"log"
)

// schemaHost splits a 'schema://host/path' prefix
// into 'schema', 'host' and '/path'
func schemaHost(prefix string) (schema string, host string) {
	var s,h, p string
	var ll []string
	l := strings.Split(prefix, "://")

	if len(l) > 2  {
		log.Print(prefix)
		log.Fatal(" the host path is invalid!")
	} else if len(l) == 2  {
		s  = l[0]
		l =  l[1:]
	}   else {
		s = "http"
	}

	ll = strings.Split(l[0], "/")
	if len(ll) == 2 {
		h, p = ll[0], ll[1]
	} else if len(ll) == 1 {
		h, p = ll[0], ""
	} else {
		log.Print(prefix)
		log.Fatal(" the host path is invalid!")
	}

	if p == "" {
		p = "/"
	}
	return s, h
}

