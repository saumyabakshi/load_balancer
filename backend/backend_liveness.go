package backend

import (
	"context"
	"net"
	"net/url"

	"github.com/saumyabakshi/load_balancer/utils"
	"go.uber.org/zap"
)

func IsBackendAlive(ctx context.Context, aliveChannel chan bool,url *url.URL) {
	var d net.Dialer
	utils.Logger.Info("Checking backend", zap.String("url", url.Host))
	conn, err := d.DialContext(ctx, "tcp", url.Host)
	if err != nil {
		utils.Logger.Debug("Error dialing the backend", zap.String("url", url.String()), zap.Error(err))
		aliveChannel <- false
		return
	}
	_ = conn.Close()
	aliveChannel <- true
}