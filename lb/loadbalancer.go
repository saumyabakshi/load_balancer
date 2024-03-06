package lb

import (
	"net/http"

	"github.com/saumyabakshi/load_balancer/servers"
	"github.com/saumyabakshi/load_balancer/utils"
)



func AllowRetry(r *http.Request) bool {
	if _, ok := r.Context().Value(3).(bool); ok {
		return false
	}
	return true
}

type LoadBalancer interface {
	Serve(http.ResponseWriter, *http.Request)
}

type loadBalancer struct {
	s servers.Servers
}

func (lb *loadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	peer := lb.s.GetNextValidPeer()
	if peer != nil {
		utils.Logger.Info("Proxying request")
		peer.Serve(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func NewLoadBalancer(sp servers.Servers) LoadBalancer {
	return &loadBalancer{
		s: sp,
	}
}
