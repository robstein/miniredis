package miniredis

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

// Test starting/stopping a server
func TestServer(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()

	c, err := redis.Dial("tcp", s.Addr())
	ok(t, err)
	_, err = c.Do("PING")
	ok(t, err)
}

func TestMultipleServers(t *testing.T) {
	s1, err := Run()
	ok(t, err)
	s2, err := Run()
	ok(t, err)
	if s1.Addr() == s2.Addr() {
		t.Fatal("Non-unique addresses", s1.Addr(), s2.Addr())
	}

	s2.Close()
	s1.Close()
	// Closing multiple times is fine
	go s1.Close()
	go s1.Close()
	s1.Close()
}

// Test simple GET/SET keys
func TestKeys(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c, err := redis.Dial("tcp", s.Addr())
	ok(t, err)

	// SET command
	_, err = c.Do("SET", "foo", "bar")
	ok(t, err)
	// GET command
	v, err := redis.String(c.Do("GET", "foo"))
	ok(t, err)
	equals(t, "bar", v)

	// Query server directly.
	equals(t, "bar", s.Get("foo"))
}
