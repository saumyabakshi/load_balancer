package servers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/saumyabakshi/load_balancer/backend"
	"github.com/saumyabakshi/load_balancer/utils"
	"go.uber.org/zap"
)

type Servers interface {
	GetBackends() []backend.Backend
	GetNextValidPeer() backend.Backend
	AddBackend(backend.Backend)
}

type roundRobin struct {
	backends []backend.Backend
	mux	  sync.RWMutex
	current int
}

type leastConnections struct {
	backends []backend.Backend
	mux	  sync.RWMutex
}

func (r *roundRobin) Rotate() backend.Backend {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.current = (r.current + 1) % len(r.backends)
	return r.backends[r.current]
}

func (r *roundRobin) GetBackends() []backend.Backend {
	return r.backends
}

func (r *roundRobin) GetNextValidPeer() backend.Backend {
	for i := 0; i < len(r.backends); i++ {
		next := r.Rotate()
		if next.IsAlive() {
			utils.Logger.Info("Bakcend selected", zap.String("url", next.GetURL().String()))
			return next
		}
	}
	return nil
}

func (r *roundRobin) AddBackend(b backend.Backend) {
	r.backends = append(r.backends, b)
}

func (l *leastConnections) GetBackends() []backend.Backend {
	return l.backends
}

func (l *leastConnections) GetNextValidPeer() backend.Backend {
	var leastConnected backend.Backend
	for _, b := range l.backends {
		if b.IsAlive() {
			leastConnected = b
			break
		}
	}
	for _, b := range l.backends {
		if !b.IsAlive() {
			continue
		}
		if leastConnected.GetActiveConns() > b.GetActiveConns() {
			leastConnected = b
		}
	}
	return leastConnected
}

func (l *leastConnections) AddBackend(b backend.Backend) {
	l.backends = append(l.backends, b)
}

func HealthCheck(ctx context.Context, s Servers) {
	aliveChannel := make(chan bool,1)

	for _, b := range s.GetBackends() {
		b := b
		rCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		status := "up"
		go backend.IsBackendAlive(rCtx, aliveChannel, b.GetURL())

		select {
		case <-ctx.Done():
			utils.Logger.Info("Health check stopped")
			return
		case alive := <-aliveChannel:
			b.SetAlive(alive)
			if !alive {
				status = "down"
			}
		}
		utils.Logger.Debug("Health check on backend", zap.String("url", b.GetURL().String()), zap.String("status", status))
}
}

func NewServers( strategy string) (Servers, error) {
	switch strategy {
	case "round-robin":
		return &roundRobin{
			backends: make([]backend.Backend, 0),
			current: 0,
		}, nil
	case "least-connections":
		return &leastConnections{
			backends: make([]backend.Backend, 0),
		},nil
	default:
		return nil, fmt.Errorf("unknown strategy %s", strategy)
	}

}