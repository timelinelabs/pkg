// +build integration

package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
}

func TestEtcdBackend(t *testing.T) {
	s, er := NewEtcdStore([]string{"bad:/url"}, 5*time.Second)
	assert.Error(t, err, "EtcdStore")
	s, er = NewEtcdStore([]string{"http://127.0.0.1:2379"}, 2*time.Second)
	assert.NoError(t, er, "EtcdStore")
	assert.Implements(t, (*Store)(nil), new(EtcdStore), "EtcdStore")
	testImplementation(t, s)
}

func TestSetTTL(t *testing.T) {

	stores := map[string]Store{
		"etcd": NewEtcdStore([]string{"http://127.0.0.1:2379"}),
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
