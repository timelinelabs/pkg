// +build integration

package store

import (
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
