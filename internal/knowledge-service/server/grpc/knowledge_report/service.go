package knowledge_report

import (
	knowledgebase_report_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-report-service"
	grpc_provider "github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/grpc-provider"
	"google.golang.org/grpc"
)

type Service struct {
	knowledgebase_report_service.UnimplementedKnowledgeBaseReportServiceServer
}

var reportService = Service{}

func init() {
	grpc_provider.AddGrpcContainer(&reportService)
}

func (s *Service) GrpcType() string {
	return "grpc_knowledge_report_service"
}

func (s *Service) Register(serv *grpc.Server) error {
	knowledgebase_report_service.RegisterKnowledgeBaseReportServiceServer(serv, s)
	return nil
}
