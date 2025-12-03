package message_builder

import (
	"context"
	"encoding/json"
)

type KnowledgeStartBuilder struct {
}

func (KnowledgeStartBuilder) MessageType() RagMessageType {
	return KnowledgeStart
}

func (KnowledgeStartBuilder) Build(ctx context.Context, ragContext *RagContext) *RagEvent {
	ragMessage := &RagMessage{
		Code:    0,
		Message: "success",
		MsgId:   ragContext.MessageId,
		MsgType: KnowledgeStart,
		History: make([]*RagHistory, 0),
		Data: &RagData{
			Output:     "",
			SearchList: make([]interface{}, 0),
		},
	}
	data, err := json.Marshal(ragMessage)
	if err != nil {
		return &RagEvent{
			Error: err,
		}
	}
	return &RagEvent{
		Message: []string{LineData(data)},
	}
}
