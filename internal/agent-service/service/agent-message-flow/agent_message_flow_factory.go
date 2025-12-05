package agent_message_flow

import (
	"context"

	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-flow/nodes"
	"github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-flow/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

const (
	keyOfPersonRender        = "persona_render"
	keyOfKnowledgeRetriever  = "knowledge_retriever"
	keyOfUploadFileRetriever = "upload_file_retriever"
	keyOfPromptVariables     = "prompt_variables"
	keyOfPromptTemplate      = "prompt_template"
	//keyOfToolsPreRetriever      = "tools_pre_retriever"
)

type AgentState struct {
	Messages                 []*schema.Message
	UserInput                *schema.Message
	ReturnDirectlyToolCallID string
}

// NewAgentMessageFlow 创建一个图信息 这个要确认是否是线程安全的，要不每次都创建一个图信息，效率会有问题，如果不是线程安全的，考虑从自己实现一个简单的责任链
func NewAgentMessageFlow() *compose.Graph[*request.AgentChatContext, []*schema.Message] {
	//创建图信息
	graph := compose.NewGraph[*request.AgentChatContext, []*schema.Message](
		compose.WithGenLocalState(func(ctx context.Context) (state *AgentState) {
			return &AgentState{}
		}))
	//增加图节点
	AddAgentFlowNode(graph)
	//连接图信息
	BuildGraphEdge(graph)
	return graph
}

func AddAgentFlowNode(graph *compose.Graph[*request.AgentChatContext, []*schema.Message]) {

	//_ = graph.AddLambdaNode(keyOfPersonRender,
	//	compose.InvokableLambda[*AgentRequest, string](personaVars.RenderPersona),
	//	compose.WithStatePreHandler(func(ctx context.Context, ar *request.AgentChatContext, state *AgentState) (*request.AgentChatContext, error) {
	//		state.UserInput = ar.Input
	//		return ar, nil
	//	}),
	//	compose.WithOutputKey(prompt.placeholderOfPersona))

	promptVars := nodes.PromptVariables{}
	_ = graph.AddLambdaNode(keyOfPromptVariables,
		compose.InvokableLambda[*request.AgentChatContext, map[string]any](promptVars.AssemblePromptVariables))

	kr := nodes.KnowledgeRetriever{}
	_ = graph.AddLambdaNode(keyOfKnowledgeRetriever,
		compose.InvokableLambda[*request.AgentChatContext, string](kr.Retrieve),
		compose.WithNodeName(keyOfKnowledgeRetriever),
		compose.WithOutputKey(prompt.PlaceholderOfKnowledge))

	//fr := nodes.UploadFileRetriever{}
	//_ = graph.AddLambdaNode(keyOfUploadFileRetriever,
	//	compose.InvokableLambda[*request.AgentChatContext, string](fr.Retrieve),
	//	compose.WithNodeName(keyOfUploadFileRetriever),
	//	compose.WithOutputKey(prompt.PlaceholderOfUploadFile))
	//_ = graph.AddLambdaNode(keyOfToolsPreRetriever,
	//	compose.InvokableLambda[*AgentRequest, []*schema.Message](tr.toolPreRetrieve),
	//	compose.WithOutputKey(keyOfToolsPreRetriever),
	//	compose.WithNodeName(keyOfToolsPreRetriever),
	//)
	_ = graph.AddChatTemplateNode(keyOfPromptTemplate, nodes.ChatPrompt)
}

func BuildGraphEdge(graph *compose.Graph[*request.AgentChatContext, []*schema.Message]) {
	//_ = graph.AddEdge(compose.START, keyOfPersonRender)
	_ = graph.AddEdge(compose.START, keyOfPromptVariables)
	_ = graph.AddEdge(compose.START, keyOfKnowledgeRetriever)
	//_ = graph.AddEdge(compose.START, keyOfUploadFileRetriever)
	//_ = graph.AddEdge(compose.START, keyOfToolsPreRetriever)

	//_ = graph.AddEdge(keyOfPersonRender, keyOfPromptTemplate)
	_ = graph.AddEdge(keyOfPromptVariables, keyOfPromptTemplate)
	_ = graph.AddEdge(keyOfKnowledgeRetriever, keyOfPromptTemplate)
	//_ = graph.AddEdge(keyOfUploadFileRetriever, keyOfPromptTemplate)
	//_ = graph.AddEdge(keyOfToolsPreRetriever, keyOfPromptTemplate)

	//添加结束节点
	_ = graph.AddEdge(keyOfPromptTemplate, compose.END)
}
