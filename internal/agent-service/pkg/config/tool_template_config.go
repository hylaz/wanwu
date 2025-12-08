package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg"
)

const (
	toolTemplateConfigPath = "configs/microservice/agent-service/tool-template/doc_parser.json"
)

type ToolTemplateConfig struct {
	ConfigPluginToolList []*request.PluginToolInfo
}

var toolTemplateConfig = ToolTemplateConfig{}

func init() {
	pkg.AddContainer(toolTemplateConfig)
}

func GetToolTemplateConfig() *ToolTemplateConfig {
	return &toolTemplateConfig
}

func (c ToolTemplateConfig) LoadType() string {
	return "tool-template-config"
}

func (c ToolTemplateConfig) Load() error {
	b, err := os.ReadFile(toolTemplateConfigPath)
	if err != nil {
		return fmt.Errorf("load tool template path %v err: %v", toolTemplateConfigPath, err)
	}
	toolConfig := fmt.Sprintf(string(b), GetConfig().ToolServer.Endpoint)
	var pluginTool = request.PluginToolInfo{}
	err = json.Unmarshal([]byte(toolConfig), &pluginTool)
	if err != nil {
		return fmt.Errorf("unmarshal tool template path %v err: %v", toolTemplateConfigPath, err)
	}
	//先简单写只加载一个
	toolTemplateConfig.ConfigPluginToolList = []*request.PluginToolInfo{&pluginTool}
	return nil
}

func (c ToolTemplateConfig) StopPriority() int {
	return pkg.DefaultPriority
}

func (c ToolTemplateConfig) Stop() error {
	return nil
}
