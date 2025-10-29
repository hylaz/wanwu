package config

import (
	"net/url"

	"github.com/UnicomAI/wanwu/pkg/i18n"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/minio"
	"github.com/UnicomAI/wanwu/pkg/util"
)

var (
	_c *Config
)

type Config struct {
	Server            ServerConfig            `json:"server" mapstructure:"server"`
	Log               LogConfig               `json:"log" mapstructure:"log"`
	JWT               JWTConfig               `json:"jwt" mapstructure:"jwt"`
	Decrypt           DecryptPasswd           `json:"decrypt-passwd" mapstructure:"decrypt-passwd"`
	I18n              i18n.Config             `json:"i18n" mapstructure:"i18n"`
	AssistantTemplate AssistantTemplateConfig `json:"assistant-template" mapstructure:"assistant-template"`
	CustomInfo        CustomInfoConfig        `json:"custom-info" mapstructure:"custom-info"`
	DocCenter         DocCenterConfig         `json:"doc-center" mapstructure:"doc-center"`
	DefaultIcon       DefaultIconConfig       `json:"default-icon" mapstructure:"default-icon"`
	// middleware
	Minio minio.Config `json:"minio" mapstructure:"minio"`
	// microservice
	Iam       ServiceConfig      `json:"iam" mapstructure:"iam"`
	Model     ModelConfig        `json:"model" mapstructure:"model"`
	MCP       ServiceConfig      `json:"mcp" mapstructure:"mcp"`
	App       ServiceConfig      `json:"app" mapstructure:"app"`
	Knowledge ServiceConfig      `json:"knowledge" mapstructure:"knowledge"`
	Rag       ServiceConfig      `json:"rag" mapstructure:"rag"`
	Assistant ServiceConfig      `json:"assistant" mapstructure:"assistant"`
	Operate   ServiceConfig      `json:"operate" mapstructure:"operate"`
	Agent     AgentServiceConfig `json:"agent" mapstructure:"agent"`

	Workflow           WorkflowServiceConfig           `json:"workflow" mapstructure:"workflow"`
	AgentScopeWorkFlow AgentScopeWorkFlowServiceConfig `json:"agentscope-workflow" mapstructure:"agentscope-workflow"`
	RagKnowledgeConfig RagKnowledgeConfig              `json:"rag-knowledge" mapstructure:"rag-knowledge"`
}

type ServerConfig struct {
	Host        string `json:"host" mapstructure:"host"`
	Port        int    `json:"port" mapstructure:"port"`
	WebBaseUrl  string `json:"web_base_url" mapstructure:"web_base_url"`
	ApiBaseUrl  string `json:"api_base_url" mapstructure:"api_base_url"`
	AppOpenUrl  string `json:"app_open_base_url" mapstructure:"app_open_base_url"`
	CallbackUrl string `json:"callback_url" mapstructure:"callback_url"`
}

type ModelConfig struct {
	Host            string `json:"host" mapstructure:"host"`
	PngTestFilePath string `json:"png_test_file_path" mapstructure:"png_test_file_path"`
	PdfTestFilePath string `json:"pdf_test_file_path" mapstructure:"pdf_test_file_path"`
}

type LogConfig struct {
	Std   bool         `json:"std" mapstructure:"std"`
	Level string       `json:"level" mapstructure:"level"`
	Logs  []log.Config `json:"logs" mapstructure:"logs"`
}

type JWTConfig struct {
	SigningKey string `json:"signing-key" mapstructure:"signing-key"`
}

type DecryptPasswd struct {
	IV  string `json:"iv" mapstructure:"iv"`
	Key string `json:"key" mapstructure:"key"`
}

type ServiceConfig struct {
	Host string `json:"host" mapstructure:"host"`
}

type WorkflowServiceConfig struct {
	Endpoint           string               `json:"endpoint" mapstructure:"endpoint"`
	MinioProxyEndpoint string               `json:"minio_proxy_endpoint" mapstructure:"minio_proxy_endpoint"`
	MinioProxyPrefix   string               `json:"minio_proxy_prefix" mapstructure:"minio_proxy_prefix"`
	ListUri            string               `json:"list_uri" mapstructure:"list_uri"`
	CreateUri          string               `json:"create_uri" mapstructure:"create_uri"`
	DeleteUri          string               `json:"delete_uri" mapstructure:"delete_uri"`
	CopyUri            string               `json:"copy_uri" mapstructure:"copy_uri"`
	ExportUri          string               `json:"export_uri" mapstructure:"export_uri"`
	ImportUri          string               `json:"import_uri" mapstructure:"import_uri"`
	TestRunUri         string               `json:"test_run_uri" mapstructure:"test_run_uri"`
	UploadActionUri    string               `json:"upload_action_uri" mapstructure:"upload_action_uri"`
	UploadCommonUri    string               `json:"upload_common_uri" mapstructure:"upload_common_uri"`
	SignImgUri         string               `json:"sign_img_uri" mapstructure:"sign_img_uri"`
	ModelParams        []WorkflowModelParam `json:"model_params" mapstructure:"model_params"`
}

type RagKnowledgeConfig struct {
	Endpoint               string `json:"endpoint" mapstructure:"endpoint"`
	ChatEndpoint           string `json:"chat-endpoint" mapstructure:"chat-endpoint"`
	SearchKnowledgeBaseUri string `json:"search-knowledge-base-uri" mapstructure:"search-knowledge-base-uri"`
	KnowledgeChatUri       string `json:"knowledge-chat-uri" mapstructure:"knowledge-chat-uri"`
}

type WorkflowModelParam struct {
	Name      string `json:"name" mapstructure:"name"`
	Desc      string `json:"desc" mapstructure:"desc"`
	Label     string `json:"label" mapstructure:"label"`
	Type      int    `json:"type" mapstructure:"type"`
	Precision int    `json:"precision" mapstructure:"precision"`
	Min       string `json:"min" mapstructure:"min"`
	Max       string `json:"max" mapstructure:"max"`

	ParamClass WorkflowModelParamClass      `json:"param_class" mapstructure:"param_class"`
	DefaultVal WorkflowModelParamDefaultVal `json:"default_val" mapstructure:"default_val"`
}

type WorkflowModelParamClass struct {
	ClassID int    `json:"class_id" mapstructure:"class_id"`
	Label   string `json:"label" mapstructure:"label"`
}

type WorkflowModelParamDefaultVal struct {
	Precise    string `json:"precise" mapstructure:"precise"`
	Balance    string `json:"balance" mapstructure:"balance"`
	Creative   string `json:"creative" mapstructure:"creative"`
	DefaultVal string `json:"default_val" mapstructure:"default_val"`
}

type AgentScopeWorkFlowServiceConfig struct {
	Endpoint                string `json:"endpoint" mapstructure:"endpoint"`
	WorkFlowListUri         string `json:"workflow_list_uri" mapstructure:"workflow_list_uri"`
	WorkFlowListUriInternal string `json:"workflow_list_uri_internal" mapstructure:"workflow_list_uri_internal"`
	DeleteWorkFlowUri       string `json:"delete_workflow_uri" mapstructure:"delete_workflow_uri"`
	PublishWorkFlowUri      string `json:"publish_workflow_uri" mapstructure:"publish_workflow_uri"`
	UnPublishWorkFlowUri    string `json:"unpublish_workflow_uri" mapstructure:"unpublish_workflow_uri"`
}

type AgentServiceConfig struct {
	Host           string    `json:"host" mapstructure:"host"`
	UploadMinioUri UriConfig `json:"upload_minio" mapstructure:"upload_minio"`
}

type UriConfig struct {
	Port string `json:"port" mapstructure:"port"`
	Uri  string `json:"uri" mapstructure:"uri"`
}

type AssistantTemplateConfig struct {
	ConfigPath string `json:"configPath" mapstructure:"configPath"`
}

type DocCenterConfig struct {
	FrontendPrefix string          `json:"frontend_prefix" mapstructure:"frontend_prefix"`
	Links          []DocLinkConfig `json:"links" mapstructure:"links"`
	docs           map[string]string
}

type DocLinkConfig struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type CustomInfoConfig struct {
	DefaultMode          string        `json:"default_mode" mapstructure:"default_mode"`
	Modes                []CustomTheme `json:"modes" mapstructure:"modes"`
	Version              string        `json:"version" mapstructure:"version"`
	RegisterByEmail      int           `json:"register_by_email" mapstructure:"register_by_email"`
	ResetPasswordByEmail int           `json:"reset_password_by_email" mapstructure:"reset_password_by_email"`
	LoginByEmail         int           `json:"login_by_email" mapstructure:"login_by_email"`
}

type CustomTheme struct {
	Mode  string      `json:"mode" mapstructure:"mode"`
	Login CustomLogin `json:"login" mapstructure:"login"`
	Home  CustomHome  `json:"home" mapstructure:"home"`
	Tab   CustomTab   `json:"tab" mapstructure:"tab"`
	About CustomAbout `json:"about" mapstructure:"about"`
}

type CustomLogin struct {
	BackgroundPath   string `json:"background_path" mapstructure:"background_path"`
	LogoPath         string `json:"logo_path" mapstructure:"logo_path"`
	LoginButtonColor string `json:"login_button_color" mapstructure:"login_button_color"`
	WelcomeText      string `json:"welcome_text" mapstructure:"welcome_text"`
	PlatformDesc     string `json:"platform_desc" mapstructure:"platform_desc"`
}

type CustomHome struct {
	LogoPath        string `json:"logo_path" mapstructure:"logo_path"`
	Title           string `json:"title" mapstructure:"title"`
	BackgroundColor string `json:"background_color" mapstructure:"background_color"`
}

type CustomTab struct {
	TabTitle    string `json:"title" mapstructure:"title"`
	TabLogoPath string `json:"logo_path" mapstructure:"logo_path"`
}

type CustomAbout struct {
	LogoPath  string `json:"logo_path" mapstructure:"logo_path"`
	Copyright string `json:"copyright" mapstructure:"copyright"`
}

type DefaultIconConfig struct {
	UserIcon     string `json:"user" mapstructure:"user"`
	RagIcon      string `json:"rag" mapstructure:"rag"`
	AgentIcon    string `json:"agent" mapstructure:"agent"`
	WorkflowIcon string `json:"workflow" mapstructure:"workflow"`
	ToolIcon     string `json:"tool" mapstructure:"tool"`
}

func LoadConfig(in string) error {
	_c = &Config{}
	if err := util.LoadConfig(in, _c); err != nil {
		return err
	}
	_c.DocCenter.docs = make(map[string]string)
	for _, link := range _c.DocCenter.Links {
		url, _ := url.JoinPath(_c.Server.WebBaseUrl, _c.DocCenter.FrontendPrefix, url.PathEscape(link.Val))
		_c.DocCenter.docs[link.Key] = url
	}
	return nil
}

func Cfg() *Config {
	if _c == nil {
		log.Panicf("cfg nil")
	}
	return _c
}

// GetDocs 返回 docs 的深拷贝
func (d *DocCenterConfig) GetDocs() map[string]string {
	if d.docs == nil {
		return nil
	}
	// 深拷贝
	result := make(map[string]string, len(d.docs))
	for k, v := range d.docs {
		result[k] = v
	}
	return result
}
