/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nodes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-flow/prompt"
	"strings"
	"time"

	"github.com/UnicomAI/wanwu/internal/agent-service/model"
	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/config"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/http"
	http_client "github.com/UnicomAI/wanwu/pkg/http-client"
	"github.com/UnicomAI/wanwu/pkg/log"
)

const (
	successCode = 0
)

type KnowledgeRetriever struct {
}

func (k *KnowledgeRetriever) Retrieve(ctx context.Context, reqContext *request.AgentChatContext) (string, error) {
	req := reqContext.AgentChatReq
	if req.KnowledgeParams == nil {
		return "", nil
	}
	req.KnowledgeParams.Question = req.Input
	req.KnowledgeParams.CustomModelInfo = &request.CustomModelInfo{
		LlmModelID: req.ModelParams.ModelId,
	}
	hit, _ := ragKnowledgeHit(ctx, req.KnowledgeParams)
	if hit == nil {
		return "", nil
	}
	reqContext.KnowledgeHitData = hit.Data
	packedRes := strings.Builder{}
	for idx, doc := range hit.Data.SearchList {
		if doc == nil {
			continue
		}
		packedRes.WriteString(fmt.Sprintf("---\nrecall slice %d: %s\n", idx+1, doc.Snippet))
	}
	knowledgeData := packedRes.String()
	if len(knowledgeData) > 0 {
		knowledgeData = fmt.Sprintf(prompt.REACT_SYSTEM_PROMPT_KNOWLEDGE, knowledgeData)
		return knowledgeData, nil
	}
	//如果没有知识库时，尽量减少输入token大小
	return "", nil
}

// RagKnowledgeHit rag命中测试
func ragKnowledgeHit(ctx context.Context, knowledgeHitParams *request.KnowledgeParams) (*model.RagKnowledgeHitResp, error) {
	ragServer := config.GetConfig().RagServer
	url := ragServer.ProxyPoint + ragServer.KnowledgeHitUri
	paramsByte, err := json.Marshal(knowledgeHitParams)
	if err != nil {
		return nil, err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_knowledge_hit",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return nil, err
	}
	var resp model.RagKnowledgeHitResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	if resp.Code != successCode {
		return nil, errors.New(resp.Message)
	}
	return &resp, nil
}
