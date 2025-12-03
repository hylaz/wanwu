package response

import "github.com/UnicomAI/wanwu/internal/bff-service/model/request"

type RagInfo struct {
	RagID string `json:"ragId" validate:"required"`
	request.AppBriefConfig
	ModelConfig           request.AppModelConfig           `json:"modelConfig" validate:"required"`           // 模型
	RerankConfig          request.AppModelConfig           `json:"rerankConfig" validate:"required"`          // Rerank模型
	QARerankConfig        request.AppModelConfig           `json:"qaRerankConfig" validate:"required"`        // 问答库Rerank模型
	KnowledgeBaseConfig   request.AppKnowledgebaseConfig   `json:"knowledgeBaseConfig" validate:"required"`   // 知识库
	QAKnowledgeBaseConfig request.AppQAKnowledgebaseConfig `json:"qaKnowledgeBaseConfig" validate:"required"` // 问答库
	SafetyConfig          request.AppSafetyConfig          `json:"safetyConfig"`                              // 敏感词表配置
}
