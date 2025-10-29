package callback

import (
	"net/http"

	"github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/callback"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func Register(openAPI *gin.RouterGroup) {
	// callback
	mid.Sub("callback").Reg(openAPI, "/model/:modelId", http.MethodGet, callback.GetModelById, "根据modelId获取模型")
	mid.Sub("callback").Reg(openAPI, "/model/:modelId/chat/completions", http.MethodPost, callback.ModelChatCompletions, "Model Chat Completions")
	mid.Sub("callback").Reg(openAPI, "/model/:modelId/embeddings", http.MethodPost, callback.ModelEmbeddings, "Model Embeddings")
	mid.Sub("callback").Reg(openAPI, "/model/:modelId/rerank", http.MethodPost, callback.ModelRerank, "Model rerank")
	mid.Sub("callback").Reg(openAPI, "/model/:modelId/ocr", http.MethodPost, callback.ModelOcr, "Model ocr")
	mid.Sub("callback").Reg(openAPI, "/model/:modelId/gui", http.MethodPost, callback.ModelGui, "Model gui")
	mid.Sub("callback").Reg(openAPI, "/model/:modelId/pdf-parser", http.MethodPost, callback.ModelPdfParser, "Model pdf文档解析")
	// workflow
	mid.Sub("callback").Reg(openAPI, "/workflow/list", http.MethodGet, callback.GetWorkflowList, "根据userId和spaceId获取Workflow")
	mid.Sub("callback").Reg(openAPI, "/workflow/tool/square", http.MethodGet, callback.GetWorkflowSquareTool, "获取内置工具详情")
	mid.Sub("callback").Reg(openAPI, "/workflow/tool/custom", http.MethodGet, callback.GetWorkflowCustomTool, "获取自定义工具详情")
	// rag bff proxy
	mid.Sub("callback").Reg(openAPI, "/rag/search-knowledge-base", http.MethodPost, callback.SearchKnowledgeBase, "查询知识库列表（命中测试）")
	mid.Sub("callback").Reg(openAPI, "/rag/knowledge/stream/search", http.MethodPost, callback.KnowledgeStreamSearch, "根据知识库id 和当前用户id 获取有权限的知识库列表信息")
}
