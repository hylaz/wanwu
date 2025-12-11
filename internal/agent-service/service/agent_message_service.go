package service

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/log"
	sse_util "github.com/UnicomAI/wanwu/pkg/sse-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

// AgentMessage 智能体消息处理
func AgentMessage(ctx *gin.Context, iter *adk.AsyncIterator[*adk.AgentEvent], req *request.AgentChatContext) (err error) {
	rawCh := make(chan string, 128)
	go func() {
		defer util.PrintPanicStack()
		defer close(rawCh)
		respContext := response.NewAgentChatRespContext()

		for {

			event, ok := iter.Next()
			if !ok {
				break
			}

			if event.Err != nil {
				log.Errorf("agent event result error %v", event.Err)
				err = event.Err
				return
			}

			output := event.Output.MessageOutput
			err := streamReceive(rawCh, output, respContext, req)
			if err != nil {
				log.Errorf("agent stream receive error %v", event.Err)
				continue
			}
		}
	}()

	// 2.流式返回结果
	return sse_util.NewSSEWriter(ctx, fmt.Sprintf("[Agent] %v ", req.AgentChatReq.Input), "").
		WriteStream(rawCh, nil, buildAgentChatRespLineProcessor(), nil)
}

func streamReceive(sseCh chan string, output *adk.MessageVariant, respContext *response.AgentChatRespContext, req *request.AgentChatContext) error {
	var msg *schema.Message
	if output.IsStreaming {
		err := streamMessage(sseCh, output.MessageStream, respContext, req)
		if err != nil {
			return err
		}

	} else {

		msg = output.Message
		respList, err := response.NewAgentChatRespWithTool(msg, respContext, req)
		if err != nil {
			log.Errorf("MessageOutput error %v", err)
			return err
		}
		for _, resp := range respList {
			sseCh <- resp
		}
	}
	return nil
}

// streamMessage 处理流信息
func streamMessage(sseCh chan string, s *schema.StreamReader[*schema.Message], respContext *response.AgentChatRespContext, req *request.AgentChatContext) error {
	defer s.Close()

	for {
		msg, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		messageJSON, _ := json.Marshal(msg)
		log.Infof("stream message %v", string(messageJSON))
		respList, err := response.NewAgentChatRespWithTool(msg, respContext, req)
		if err != nil {
			log.Errorf("MessageOutput error %v", err)
			return err
		}
		for _, resp := range respList {
			sseCh <- resp
		}
	}
	return nil
}

// buildAgentChatRespLineProcessor 构造rag对话结果行处理器
func buildAgentChatRespLineProcessor() func(*gin.Context, string, interface{}) (string, bool, error) {
	return func(c *gin.Context, lineText string, params interface{}) (string, bool, error) {
		if strings.HasPrefix(lineText, "error:") {
			errorText := fmt.Sprintf("data: {\"code\": -1, \"message\": \"%s\"}\n\n", strings.TrimPrefix(lineText, "error:"))
			return errorText, false, nil
		}
		if strings.HasPrefix(lineText, "data:") {
			return lineText + "\n", false, nil
		}
		return lineText + "\n", false, nil
	}
}
