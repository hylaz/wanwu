package config

import (
	smtp_util "github.com/UnicomAI/wanwu/internal/iam-service/pkg/util/smtp-util"
	"github.com/UnicomAI/wanwu/pkg/db"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/redis"
	"github.com/UnicomAI/wanwu/pkg/util"
)

var (
	_c *Config
)

type Config struct {
	Server   ServerConfig     `json:"server" mapstructure:"server"`
	Log      LogConfig        `json:"log" mapstructure:"log"`
	DB       db.Config        `json:"db" mapstructure:"db"`
	Redis    redis.Config     `json:"redis" mapstructure:"redis"`
	SMTP     smtp_util.Config `json:"smtp" mapstructure:"smtp"`
	Register RegisterConfig   `json:"register" mapstructure:"register"`
	Password PasswordConfig   `json:"password" mapstructure:"password"`
}

type ServerConfig struct {
	GrpcEndpoint   string `json:"grpc_endpoint" mapstructure:"grpc_endpoint"`
	MaxRecvMsgSize int    `json:"max_recv_msg_size" mapstructure:"max_recv_msg_size"`
}

type LogConfig struct {
	Std   bool         `json:"std" mapstructure:"std"`
	Level string       `json:"level" mapstructure:"level"`
	Logs  []log.Config `json:"logs" mapstructure:"logs"`
}

type DBConfig struct {
	Name string `json:"name" mapstructure:"name"`
}

type RegisterConfig struct {
	Email RegisterByEmail `json:"email" mapstructure:"email"`
}

type RegisterByEmail struct {
	CodeLength     int                     `json:"code_length" mapstructure:"code_length"`
	PasswordLength int                     `json:"password_length" mapstructure:"password_length"`
	Template       RegisterByEmailTemplate `json:"template" mapstructure:"template"`
}

type PasswordConfig struct {
	Email RegisterByEmail `json:"email" mapstructure:"email"`
}

type PasswordByEmail struct {
	CodeLength int                     `json:"code_length" mapstructure:"code_length"`
	Template   RegisterByEmailTemplate `json:"template" mapstructure:"template"`
}

type RegisterByEmailTemplate struct {
	Subject     string `json:"subject" mapstructure:"subject"`
	ContentType string `json:"content_type" mapstructure:"content_type"`
	Body        string `json:"body" mapstructure:"body"`
}

func LoadConfig(in string) error {
	_c = &Config{}
	return util.LoadConfig(in, _c)
}

func Cfg() *Config {
	if _c == nil {
		log.Panicf("cfg nil")
	}
	return _c
}
