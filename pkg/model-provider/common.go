package mp

// model type
const (
	ModelTypeLLM        = "llm"
	ModelTypeEmbedding  = "embedding"
	ModelTypeRerank     = "rerank"
	ModelTypeOcr        = "ocr"
	ModelTypeGui        = "gui"
	ModelTypePdfParser  = "pdf-parser"
	ModelTypeAsr        = "asr"
	ModelTypeText2Image = "text2image"
)

// model provider
const (
	ProviderOpenAICompatible = "OpenAI-API-compatible" // openai
	ProviderYuanJing         = "YuanJing"              //
	ProviderHuoshan          = "Huoshan"               //火山
	ProviderOllama           = "Ollama"
	ProviderQwen             = "Qwen" //通义大模型
	ProviderInfini           = "Infini"
)

var (
	_callbackUrl string
)

func Init(callbackUrl string) {
	if _callbackUrl != "" {
		panic("model provider already init")
	}
	_callbackUrl = callbackUrl
}
