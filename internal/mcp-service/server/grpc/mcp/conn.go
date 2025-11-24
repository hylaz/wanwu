package mcp

import (
	"fmt"

	app_service "github.com/UnicomAI/wanwu/api/proto/app-service"
	"github.com/UnicomAI/wanwu/internal/mcp-service/config"
	"github.com/UnicomAI/wanwu/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	maxMsgSize            = 1024 * 1024 * 4 // 4M
	headlessServiceSchema = "dns:///"
)

var (
	app app_service.AppServiceClient
)

func StartService() error {
	// grpc connections
	AppConn, err := newConn(config.Cfg().App.Host)
	if err != nil {
		return fmt.Errorf("init app-service connection err: %v", err)
	}
	app = app_service.NewAppServiceClient(AppConn)
	log.Infof("App init success")
	log.Infof("App: %s", config.Cfg().App.Host)
	return nil
}

func newConn(host string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(headlessServiceSchema+host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize)),
	)
	if err != nil {
		return nil, err
	}
	return conn, err
}
