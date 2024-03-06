package backend

import (
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBackend(t *testing.T) {
	url, _ := url.Parse("http://localhost:3000/")
	rp := httputil.NewSingleHostReverseProxy(url)
	b := NewBackend(url, rp)
	assert.Equal(t, b.GetURL().String(), "http://localhost:3000/")
	assert.Equal(t, b.IsAlive(), true)
}

func TestBackend_IsAlive(t *testing.T) {
	url, _ := url.Parse("http://localhost:3000/")
	url2, _ := url.Parse("http://localhost:3001/")
	rp := httputil.NewSingleHostReverseProxy(url)
	b := NewBackend(url, rp)
	b2 := NewBackend(url2, httputil.NewSingleHostReverseProxy(url2))
	b.SetAlive(true)
	b2.SetAlive(true)
	assert.Equal(t, b2.IsAlive(), b.IsAlive())

	
}
