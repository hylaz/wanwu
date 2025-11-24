package knowledge_permission

import (
	knowledgebase_permission_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-permission-service"
	grpc_provider "github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/grpc-provider"
	"google.golang.org/grpc"
)

type Service struct {
	knowledgebase_permission_service.UnimplementedKnowledgeBasePermissionServiceServer
}

var docService = Service{}

func init() {
	grpc_provider.AddGrpcContainer(&docService)
}

func (s *Service) GrpcType() string {
	return "grpc_knowledge_permission_service"
}

func (s *Service) Register(serv *grpc.Server) error {
	knowledgebase_permission_service.RegisterKnowledgeBasePermissionServiceServer(serv, s)
	return nil
}
