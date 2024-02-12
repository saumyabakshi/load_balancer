package backend


import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)
type backend struct {
	url *url.URL
	alive bool
	mux sync.RWMutex
	connections int32
}

type Backend interface {
	IsAlive() bool
	SetAlive(bool)
	GetURL() *url.URL
	GetActiveConns() int32
	Serve(http.ResponseWriter, *http.Request)
}

func (b *backend) IsAlive() bool {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.alive 
}

func (b *backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.alive = alive
	defer b.mux.Unlock()
}

func (b *backend) GetURL() *url.URL {
	return b.url
}

func (b *backend) GetActiveConns() int32 {
	b.mux.RLock()
	conns := b.connections
	defer b.mux.RUnlock()
	return conns
}

func (b *backend) Serve(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&b.connections, 1)
	defer atomic.AddInt32(&b.connections, -1)
	b.reverseProxy.ServeHTTP(w, r)

}

func NewBackend(url *url.URL, rp *httputil.reverseProxy) Backend {
	b := &backend{
		url: url,
		alive: true,
		reverseProxy: rp,
	}
	return b
}

