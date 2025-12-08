package agent_log

import (
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/config"
	"github.com/UnicomAI/wanwu/pkg/log"
)

// 打印等级 Panic > Error > Warn > Info > Debug

var agentLog = AgentLog{}

type AgentLog struct {
}

func init() {
	pkg.AddContainer(agentLog)
}

func (c AgentLog) LoadType() string {
	return "agent-log-config"
}

func (c AgentLog) Load() error {
	configInfo := config.GetConfig()
	logConfig := configInfo.Log
	return log.InitLog(logConfig.Std, logConfig.Level, logConfig.Logs...)
}

func (c AgentLog) StopPriority() int {
	return pkg.DefaultPriority
}

func (c AgentLog) Stop() error {
	return nil
}
