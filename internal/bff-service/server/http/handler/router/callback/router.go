package callback

import (
	"net/http"

	"github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/callback"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func Register(callbackAPI *gin.RouterGroup) {
	// callback
	mid.Sub("callback").Reg(callbackAPI, "/file/url/base64", http.MethodPost, callback.FileUrlConvertBase64, "文件URL转换为base64")
	mid.Sub("callback").Reg(callbackAPI, "/tool/builtin/asr", http.MethodPost, callback.AudioBase64ConvertText, "语音文件（base64格式）转文本内置工具服务")

	// model
	mid.Sub("callback").Reg(callbackAPI, "/model/:modelId", http.MethodGet, callback.GetModelById, "根据modelId获取模型")
	mid.Sub("callback").Reg(callbackAPI, "/model/:modelId/chat/completions", http.MethodPost, callback.ModelChatCompletions, "Model Chat Completions")
	mid.Sub("callback").Reg(callbackAPI, "/model/:modelId/embeddings", http.MethodPost, callback.ModelEmbeddings, "Model Embeddings")
	mid.Sub("callback").Reg(callbackAPI, "/model/:modelId/rerank", http.MethodPost, callback.ModelRerank, "Model rerank")
	mid.Sub("callback").Reg(callbackAPI, "/model/:modelId/ocr", http.MethodPost, callback.ModelOcr, "Model ocr")
	mid.Sub("callback").Reg(callbackAPI, "/model/:modelId/gui", http.MethodPost, callback.ModelGui, "Model gui")
	mid.Sub("callback").Reg(callbackAPI, "/model/:modelId/pdf-parser", http.MethodPost, callback.ModelPdfParser, "Model pdf文档解析")
	//mid.Sub("callback").Reg(callbackAPI, "/model/:modelId/asr", http.MethodPost, callback.ModelAsr, "Model asr")
	//mid.Sub("callback").Reg(callbackAPI, "/model/:modelId/text2image", http.MethodPost, callback.ModelText2Image, "Model asr")
	// workflow
	mid.Sub("callback").Reg(callbackAPI, "/workflow/list", http.MethodGet, callback.GetWorkflowList, "根据userId和spaceId获取Workflow")
	mid.Sub("callback").Reg(callbackAPI, "/workflow/tool/square", http.MethodGet, callback.GetWorkflowSquareTool, "获取内置工具详情")
	mid.Sub("callback").Reg(callbackAPI, "/workflow/tool/custom", http.MethodGet, callback.GetWorkflowCustomTool, "获取自定义工具详情")
	mid.Sub("callback").Reg(callbackAPI, "/workflow/upload/file", http.MethodPost, callback.WorkflowUploadFile, "通过二进制上传文件")
	mid.Sub("callback").Reg(callbackAPI, "/workflow/upload/file/base64", http.MethodPost, callback.WorkflowUploadFileByBase64, "通过base64上传文件")
	// chatflow
	mid.Sub("callback").Reg(callbackAPI, "/chatflow/list", http.MethodGet, callback.GetChatflowList, "根据userId和spaceId获取Chatflow")
	// rag bff proxy
	mid.Sub("callback").Reg(callbackAPI, "/rag/search-knowledge-base", http.MethodPost, callback.SearchKnowledgeBase, "查询知识库列表（命中测试）")
	mid.Sub("callback").Reg(callbackAPI, "/rag/knowledge/stream/search", http.MethodPost, callback.KnowledgeStreamSearch, "根据知识库id 和当前用户id 获取有权限的知识库列表信息")
	// rag bff proxy
	mid.Sub("callback").Reg(callbackAPI, "/rag/search-QA-base", http.MethodPost, callback.SearchQABase, "查询问答库列表（命中测试）")
}
