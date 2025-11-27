package knowledge_qa

import (
	knowledgebase_qa_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-qa-service"
	grpc_provider "github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/grpc-provider"
	"google.golang.org/grpc"
)

type Service struct {
	knowledgebase_qa_service.UnimplementedKnowledgeBaseQAServiceServer
}

var docService = Service{}

func init() {
	grpc_provider.AddGrpcContainer(&docService)
}

func (s *Service) GrpcType() string {
	return "grpc_knowledge_qa_service"
}

func (s *Service) Register(serv *grpc.Server) error {
	knowledgebase_qa_service.RegisterKnowledgeBaseQAServiceServer(serv, s)
	return nil
}
