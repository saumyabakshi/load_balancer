package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saumyabakshi/load_balancer/backend"
	"github.com/saumyabakshi/load_balancer/lb"
	"github.com/saumyabakshi/load_balancer/servers"
	"github.com/saumyabakshi/load_balancer/utils"
	"go.uber.org/zap"
)

func main() {
	logger := utils.InitLogger()
	defer logger.Sync()

	var (
		strategy = flag.String("strategy", "round-robin", "Load balancing strategy")
		port = flag.Int("port", 8000, "Port to serve on")
	)



	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverPool, err := servers.NewServers(*strategy)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	loadBalancer := lb.NewLoadBalancer(serverPool)

	backends := []string{
		"http://localhost:3333",
		"http://localhost:3332",
	}

	for _, u := range backends {
		endpoint, err := url.Parse(u)
		if err != nil {
			utils.Logger.Fatal(err.Error(), zap.String("URL", u))
		}
		rp := httputil.NewSingleHostReverseProxy(endpoint)
		backendServer := backend.NewBackend(endpoint, rp)
		rp.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			utils.Logger.Error("error handling the request",
				zap.String("host", endpoint.Host),
				zap.Error(e),
			)
			backendServer.SetAlive(false)

			if !lb.AllowRetry(request) {
				utils.Logger.Info(
					"Max retry attempts reached, terminating",
					zap.String("address", request.RemoteAddr),
					zap.String("path", request.URL.Path),
				)
				http.Error(writer, "Service not available", http.StatusServiceUnavailable)
				return
			}

			utils.Logger.Info(
				"Attempting retry",
				zap.String("address", request.RemoteAddr),
				zap.String("URL", request.URL.Path),
				zap.Bool("retry", true),
			)
			loadBalancer.Serve(
				writer,
				request.WithContext(
					context.WithValue(request.Context(), 3, true),
				),
			)
		}

		serverPool.AddBackend(backendServer)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: http.HandlerFunc(loadBalancer.Serve),
	}

	go servers.StartHealthCheck(ctx, serverPool)

	go func() {
		<-ctx.Done()
		shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		utils.Logger.Fatal("ListenAndServe() error", zap.Error(err))
	}
}