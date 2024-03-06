package servers

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"testing"

	"github.com/saumyabakshi/load_balancer/backend"
	"github.com/stretchr/testify/assert"
)

func TestRoundRobinTestCreation(t *testing.T) {
	s,_ := NewServers("round-robin")
	url,_ := url.Parse("http://localhost:3000")
	b := backend.NewBackend(url, httputil.NewSingleHostReverseProxy(url))
	s.AddBackend(b)
	assert.Equal(t, 1, len(s.GetBackends()))
}

func TestLeastConnectionsTestCreation(t *testing.T) {
	s,_ := NewServers("least-connections")
	url,_ := url.Parse("http://localhost:3010")
	b := backend.NewBackend(url, httputil.NewSingleHostReverseProxy(url))
	s.AddBackend(b)
	assert.Equal(t, 1, len(s.GetBackends()))
}

func TestRoundRobinTestRotate(t *testing.T) {
	s,_ := NewServers("round-robin")
	url,_ := url.Parse("http://localhost:3000")
	b := backend.NewBackend(url, httputil.NewSingleHostReverseProxy(url))
	s.AddBackend(b)
	
	url,_ = url.Parse("http://localhost:3001")
	b2 := backend.NewBackend(url, httputil.NewSingleHostReverseProxy(url))
	s.AddBackend(b2)

	url,_ = url.Parse("http://localhost:3002")
	b3 := backend.NewBackend(url, httputil.NewSingleHostReverseProxy(url))
	s.AddBackend(b3)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			s.GetNextValidPeer()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			s.GetNextValidPeer()
		}
	}()

	wg.Wait()
	assert.Equal(t, b.GetURL().String(), s.GetNextValidPeer().GetURL().String())
	
}




