package response

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
)

type CustomPromptIDResp struct {
	CustomPromptID string `json:"customPromptId"` // 自定义提示词ID
}

type CustomPrompt struct {
	CustomPromptIDResp                // 自定义提示词ID
	Avatar             request.Avatar `json:"avatar"`   // 图标
	Name               string         `json:"name"`     // 名称
	Desc               string         `json:"desc"`     // 描述
	Prompt             string         `json:"prompt"`   // 提示词
	UpdateAt           string         `json:"updateAt"` // 更新时间
}

type CustomPromptOpt struct {
	Code     *int                       `json:"code"`     // 状态码
	Message  string                     `json:"message"`  // 状态描述
	Response string                     `json:"response"` // 响应内容
	Finish   string                     `json:"finish"`   // 结束标志
	Usage    *mp_common.OpenAIRespUsage `json:"usage"`    // token使用统计
}
