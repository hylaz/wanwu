package mp

// model type
const (
	ModelTypeLLM       = "llm"
	ModelTypeEmbedding = "embedding"
	ModelTypeRerank    = "rerank"
	ModelTypeOcr       = "ocr"
	ModelTypeGui       = "gui"
	ModelTypePdfParser = "pdf-parser"
)

// model provider
const (
	ProviderOpenAICompatible = "OpenAI-API-compatible"
	ProviderYuanJing         = "YuanJing"
	ProviderHuoshan          = "Huoshan"
	ProviderOllama           = "Ollama"
	ProviderQwen             = "Qwen"
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
