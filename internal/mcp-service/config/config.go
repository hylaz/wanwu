package config

import (
	"context"
	"fmt"

	"github.com/UnicomAI/wanwu/pkg/db"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
)

var (
	_c *Config
)

type Config struct {
	Server  ServerConfig  `json:"server" mapstructure:"server"`
	Log     LogConfig     `json:"log" mapstructure:"log"`
	DB      db.Config     `json:"db" mapstructure:"db"`
	Mcps    []*McpConfig  `json:"mcps" mapstructure:"mcps"`
	Tools   []*ToolConfig `json:"tools" mapstructure:"tools"`
	McpCfg  McpCfg        `json:"mcp" mapstructure:"mcp"`
	ToolCfg ToolCfg       `json:"tool" mapstructure:"tool"`
	App     App           `json:"app" mapstructure:"app"`
}

type App struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
}

type McpCfg struct {
	ConfigPath string `json:"configPath" mapstructure:"configPath"`
}

type ToolCfg struct {
	ConfigPath string `json:"configPath" mapstructure:"configPath"`
}

type ServerConfig struct {
	GrpcEndpoint   string `json:"grpc_endpoint" mapstructure:"grpc_endpoint"`
	MaxRecvMsgSize int    `json:"max_recv_msg_size" mapstructure:"max_recv_msg_size"`
	ApiBaseUrl     string `json:"api_base_url" mapstructure:"api_base_url"`
}

type LogConfig struct {
	Std   bool         `json:"std" mapstructure:"std"`
	Level string       `json:"level" mapstructure:"level"`
	Logs  []log.Config `json:"logs" mapstructure:"logs"`
}

func LoadConfig(ctx context.Context, in string) error {
	_c = &Config{}
	if err := util.LoadConfig(in, _c); err != nil {
		return err
	}
	mcpIn := _c.McpCfg.ConfigPath
	toolIn := _c.ToolCfg.ConfigPath
	if err := util.LoadConfig(mcpIn, _c); err != nil {
		return fmt.Errorf("load mcp config err: %v", err)
	}
	if err := util.LoadConfig(toolIn, _c); err != nil {
		return fmt.Errorf("load tool config err: %v", err)
	}
	for _, mcp := range _c.Mcps {
		if err := mcp.load(); err != nil {
			return err
		}
	}
	for _, tool := range _c.Tools {
		if err := tool.load(ctx); err != nil {
			return err
		}
	}
	return nil
}

func Cfg() *Config {
	if _c == nil {
		log.Panicf("cfg nil")
	}
	return _c
}

func (c *Config) MCP(mcpSquareID string) (McpConfig, bool) {
	for _, mcp := range c.Mcps {
		if mcp.McpSquareId == mcpSquareID {
			return *mcp, true
		}
	}
	return McpConfig{}, false
}

func (c *Config) Tool(toolSquareID string) (ToolConfig, bool) {
	for _, tool := range c.Tools {
		if tool.ToolSquareId == toolSquareID {
			return *tool, true
		}
	}
	return ToolConfig{}, false
}
