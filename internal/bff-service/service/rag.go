package service

import (
	knowledgeBase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	rag_service "github.com/UnicomAI/wanwu/api/proto/rag-service"
	safety_service "github.com/UnicomAI/wanwu/api/proto/safety-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/constant"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/gin-gonic/gin"
)

func CreateRag(ctx *gin.Context, userId, orgId string, req request.AppBriefConfig) (*request.RagReq, error) {
	resp, err := rag.CreateRag(ctx.Request.Context(), &rag_service.CreateRagReq{
		AppBrief: appBriefConfigModel2Proto(req),
		Identity: &rag_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &request.RagReq{
		RagID: resp.RagId,
	}, err
}

func UpdateRag(ctx *gin.Context, req request.RagBrief, userId, orgId string) error {
	_, err := rag.UpdateRag(ctx.Request.Context(), &rag_service.UpdateRagReq{
		RagId:    req.RagID,
		AppBrief: appBriefConfigModel2Proto(req.AppBriefConfig),
		Identity: &rag_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func UpdateRagConfig(ctx *gin.Context, req request.RagConfig) error {
	modelConfig, err := appModelConfigModel2Proto(req.ModelConfig)
	if err != nil {
		return err
	}
	rerankConfig, err := appModelConfigModel2Proto(req.RerankConfig)
	if err != nil {
		return err
	}
	qaRerankConfig, err := appModelConfigModel2Proto(req.QARerankConfig)
	if err != nil {
		return err
	}
	_, err = rag.UpdateRagConfig(ctx.Request.Context(), &rag_service.UpdateRagConfigReq{
		RagId:                 req.RagID,
		ModelConfig:           modelConfig,
		RerankConfig:          rerankConfig,
		QArerankConfig:        qaRerankConfig,
		KnowledgeBaseConfig:   ragKBConfigToProto(req.KnowledgeBaseConfig),
		QAknowledgeBaseConfig: ragQAKBConfigToProto(req.QAKnowledgeBaseConfig),
		SensitiveConfig:       ragSensitiveConfigToProto(req.SafetyConfig),
	})
	return err
}

func ragSensitiveConfigToProto(req request.AppSafetyConfig) *rag_service.RagSensitiveConfig {
	var sensitiveTableIds []string
	for _, v := range req.Tables {
		sensitiveTableIds = append(sensitiveTableIds, v.TableId)
	}
	sensitiveConfig := &rag_service.RagSensitiveConfig{
		Enable:   req.Enable,
		TableIds: sensitiveTableIds,
	}
	return sensitiveConfig
}

func ragKBConfigToProto(knowledgeConfig request.AppKnowledgebaseConfig) *rag_service.RagKnowledgeBaseConfig {
	result := &rag_service.RagKnowledgeBaseConfig{
		PerKnowledgeConfigs: make([]*rag_service.RagPerKnowledgeConfig, 0, len(knowledgeConfig.Knowledgebases)),
	}
	for _, knowledge := range knowledgeConfig.Knowledgebases {
		// 初始化单个知识库配置
		perConfig := &rag_service.RagPerKnowledgeConfig{
			KnowledgeId: knowledge.ID,
			GraphSwitch: knowledge.GraphSwitch,
		}
		// 构建元数据过滤条件（如果启用）
		if metaFilter := buildRagMetaFilter(knowledge.MetaDataFilterParams); metaFilter != nil {
			perConfig.RagMetaFilter = metaFilter
		}
		// 单个知识库配置添加到result
		result.PerKnowledgeConfigs = append(result.PerKnowledgeConfigs, perConfig)
	}
	result.GlobalConfig = buildRagGlobalConfig(knowledgeConfig.Config)
	return result
}

func ragQAKBConfigToProto(qaKnowledgeConfig request.AppQAKnowledgebaseConfig) *rag_service.RagQAKnowledgeBaseConfig {
	result := &rag_service.RagQAKnowledgeBaseConfig{
		PerKnowledgeConfigs: make([]*rag_service.RagPerQAKnowledgeConfig, 0, len(qaKnowledgeConfig.Knowledgebases)),
	}
	for _, knowledge := range qaKnowledgeConfig.Knowledgebases {
		// 初始化单个问答库配置
		perConfig := &rag_service.RagPerQAKnowledgeConfig{
			KnowledgeId: knowledge.ID,
		}
		// 构建元数据过滤条件（如果启用）
		if metaFilter := buildRagMetaFilter(knowledge.MetaDataFilterParams); metaFilter != nil {
			perConfig.RagMetaFilter = metaFilter
		}
		// 单个知识库配置添加到result
		result.PerKnowledgeConfigs = append(result.PerKnowledgeConfigs, perConfig)
	}
	result.GlobalConfig = buildRagQAGlobalConfig(qaKnowledgeConfig.Config)
	return result
}

// 构建单个知识库的元数据过滤条件
func buildRagMetaFilter(params *request.MetaDataFilterParams) *rag_service.RagMetaFilter {
	// 检查过滤参数是否有效（未启用则返回nil）
	if params == nil {
		return nil
	}
	if params.MetaFilterParams == nil {
		return &rag_service.RagMetaFilter{
			FilterEnable:    params.FilterEnable,
			FilterLogicType: params.FilterLogicType,
			FilterItems:     make([]*rag_service.RagMetaFilterItem, 0),
		}
	}
	// 转换过滤条件项
	filterItems := make([]*rag_service.RagMetaFilterItem, 0, len(params.MetaFilterParams))
	for _, metaParam := range params.MetaFilterParams {
		filterItems = append(filterItems, &rag_service.RagMetaFilterItem{
			Key:       metaParam.Key,
			Type:      metaParam.Type,
			Value:     metaParam.Value,
			Condition: metaParam.Condition,
		})
	}
	return &rag_service.RagMetaFilter{
		FilterEnable:    params.FilterEnable,
		FilterLogicType: params.FilterLogicType,
		FilterItems:     filterItems,
	}
}

func buildRagGlobalConfig(kbConfig request.AppKnowledgebaseParams) *rag_service.RagGlobalConfig {
	return &rag_service.RagGlobalConfig{
		MaxHistory:        kbConfig.MaxHistory,
		Threshold:         kbConfig.Threshold,
		TopK:              kbConfig.TopK,
		MatchType:         kbConfig.MatchType,
		KeywordPriority:   kbConfig.KeywordPriority,
		PriorityMatch:     kbConfig.PriorityMatch,
		SemanticsPriority: kbConfig.SemanticsPriority,
		TermWeight:        kbConfig.TermWeight,
		TermWeightEnable:  kbConfig.TermWeightEnable,
		UseGraph:          kbConfig.UseGraph,
		ChiChat:           kbConfig.ChiChat,
	}
}

func buildRagQAGlobalConfig(kbConfig request.AppQAKnowledgebaseParams) *rag_service.RagQAGlobalConfig {
	return &rag_service.RagQAGlobalConfig{
		MaxHistory:        kbConfig.MaxHistory,
		Threshold:         kbConfig.Threshold,
		TopK:              kbConfig.TopK,
		MatchType:         kbConfig.MatchType,
		KeywordPriority:   kbConfig.KeywordPriority,
		PriorityMatch:     kbConfig.PriorityMatch,
		SemanticsPriority: kbConfig.SemanticsPriority,
	}
}

func DeleteRag(ctx *gin.Context, req request.RagReq) error {
	_, err := rag.DeleteRag(ctx.Request.Context(), &rag_service.RagDeleteReq{
		RagId: req.RagID,
	})
	return err
}

func GetRag(ctx *gin.Context, req request.RagReq) (*response.RagInfo, error) {
	resp, err := rag.GetRagDetail(ctx.Request.Context(), &rag_service.RagDetailReq{RagId: req.RagID})
	if err != nil {
		return nil, err
	}
	modelConfig, rerankConfig, qaRerankConfig, err := appModelRerankProto2Model(ctx, resp)
	if err != nil {
		log.Errorf("ragId: %v gets config fail: %v", req.RagID, err.Error())
	}
	ragInfo := &response.RagInfo{
		RagID:                 resp.RagId,
		AppBriefConfig:        appBriefConfigProto2Model(ctx, resp.BriefConfig, constant.AppTypeRag),
		ModelConfig:           modelConfig,
		RerankConfig:          rerankConfig,
		QARerankConfig:        qaRerankConfig,
		KnowledgeBaseConfig:   ragKBConfigProto2Model(ctx, resp.KnowledgeBaseConfig),
		QAKnowledgeBaseConfig: ragKBQAConfigProto2Model(ctx, resp.QAknowledgeBaseConfig),
		SafetyConfig:          ragSafetyConfigProto2Model(ctx, resp.SensitiveConfig),
	}

	return ragInfo, nil
}

func appModelRerankProto2Model(ctx *gin.Context, resp *rag_service.RagInfo) (request.AppModelConfig, request.AppModelConfig, request.AppModelConfig, error) {
	var modelConfig, rerankConfig, qaRerankConfig request.AppModelConfig
	if resp.ModelConfig.ModelId != "" {
		modelInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: resp.ModelConfig.ModelId})
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, request.AppModelConfig{}, err
		}
		modelConfig, err = appModelConfigProto2Model(resp.ModelConfig, modelInfo.DisplayName)
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, request.AppModelConfig{}, err
		}
	}
	if resp.RerankConfig.ModelId != "" {
		rerankInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: resp.RerankConfig.ModelId})
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, request.AppModelConfig{}, err
		}
		rerankConfig, err = appModelConfigProto2Model(resp.RerankConfig, rerankInfo.DisplayName)
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, request.AppModelConfig{}, err
		}
	}
	if resp.QArerankConfig.ModelId != "" {
		qaRerankInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: resp.QArerankConfig.ModelId})
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, request.AppModelConfig{}, err
		}
		qaRerankConfig, err = appModelConfigProto2Model(resp.QArerankConfig, qaRerankInfo.DisplayName)
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, request.AppModelConfig{}, err
		}
	}
	return modelConfig, rerankConfig, qaRerankConfig, nil
}

func ragSafetyConfigProto2Model(ctx *gin.Context, sensitiveCfg *rag_service.RagSensitiveConfig) request.AppSafetyConfig {
	var sensitiveTableList []request.SensitiveTable
	if len(sensitiveCfg.GetTableIds()) != 0 {
		sensitiveWordTable, _ := safety.GetSensitiveWordTableListByIDs(ctx, &safety_service.GetSensitiveWordTableListByIDsReq{TableIds: sensitiveCfg.GetTableIds()})
		if sensitiveWordTable != nil {
			for _, table := range sensitiveWordTable.List {
				sensitiveTableList = append(sensitiveTableList, request.SensitiveTable{
					TableId:   table.TableId,
					TableName: table.TableName,
				})
			}
		}
	}
	safetyConfig := request.AppSafetyConfig{
		Enable: sensitiveCfg.Enable,
		Tables: sensitiveTableList,
	}
	return safetyConfig
}

func ragKBConfigProto2Model(ctx *gin.Context, kbConfig *rag_service.RagKnowledgeBaseConfig) request.AppKnowledgebaseConfig {
	if kbConfig == nil {
		return request.AppKnowledgebaseConfig{
			Knowledgebases: make([]request.AppKnowledgeBase, 0),
			Config:         request.AppKnowledgebaseParams{},
		}
	}
	knowledgeList := make([]request.AppKnowledgeBase, 0, len(kbConfig.PerKnowledgeConfigs))

	// 转换每个知识库的单独配置
	for _, perConfig := range kbConfig.PerKnowledgeConfigs {
		kbInfo, err := knowledgeBase.SelectKnowledgeDetailById(ctx, &knowledgeBase_service.KnowledgeDetailSelectReq{
			KnowledgeId: perConfig.KnowledgeId,
		})
		if err != nil {
			log.Errorf("select knowledge detail error: %v", err)
			return request.AppKnowledgebaseConfig{
				Knowledgebases: make([]request.AppKnowledgeBase, 0),
				Config:         request.AppKnowledgebaseParams{},
			}
		}
		// 基础信息映射
		knowledge := request.AppKnowledgeBase{
			ID:          perConfig.KnowledgeId,
			Name:        kbInfo.Name,
			GraphSwitch: kbInfo.GraphSwitch,
		}
		// 转换元数据过滤配置
		metaFilter := perConfig.RagMetaFilter
		knowledge.MetaDataFilterParams = convertRagMetaFilterToParams(metaFilter)

		knowledgeList = append(knowledgeList, knowledge)
	}
	globalConfig := kbConfig.GlobalConfig
	if globalConfig == nil {
		globalConfig = &rag_service.RagGlobalConfig{}
	}
	appConfig := request.AppKnowledgebaseParams{
		MaxHistory:        globalConfig.MaxHistory,
		Threshold:         globalConfig.Threshold,
		TopK:              globalConfig.TopK,
		MatchType:         globalConfig.MatchType,
		KeywordPriority:   globalConfig.KeywordPriority,
		PriorityMatch:     globalConfig.PriorityMatch,
		SemanticsPriority: globalConfig.SemanticsPriority,
		TermWeight:        globalConfig.TermWeight,
		TermWeightEnable:  globalConfig.TermWeightEnable,
		UseGraph:          globalConfig.UseGraph,
		ChiChat:           globalConfig.ChiChat,
	}
	return request.AppKnowledgebaseConfig{
		Knowledgebases: knowledgeList,
		Config:         appConfig,
	}
}

func ragKBQAConfigProto2Model(ctx *gin.Context, kbConfig *rag_service.RagQAKnowledgeBaseConfig) request.AppQAKnowledgebaseConfig {
	if kbConfig == nil {
		return request.AppQAKnowledgebaseConfig{
			Knowledgebases: make([]request.AppQAKnowledgeBase, 0),
			Config:         request.AppQAKnowledgebaseParams{},
		}
	}
	knowledgeList := make([]request.AppQAKnowledgeBase, 0, len(kbConfig.PerKnowledgeConfigs))

	// 转换每个问答库的单独配置
	for _, perConfig := range kbConfig.PerKnowledgeConfigs {
		kbInfo, err := knowledgeBase.SelectKnowledgeDetailById(ctx, &knowledgeBase_service.KnowledgeDetailSelectReq{
			KnowledgeId: perConfig.KnowledgeId,
		})
		if err != nil {
			log.Errorf("select qa detail error: %v", err)
			return request.AppQAKnowledgebaseConfig{
				Knowledgebases: make([]request.AppQAKnowledgeBase, 0),
				Config:         request.AppQAKnowledgebaseParams{},
			}
		}
		// 基础信息映射
		knowledge := request.AppQAKnowledgeBase{
			ID:   perConfig.KnowledgeId,
			Name: kbInfo.Name,
		}
		// 转换元数据过滤配置
		metaFilter := perConfig.RagMetaFilter
		knowledge.MetaDataFilterParams = convertRagMetaFilterToParams(metaFilter)

		knowledgeList = append(knowledgeList, knowledge)
	}
	globalConfig := kbConfig.GlobalConfig
	if globalConfig == nil {
		globalConfig = &rag_service.RagQAGlobalConfig{}
	}
	appConfig := request.AppQAKnowledgebaseParams{
		MaxHistory:        globalConfig.MaxHistory,
		Threshold:         globalConfig.Threshold,
		TopK:              globalConfig.TopK,
		MatchType:         globalConfig.MatchType,
		KeywordPriority:   globalConfig.KeywordPriority,
		PriorityMatch:     globalConfig.PriorityMatch,
		SemanticsPriority: globalConfig.SemanticsPriority,
	}
	return request.AppQAKnowledgebaseConfig{
		Knowledgebases: knowledgeList,
		Config:         appConfig,
	}
}

func convertRagMetaFilterToParams(metaFilter *rag_service.RagMetaFilter) *request.MetaDataFilterParams {
	if metaFilter == nil {
		return nil
	}
	// 转换过滤条件项
	filterParams := make([]*request.MetaFilterParams, 0, len(metaFilter.FilterItems))
	for _, item := range metaFilter.FilterItems {
		filterParams = append(filterParams, &request.MetaFilterParams{
			Key:       item.Key,
			Type:      item.Type,
			Value:     item.Value,
			Condition: item.Condition,
		})
	}
	return &request.MetaDataFilterParams{
		FilterEnable:     metaFilter.FilterEnable,
		FilterLogicType:  metaFilter.FilterLogicType,
		MetaFilterParams: filterParams, // 映射过滤条件列表
	}
}

func CopyRag(ctx *gin.Context, userId, orgId string, req request.RagReq) (*request.RagReq, error) {
	resp, err := rag.CopyRag(ctx.Request.Context(), &rag_service.CopyRagReq{
		Identity: &rag_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
		RagId: req.RagID,
	})
	if err != nil {
		return nil, err
	}
	return &request.RagReq{
		RagID: resp.RagId,
	}, err
}
