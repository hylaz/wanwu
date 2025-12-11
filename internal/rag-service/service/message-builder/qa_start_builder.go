package message_builder

import (
	"context"
	"encoding/json"
)

type QAStartBuilder struct {
}

func (QAStartBuilder) MessageType() RagMessageType {
	return QAStart
}

func (QAStartBuilder) Build(ctx context.Context, ragContext *RagContext) *RagEvent {
	if len(ragContext.QAIds) == 0 {
		return &RagEvent{Skip: true}
	}

	ragMessage := &RagMessage{
		Code:    0,
		Message: "success",
		MsgId:   ragContext.MessageId,
		MsgType: QAStart,
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
