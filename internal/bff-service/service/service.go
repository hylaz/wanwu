package service

import (
	"fmt"
	knowledgebase_report_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-report-service"

	app_service "github.com/UnicomAI/wanwu/api/proto/app-service"
	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	knowledgebase_doc_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-doc-service"
	knowledgebase_keywords_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-keywords-service"
	knowledgebase_permission_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-permission-service"
	knowledgebase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	knowledgebase_splitter_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-splitter-service"
	knowledgebase_tag_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-tag-service"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	operate_service "github.com/UnicomAI/wanwu/api/proto/operate-service"
	perm_service "github.com/UnicomAI/wanwu/api/proto/perm-service"
	rag_service "github.com/UnicomAI/wanwu/api/proto/rag-service"
	safety_service "github.com/UnicomAI/wanwu/api/proto/safety-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	maxMsgSize            = 1024 * 1024 * 4 // 4M
	headlessServiceSchema = "dns:///"
)

var (
	iam                     iam_service.IAMServiceClient
	perm                    perm_service.PermServiceClient
	model                   model_service.ModelServiceClient
	mcp                     mcp_service.MCPServiceClient
	knowledgeBase           knowledgebase_service.KnowledgeBaseServiceClient
	knowledgeBaseDoc        knowledgebase_doc_service.KnowledgeBaseDocServiceClient
	knowledgeBaseTag        knowledgebase_tag_service.KnowledgeBaseTagServiceClient
	knowledgeBaseSplitter   knowledgebase_splitter_service.KnowledgeBaseSplitterServiceClient
	knowledgeBasePermission knowledgebase_permission_service.KnowledgeBasePermissionServiceClient
	knowledgeBaseKeywords   knowledgebase_keywords_service.KnowledgeBaseKeywordsServiceClient
	knowledgeBaseReport     knowledgebase_report_service.KnowledgeBaseReportServiceClient
	app                     app_service.AppServiceClient
	rag                     rag_service.RagServiceClient
	assistant               assistant_service.AssistantServiceClient
	safety                  safety_service.SafetyServiceClient
	operate                 operate_service.OperateServiceClient
)

// --- API ---

func Init() error {
	// grpc connections
	iamConn, err := newConn(config.Cfg().Iam.Host)
	if err != nil {
		return fmt.Errorf("init iam-service connection err: %v", err)
	}
	appConn, err := newConn(config.Cfg().App.Host)
	if err != nil {
		return fmt.Errorf("init app-service connection err: %v", err)
	}
	modelConn, err := newConn(config.Cfg().Model.Host)
	if err != nil {
		return fmt.Errorf("init model-service connection err: %v", err)
	}
	mcpConn, err := newConn(config.Cfg().MCP.Host)
	if err != nil {
		return fmt.Errorf("init mcp-service connection err: %v", err)
	}
	knowledgeBaseConn, err := newConn(config.Cfg().Knowledge.Host)
	if err != nil {
		return fmt.Errorf("init knowledgeBase-service connection err: %v", err)
	}
	ragConn, err := newConn(config.Cfg().Rag.Host)
	if err != nil {
		return fmt.Errorf("init rag-service connection err: %v", err)
	}
	assistantConn, err := newConn(config.Cfg().Assistant.Host)
	if err != nil {
		return fmt.Errorf("init assistant-service connection err: %v", err)
	}
	operateConn, err := newConn(config.Cfg().Operate.Host)
	if err != nil {
		return fmt.Errorf("init operate-service connection err: %v", err)
	}
	// grpc clients
	iam = iam_service.NewIAMServiceClient(iamConn)
	perm = perm_service.NewPermServiceClient(iamConn)
	model = model_service.NewModelServiceClient(modelConn)
	mcp = mcp_service.NewMCPServiceClient(mcpConn)
	app = app_service.NewAppServiceClient(appConn)
	knowledgeBase = knowledgebase_service.NewKnowledgeBaseServiceClient(knowledgeBaseConn)
	knowledgeBaseDoc = knowledgebase_doc_service.NewKnowledgeBaseDocServiceClient(knowledgeBaseConn)
	knowledgeBaseTag = knowledgebase_tag_service.NewKnowledgeBaseTagServiceClient(knowledgeBaseConn)
	knowledgeBaseKeywords = knowledgebase_keywords_service.NewKnowledgeBaseKeywordsServiceClient(knowledgeBaseConn)
	knowledgeBaseSplitter = knowledgebase_splitter_service.NewKnowledgeBaseSplitterServiceClient(knowledgeBaseConn)
	knowledgeBasePermission = knowledgebase_permission_service.NewKnowledgeBasePermissionServiceClient(knowledgeBaseConn)
	rag = rag_service.NewRagServiceClient(ragConn)
	assistant = assistant_service.NewAssistantServiceClient(assistantConn)
	safety = safety_service.NewSafetyServiceClient(appConn)
	operate = operate_service.NewOperateServiceClient(operateConn)
	return nil
}

// --- internal ---

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

func toIDName(idName *iam_service.IDName) response.IDName {
	return response.IDName{
		ID:   idName.Id,
		Name: idName.Name,
	}
}

func toIDNames(idNames []*iam_service.IDName) []response.IDName {
	var ret []response.IDName
	for _, idName := range idNames {
		ret = append(ret, toIDName(idName))
	}
	return ret
}
