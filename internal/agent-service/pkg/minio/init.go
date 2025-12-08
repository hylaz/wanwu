package minio

import (
	"context"

	"github.com/UnicomAI/wanwu/internal/agent-service/pkg"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/config"
	"github.com/UnicomAI/wanwu/pkg/minio"
)

var minioClient = ClientMinio{}

type ClientMinio struct {
}

func init() {
	pkg.AddContainer(minioClient)
}

func (c ClientMinio) LoadType() string {
	return "minioClient"
}

func (c ClientMinio) Load() error {
	minioConfig := config.GetConfig().Minio
	//初始化知识库内部使用minio-bucket
	err := minio.InitAgent(context.Background(), minio.Config{
		Endpoint: minioConfig.EndPoint,
		User:     minioConfig.User,
		Password: minioConfig.Password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c ClientMinio) StopPriority() int {
	return pkg.DefaultPriority
}

func (c ClientMinio) Stop() error {
	return nil
}
