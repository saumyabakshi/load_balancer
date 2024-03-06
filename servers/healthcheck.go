package servers

import (
	"context"
	"time"

	"github.com/saumyabakshi/load_balancer/utils"
)

func StartHealthCheck(ctx context.Context, s Servers) {
	t := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-t.C:
			utils.Logger.Info("Running health check")
			go HealthCheck(ctx, s)
		case <-ctx.Done():
			utils.Logger.Info("Stopping health check")
			return
		}
	}
}