package modules

import (
	"testing"
)

func TestHostParse(t *testing.T) {

	var s, h, p string
	prefix1 := "http://127.0.0.1:80/path"
	prefix2 := "https://127.0.0.1/path"
	prefix3 := "127.0.0.1/path"
	prefix4 := "127.0.0.1/"
	prefix5 := "127.0.0.1"
	prefix6 := "http://127.0.0.1"

	s, h = schemaHost(prefix1)
	if s != "http" || h != "127.0.0.1:80"  {
		t.Fatal(s, h, p)
	}

	s, h = schemaHost(prefix2)
	if s != "https" || h != "127.0.0.1"  {
		t.Fatal(s, h, p)
	}

	s, h = schemaHost(prefix3)
	if s != "http" || h != "127.0.0.1"  {
		t.Fatal(s, h, p)
	}

	s, h = schemaHost(prefix4)
	if s != "http" || h != "127.0.0.1" {
		t.Fatal(s, h, p)
	}

	s, h = schemaHost(prefix5)
	if s != "http" || h != "127.0.0.1" {
		t.Fatal(s, h, p)
	}

	s, h = schemaHost(prefix6)
	if s != "http" || h != "127.0.0.1" {
		t.Fatal(s, h, p)
	}
}