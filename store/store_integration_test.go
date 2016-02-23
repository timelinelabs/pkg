// +build integration

package store

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	timeout = time.Duration(1 * 1000000000)
)

func init() {
}

func TestEtcdBackend(t *testing.T) {
	// This won't work - it only returns error on use not instantiation
	//s, er := NewEtcdStore([]string{"bad:/url"}, 5*time.Second)
	//assert.Error(t, er, "EtcdStore")
	s, er := NewEtcdStore([]string{"http://127.0.0.1:2379"}, 2*time.Second)
	assert.NoError(t, er, "EtcdStore")
	assert.Implements(t, (*Store)(nil), new(EtcdStore), "EtcdStore")
	testImplementation(t, s)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	s, er := NewEtcdStore([]string{"http://127.0.0.1:2379"}, 1*time.Second)
	assert.NoError(t, er, "EtcdStore")

	for i := 0; i < 100; i++ {
		key := RandStringBytes(8)

		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			value := []byte("value")
			dur := 1 * time.Second
			err := s.Set(key, value, dur)
			assert.NoError(t, err)
		}(key)

		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			dur := 1 * time.Second
			time.Sleep(dur)
			reader := s.Get(key)
			actual := streamToString(reader)
			assert.True(t, actual == "" || actual == "value")
		}(key)
	}
}

func TestSetTTL(t *testing.T) {

	etcd, _ := NewEtcdStore([]string{"http://127.0.0.1:2379"}, timeout)

	stores := map[string]Store{
		"etcd": etcd,
		"mem":  NewMemStore(),
	}

	for _, b := range stores {
		deleteKeys(b, "/howdy")
		value := []byte("value")
		dur, _ := time.ParseDuration("1s")
		err := b.Set("/howdy", value, dur)
		assert.NoError(t, err)
		reader := b.Get("/howdy")
		time.Sleep(dur * 2)
		reader = b.Get("howdy")
		actual := streamToString(reader)
		expected := ""
		assert.Equal(t, expected, actual)
	}
}
