package sse_util

import (
	"fmt"

	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/gin-gonic/gin"
)

const DONE_MSG = "data: [DONE]\n"

// SSEWriter 设计sse writer 目标可以规范化统一标准输出方法（所有sse 返回都能用），同时与业务尽可能解耦
type SSEWriter struct {
	client  *gin.Context
	label   string // 用于SSE日志中的标记
	doneMsg string // SSE结束时，发送给前端的结束消息，空不发送；一般为 "data: [DONE]\n"
}

func NewSSEWriter(c *gin.Context, label, doneMsg string) *SSEWriter {
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	return &SSEWriter{
		client:  c,
		label:   label,
		doneMsg: doneMsg,
	}
}

// WriteStream 流式写入，识别channel 循环写入给前端
func (sw *SSEWriter) WriteStream(sseCh <-chan string, streamContextParams interface{},
	lineBuilder func(*gin.Context, string, interface{}) (string, bool, error),
	doneProcessor func(*gin.Context, interface{}) error) error {
	for s := range sseCh {
		var lineText = s
		if lineBuilder != nil {
			line, skip, err := lineBuilder(sw.client, s, streamContextParams)
			if err != nil {
				log.Errorf("[SSE]%v line %v build err: %v", sw.label, err)
				return err
			}
			if skip {
				continue
			}
			lineText = line
		}
		if err := sw.WriteLine(lineText, false, streamContextParams, doneProcessor); err != nil {
			return err
		}
	}
	return sw.WriteLine("", true, streamContextParams, doneProcessor)
}

// WriteLine 写入一行给客户端
func (sw *SSEWriter) WriteLine(lineText string, done bool, streamProcessParams interface{},
	doneProcessor func(*gin.Context, interface{}) error) error {

	var err error
	defer func() {
		if err != nil {
			log.Errorf("[SSE]%v err: %v", sw.label, err)
		} else if done {
			log.Debugf("[SSE]%v done", sw.label)
		} else {
			return
		}
		// err 或 done 执行 doneProcessor
		if doneProcessor != nil {
			if err := doneProcessor(sw.client, streamProcessParams); err != nil {
				log.Errorf("[SSE]%v doneProcessor err: %v", sw.label, err)
			}
		}
	}()

	if done {
		lineText = fmt.Sprintf("%v%v", lineText, sw.doneMsg)
	}
	// 写入数据
	log.Infof("[SSE]%v write: %v", sw.label, lineText)
	_, err = sw.client.Writer.Write([]byte(lineText))
	if err != nil {
		err = fmt.Errorf("connection closed by web: %v", err)
		return err
	}
	sw.client.Writer.Flush()
	return nil
}
