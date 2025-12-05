package http

import (
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg"
	http_client "github.com/UnicomAI/wanwu/pkg/http-client"
)

var localHttpClient = LocalHttpClient{}

type LocalHttpClient struct {
	Client *http_client.HttpClient
}

func init() {
	pkg.AddContainer(localHttpClient)
}

func (c LocalHttpClient) LoadType() string {
	return "http-client"
}

func (c LocalHttpClient) Load() error {
	localHttpClient.Client = http_client.CreateDefault()
	return nil
}

func (c LocalHttpClient) Stop() error {
	return nil
}

func (c LocalHttpClient) StopPriority() int {
	return pkg.DefaultPriority
}

func GetClient() *http_client.HttpClient {
	return localHttpClient.Client
}
