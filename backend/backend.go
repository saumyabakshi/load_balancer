package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

)
type backend struct {
	url         *url.URL
	alive       bool
	mux         sync.RWMutex
	connections int
	reverseProxy *httputil.ReverseProxy
}

type Backend interface {
	IsAlive() bool
	SetAlive(bool)
	GetURL() *url.URL
	GetActiveConns() int
	Serve(http.ResponseWriter, *http.Request)
}

func (b *backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.alive
	defer b.mux.RUnlock()
	return alive 
}

func (b *backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.alive = alive
	b.mux.Unlock()
}	

func (b *backend) GetURL() *url.URL {
	return b.url
}

func (b *backend) GetActiveConns() int {
	// RLock is used to allow multiple go routines to read the connections but not write
	b.mux.RLock()
	conns := b.connections
	b.mux.RUnlock()
	return conns
}

func (b *backend) Serve(w http.ResponseWriter, r *http.Request) {
	defer func() {
		b.mux.Lock()
		b.connections--
		b.mux.Unlock()
	
	}()
	b.mux.Lock()
	b.connections++
	b.mux.Unlock()
	b.reverseProxy.ServeHTTP(w, r)

}

func NewBackend(url *url.URL, rp *httputil.ReverseProxy) Backend {
	b := &backend{
		url: url,
		alive: true,
		reverseProxy: rp,
	}
	return b
}

