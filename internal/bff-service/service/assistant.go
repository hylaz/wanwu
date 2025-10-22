package service

import (
	"encoding/json"
	"sort"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	knowledgeBase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	safety_service "github.com/UnicomAI/wanwu/api/proto/safety-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	bff_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/util"
	"github.com/UnicomAI/wanwu/pkg/constant"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

func AssistantCreate(ctx *gin.Context, userId, orgId string, req request.AppBriefConfig) (*response.AssistantCreateResp, error) {
	resp, err := assistant.AssistantCreate(ctx.Request.Context(), &assistant_service.AssistantCreateReq{
		AssistantBrief: appBriefConfigModel2Proto(req),
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &response.AssistantCreateResp{
		AssistantId: resp.AssistantId,
	}, nil
}

func AssistantUpdate(ctx *gin.Context, userId, orgId string, req request.AssistantBrief) (interface{}, error) {
	_, err := assistant.AssistantUpdate(ctx.Request.Context(), &assistant_service.AssistantUpdateReq{
		AssistantId:    req.AssistantId,
		AssistantBrief: appBriefConfigModel2Proto(req.AppBriefConfig),
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return nil, err
}

func AssistantConfigUpdate(ctx *gin.Context, userId, orgId string, req request.AssistantConfig) (interface{}, error) {
	modelConfig, err := appModelConfigModel2Proto(req.ModelConfig)
	if err != nil {
		return nil, err
	}
	rerankConfig, err := appModelConfigModel2Proto(req.RerankConfig)
	if err != nil {
		return nil, err
	}
	_, err = assistant.AssistantConfigUpdate(ctx.Request.Context(), &assistant_service.AssistantConfigUpdateReq{
		AssistantId:         req.AssistantId,
		Prologue:            req.Prologue,
		Instructions:        req.Instructions,
		RecommendQuestion:   req.RecommendQuestion,
		ModelConfig:         modelConfig,
		KnowledgeBaseConfig: transKnowledgebases2Proto(req.KnowledgeBaseConfig),
		RerankConfig:        rerankConfig,
		OnlineSearchConfig: &assistant_service.AssistantOnlineSearchConfig{
			SearchUrl:      req.OnlineSearchConfig.SearchUrl,
			SearchKey:      req.OnlineSearchConfig.SearchKey,
			Enable:         req.OnlineSearchConfig.Enable,
			SearchRerankId: req.OnlineSearchConfig.SearchRerankId,
		},
		SafetyConfig: &assistant_service.AssistantSafetyConfig{
			Enable:         req.SafetyConfig.Enable,
			SensitiveTable: transSafetyConfig2Proto(req.SafetyConfig.Tables),
		},
		VisionConfig: &assistant_service.AssistantVisionConfig{
			PicNum: req.VisionConfig.PicNum,
		},
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})

	return nil, err
}

func GetAssistantInfo(ctx *gin.Context, userId, orgId string, req request.AssistantIdRequest) (*response.Assistant, error) {
	resp, err := assistant.GetAssistantInfo(ctx.Request.Context(), &assistant_service.GetAssistantInfoReq{
		AssistantId: req.AssistantId,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	if err != nil {
		return nil, err
	}
	return transAssistantResp2Model(ctx, resp)
}

func AssistantWorkFlowCreate(ctx *gin.Context, userId, orgId string, req request.AssistantWorkFlowAddRequest) error {
	_, err := assistant.AssistantWorkFlowCreate(ctx.Request.Context(), &assistant_service.AssistantWorkFlowCreateReq{
		AssistantId: req.AssistantId,
		WorkFlowId:  req.WorkFlowId,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func AssistantWorkFlowDelete(ctx *gin.Context, userId, orgId string, req request.AssistantWorkFlowDelRequest) error {
	_, err := assistant.AssistantWorkFlowDelete(ctx.Request.Context(), &assistant_service.AssistantWorkFlowDeleteReq{
		AssistantId: req.AssistantId,
		WorkFlowId:  req.WorkFlowId,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func AssistantWorkFlowEnableSwitch(ctx *gin.Context, userId, orgId string, req request.AssistantWorkFlowToolEnableRequest) error {
	_, err := assistant.AssistantWorkFlowEnableSwitch(ctx.Request.Context(), &assistant_service.AssistantWorkFlowEnableSwitchReq{
		AssistantId: req.AssistantId,
		WorkFlowId:  req.WorkFlowId,
		Enable:      req.Enable,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func AssistantMCPCreate(ctx *gin.Context, userId, orgId string, req request.AssistantMCPToolAddRequest) error {
	_, err := assistant.AssistantMCPCreate(ctx.Request.Context(), &assistant_service.AssistantMCPCreateReq{
		AssistantId: req.AssistantId,
		McpId:       req.MCPId,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func AssistantMCPDelete(ctx *gin.Context, userId, orgId string, req request.AssistantMCPToolDelRequest) error {
	_, err := assistant.AssistantMCPDelete(ctx.Request.Context(), &assistant_service.AssistantMCPDeleteReq{
		AssistantId: req.AssistantId,
		McpId:       req.MCPId,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func AssistantMCPEnableSwitch(ctx *gin.Context, userId, orgId string, req request.AssistantMCPToolEnableRequest) error {
	_, err := assistant.AssistantMCPEnableSwitch(ctx.Request.Context(), &assistant_service.AssistantMCPEnableSwitchReq{
		AssistantId: req.AssistantId,
		McpId:       req.MCPId,
		Enable:      req.Enable,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func AssistantCustomToolCreate(ctx *gin.Context, userId, orgId string, req request.AssistantCustomToolAddRequest) error {
	_, err := assistant.AssistantCustomToolCreate(ctx.Request.Context(), &assistant_service.AssistantCustomToolCreateReq{
		AssistantId:  req.AssistantId,
		CustomToolId: req.CustomToolId,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func AssistantCustomToolDelete(ctx *gin.Context, userId, orgId string, req request.AssistantCustomToolDelRequest) error {
	_, err := assistant.AssistantCustomToolDelete(ctx.Request.Context(), &assistant_service.AssistantCustomToolDeleteReq{
		AssistantId:  req.AssistantId,
		CustomToolId: req.CustomToolId,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func AssistantCustomToolEnableSwitch(ctx *gin.Context, userId, orgId string, req request.AssistantCustomToolEnableRequest) error {
	_, err := assistant.AssistantCustomToolEnableSwitch(ctx.Request.Context(), &assistant_service.AssistantCustomToolEnableSwitchReq{
		AssistantId:  req.AssistantId,
		CustomToolId: req.CustomToolId,
		Enable:       req.Enable,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	return err
}

func assistantMCPConvert(ctx *gin.Context, assistantMCPInfos []*assistant_service.AssistantMCPInfos) ([]*response.MCPInfos, error) {
	// 若查询结果为空，返回空列表
	if len(assistantMCPInfos) == 0 {
		return nil, nil
	}

	// 提取MCP ID列表
	var mcpIds []string
	for _, m := range assistantMCPInfos {
		mcpIds = append(mcpIds, m.McpId)
	}

	// 批量查询MCP详情
	mcpResp, err := mcp.GetCustomMCPByMCPIdList(ctx.Request.Context(), &mcp_service.GetCustomMCPByMCPIdListReq{
		McpIdList: mcpIds,
	})

	// 构建MCP详情映射
	mcpDetailMap := make(map[string]*mcp_service.CustomMCPInfo)
	if err == nil && mcpResp != nil { // 仅当查询成功且响应有效时才构建映射
		for _, item := range mcpResp.Infos {
			mcpDetailMap[item.McpId] = item
		}
	}

	// 构建返回结果
	var retMCPInfos []*response.MCPInfos
	for _, m := range assistantMCPInfos {
		item, exists := mcpDetailMap[m.McpId]
		if exists {
			// 有效MCP
			retMCPInfos = append(retMCPInfos, &response.MCPInfos{
				MCPId:         m.McpId,
				UniqueId:      bff_util.ConcatAssistantToolUniqueId("mcp", m.McpId),
				MCPSquareId:   item.Info.McpSquareId,
				Enable:        m.Enable,
				MCPName:       item.Info.Name,
				MCPDesc:       item.Info.Desc,
				MCPServerFrom: item.Info.From,
				MCPServerUrl:  item.SseUrl,
			})
		}
	}

	// 即使详情查询失败，也返回组装后的结果
	return retMCPInfos, nil
}

func assistantCustomConvert(ctx *gin.Context, assistantCustomInfos []*assistant_service.AssistantCustomToolInfos) ([]*response.CustomInfos, error) {
	// 若查询为空，返回空列表
	if len(assistantCustomInfos) == 0 {
		return nil, nil
	}

	// 提取自定义工具ID列表
	var customToolIds []string
	for _, c := range assistantCustomInfos {
		customToolIds = append(customToolIds, c.CustomToolId)
	}

	// 批量查询自定义工具详情
	mcpResp, err := mcp.GetCustomToolByCustomToolIdList(ctx.Request.Context(), &mcp_service.GetCustomToolByCustomToolIdListReq{
		CustomToolIdList: customToolIds,
	})

	// 构建ID到工具信息的映射
	customToolMap := make(map[string]*mcp_service.GetCustomToolItem)
	if err == nil && mcpResp != nil { // 仅当查询成功且响应有效时才构建映射
		for _, item := range mcpResp.List {
			customToolMap[item.CustomToolId] = item
		}
	}

	// 组装返回结果
	var retCustomInfos []*response.CustomInfos
	for _, c := range assistantCustomInfos {
		item, exists := customToolMap[c.CustomToolId]
		if exists {
			// 有效工具
			retCustomInfos = append(retCustomInfos, &response.CustomInfos{
				CustomId:   c.CustomToolId,
				UniqueId:   bff_util.ConcatAssistantToolUniqueId("custom", c.CustomToolId),
				Enable:     c.Enable,
				CustomName: item.Name,
				CustomDesc: item.Description,
			})
		}
	}

	return retCustomInfos, nil
}

func ConversationCreate(ctx *gin.Context, userId, orgId string, req request.ConversationCreateRequest) (response.ConversationCreateResp, error) {
	resp, err := assistant.ConversationCreate(ctx.Request.Context(), &assistant_service.ConversationCreateReq{
		AssistantId: req.AssistantId,
		Prompt:      req.Prompt,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	if err != nil {
		return response.ConversationCreateResp{}, err
	}
	return response.ConversationCreateResp{
		ConversationId: resp.ConversationId,
	}, nil
}

func ConversationDelete(ctx *gin.Context, userId, orgId string, req request.ConversationIdRequest) (interface{}, error) {
	_, err := assistant.ConversationDelete(ctx.Request.Context(), &assistant_service.ConversationDeleteReq{
		ConversationId: req.ConversationId,
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func GetConversationList(ctx *gin.Context, userId, orgId string, req request.ConversationGetListRequest) (response.PageResult, error) {
	resp, err := assistant.GetConversationList(ctx.Request.Context(), &assistant_service.GetConversationListReq{
		AssistantId: req.AssistantId,
		PageSize:    int32(req.PageSize),
		PageNo:      int32(req.PageNo),
		Identity: &assistant_service.Identity{
			UserId: userId,
		},
	})
	if err != nil {
		return response.PageResult{}, err
	}
	return response.PageResult{Total: resp.Total, List: resp.Data}, nil
}

func GetConversationDetailList(ctx *gin.Context, userId, orgId string, req request.ConversationGetDetailListRequest) (response.PageResult, error) {
	resp, err := assistant.GetConversationDetailList(ctx.Request.Context(), &assistant_service.GetConversationDetailListReq{
		ConversationId: req.ConversationId,
		PageSize:       int32(req.PageSize),
		PageNo:         int32(req.PageNo),
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	if err != nil {
		return response.PageResult{}, err
	}

	// 转换resp.Data为自定义的ConversionDetailInfo结构体数组
	var convertedList []response.ConversationDetailInfo
	for _, item := range resp.Data {
		convertedItem := response.ConversationDetailInfo{
			Id:             item.Id,
			AssistantId:    item.AssistantId,
			ConversationId: item.ConversationId,
			Prompt:         item.Prompt,
			SysPrompt:      item.SysPrompt,
			Response:       item.Response,
			QaType:         item.QaType,
			CreatedBy:      item.CreatedBy,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
			RequestFiles:   transRequestFiles(item.RequestFiles),
			FileSize:       item.FileSize,
			FileName:       item.FileName,
		}

		// 将SearchList从string转换为interface{}
		if item.SearchList != "" {
			var searchList interface{}
			if err := json.Unmarshal([]byte(item.SearchList), &searchList); err != nil {
				log.Warnf("解析SearchList失败，使用原始字符串，error: %v, searchList: %s", err, item.SearchList)
				convertedItem.SearchList = item.SearchList
			} else {
				convertedItem.SearchList = searchList
			}
		} else {
			convertedItem.SearchList = nil
		}

		convertedList = append(convertedList, convertedItem)

		// 对切片进行排序
		sort.Slice(convertedList, func(i, j int) bool {
			// CreatedAt值小的时间更早，排在前面
			return convertedList[i].CreatedAt < convertedList[j].CreatedAt
		})
	}

	return response.PageResult{Total: resp.Total, List: convertedList, PageNo: req.PageNo, PageSize: req.PageSize}, nil
}

func transKnowledgebases2Proto(kbConfig request.AppKnowledgebaseConfig) *assistant_service.AssistantKnowledgeBaseConfig {
	var knowIds []string
	if len(kbConfig.Knowledgebases) > 0 {
		for _, v := range kbConfig.Knowledgebases {
			knowIds = append(knowIds, v.ID)
		}
	}
	return &assistant_service.AssistantKnowledgeBaseConfig{
		KnowledgeBaseIds:     knowIds,
		MaxHistory:           kbConfig.Config.MaxHistory,
		Threshold:            kbConfig.Config.Threshold,
		TopK:                 kbConfig.Config.TopK,
		MatchType:            kbConfig.Config.MatchType,
		KeywordPriority:      kbConfig.Config.KeywordPriority,
		PriorityMatch:        kbConfig.Config.PriorityMatch,
		SemanticsPriority:    kbConfig.Config.SemanticsPriority,
		AppKnowledgeBaseList: transKnowledgeParams(kbConfig.Knowledgebases),
	}
}

func transKnowledgeParams(paramsList []request.AppKnowledgeBase) []*assistant_service.AppKnowledgeBase {
	if len(paramsList) == 0 {
		return nil
	}
	var retList []*assistant_service.AppKnowledgeBase
	for _, base := range paramsList {
		retList = append(retList, &assistant_service.AppKnowledgeBase{
			KnowledgeBaseId:      base.ID,
			KnowledgeBaseName:    base.Name,
			MetaDataFilterParams: transKnowledgeMetaParams(base.MetaDataFilterParams),
		})
	}
	return retList
}

func transKnowledgeMetaParams(baseInfo *request.MetaDataFilterParams) *assistant_service.MetaDataFilterParams {
	if baseInfo == nil {
		return nil
	}
	return &assistant_service.MetaDataFilterParams{
		FilterEnable:     baseInfo.FilterEnable,
		FilterLogicType:  baseInfo.FilterLogicType,
		MetaFilterParams: transMetaFilterParams(baseInfo.MetaFilterParams),
	}
}

func transMetaFilterParams(metaFilterList []*request.MetaFilterParams) []*assistant_service.MetaFilterParams {
	if metaFilterList == nil {
		return nil
	}
	var metaList []*assistant_service.MetaFilterParams
	for _, m := range metaFilterList {
		metaList = append(metaList, &assistant_service.MetaFilterParams{
			Condition: m.Condition,
			Key:       m.Key,
			Type:      m.Type,
			Value:     m.Value,
		})
	}
	return metaList
}
func transSafetyConfig2Proto(tables []request.SensitiveTable) []*assistant_service.SensitiveTable {
	if tables == nil {
		return nil
	}
	result := make([]*assistant_service.SensitiveTable, 0, len(tables))
	for _, table := range tables {
		result = append(result, &assistant_service.SensitiveTable{
			TableId:   table.TableId,
			TableName: table.TableName,
		})
	}
	return result
}

func transAssistantResp2Model(ctx *gin.Context, resp *assistant_service.AssistantInfo) (*response.Assistant, error) {
	log.Debugf("开始转换Assistant响应到模型，响应内容: %+v", resp)
	if resp == nil {
		log.Debugf("Assistant响应为空，返回空Assistant模型")
		return nil, nil
	}
	var modelConfig request.AppModelConfig
	if resp.ModelConfig != nil && resp.ModelConfig.ModelId != "" {
		log.Debugf("检测到模型配置，模型ID: %s", resp.ModelConfig.ModelId)
		modelInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: resp.ModelConfig.ModelId})
		if err != nil {
			log.Errorf("获取模型信息失败，模型ID: %s, 错误: %v", resp.ModelConfig.ModelId, err)
		}
		if modelInfo != nil {
			modelConfig, err = appModelConfigProto2Model(resp.ModelConfig, modelInfo.DisplayName)
			if err != nil {
				log.Errorf("模型配置Proto转换到模型失败，模型ID: %s, 错误: %v", resp.ModelConfig.ModelId, err)
				return nil, err
			}
			log.Debugf("模型配置转换成功: %+v", modelConfig)
		}
	} else {
		log.Debugf("模型配置为空或模型ID为空")
	}
	var rerankConfig request.AppModelConfig
	if resp.RerankConfig != nil && resp.RerankConfig.ModelId != "" {
		log.Debugf("检测到Rerank配置，模型ID: %s", resp.RerankConfig.ModelId)
		modelInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: resp.RerankConfig.ModelId})
		if err != nil {
			log.Errorf("获取Rerank模型信息失败，模型ID: %s, 错误: %v", resp.RerankConfig.ModelId, err)
			return nil, err
		}
		rerankConfig, err = appModelConfigProto2Model(resp.RerankConfig, modelInfo.DisplayName)
		if err != nil {
			log.Errorf("Rerank配置Proto转换到模型失败，模型ID: %s, 错误: %v", resp.RerankConfig.ModelId, err)
			return nil, err
		}
		log.Debugf("Rerank配置转换成功: %+v", rerankConfig)
	} else {
		log.Debugf("Rerank配置为空或模型ID为空")
	}

	var assistantWorkFlowInfos []*response.WorkFlowInfos
	if len(resp.WorkFlowInfos) > 0 {
		var workflowIds []string
		for _, wf := range resp.WorkFlowInfos {
			workflowIds = append(workflowIds, wf.WorkFlowId)
		}
		cozeWorkflowList, err := ListWorkflowByIDs(ctx, "", workflowIds)
		if err != nil {
			return nil, err
		}
		for _, wf := range resp.WorkFlowInfos {
			workFlowInfo := &response.WorkFlowInfos{
				WorkFlowId: wf.WorkFlowId,
				ApiName:    wf.ApiName,
				Enable:     wf.Enable,
				UniqueId:   bff_util.ConcatAssistantToolUniqueId("workflow", wf.WorkFlowId),
			}

			for _, info := range cozeWorkflowList.Workflows {
				if info.WorkflowId == wf.WorkFlowId {
					// 找到匹配的工作流，设置名称和描述
					workFlowInfo.WorkFlowName = info.Name
					workFlowInfo.WorkFlowDesc = info.Desc
				}
			}

			assistantWorkFlowInfos = append(assistantWorkFlowInfos, workFlowInfo)
			log.Debugf("添加工作流信息: WorkFlowId=%s, ApiName=%s", wf.WorkFlowId, wf.ApiName)
		}
		log.Debugf("总共添加 %d 个工作流信息", len(assistantWorkFlowInfos))
	} else {
		log.Debugf("工作流信息为空")
	}

	// 查询该用户所有权限的所有 MCP
	assistantMCPInfos, err := assistantMCPConvert(ctx, resp.McpInfos)
	if err != nil {
		return nil, err
	}
	// 查询该用户所有权限的 Custom
	assistantCustomInfos, err := assistantCustomConvert(ctx, resp.CustomToolInfos)
	if err != nil {
		return nil, err
	}

	var onlineSearchConfig request.OnlineSearchConfig
	if resp.OnlineSearchConfig != nil {
		onlineSearchConfig = request.OnlineSearchConfig{
			SearchUrl:      resp.OnlineSearchConfig.SearchUrl,
			SearchKey:      resp.OnlineSearchConfig.SearchKey,
			Enable:         resp.OnlineSearchConfig.Enable,
			SearchRerankId: resp.OnlineSearchConfig.SearchRerankId,
		}
	}
	var sensitiveWordTable *safety_service.SensitiveWordTables
	if len(resp.SafetyConfig.GetSensitiveTable()) != 0 {
		var tableIds []string
		for _, table := range resp.SafetyConfig.SensitiveTable {
			tableIds = append(tableIds, table.TableId)
		}
		sensitiveWordTable, _ = safety.GetSensitiveWordTableListByIDs(ctx, &safety_service.GetSensitiveWordTableListByIDsReq{TableIds: tableIds})
	}
	knowledgeBaseConfig, err := transKnowledgeBases2Model(ctx, resp.KnowledgeBaseConfig)
	if err != nil {
		return nil, err
	}
	var visionConfig response.VisionConfig
	if resp.VisionConfig != nil {
		visionConfig = response.VisionConfig{
			MaxPicNum: resp.VisionConfig.MaxPicNum,
			PicNum:    resp.VisionConfig.PicNum,
		}
	}
	assistantModel := response.Assistant{
		AssistantId:         resp.AssistantId,
		AppBriefConfig:      appBriefConfigProto2Model(ctx, resp.AssistantBrief, constant.AppTypeAgent),
		Prologue:            resp.Prologue,
		Instructions:        resp.Instructions,
		RecommendQuestion:   resp.RecommendQuestion,
		KnowledgeBaseConfig: knowledgeBaseConfig,
		ModelConfig:         modelConfig,
		RerankConfig:        rerankConfig,
		OnlineSearchConfig:  onlineSearchConfig,
		SafetyConfig:        request.AppSafetyConfig{Enable: resp.SafetyConfig.GetEnable()},
		VisionConfig:        visionConfig,
		Scope:               resp.Scope,
		WorkFlowInfos:       assistantWorkFlowInfos,
		MCPInfos:            assistantMCPInfos,
		CustomInfos:         assistantCustomInfos,
		CreatedAt:           util.Time2Str(resp.CreatTime),
		UpdatedAt:           util.Time2Str(resp.UpdateTime),
	}
	if sensitiveWordTable != nil {
		var sensitiveTableList []request.SensitiveTable
		for _, table := range sensitiveWordTable.List {
			sensitiveTableList = append(sensitiveTableList, request.SensitiveTable{
				TableId:   table.TableId,
				TableName: table.TableName,
			})
		}
		assistantModel.SafetyConfig.Tables = sensitiveTableList
	}
	log.Debugf("Assistant响应到模型转换完成，结果: %+v", assistantModel)
	return &assistantModel, nil
}

func transKnowledgeBases2Model(ctx *gin.Context, kbConfig *assistant_service.AssistantKnowledgeBaseConfig) (request.AppKnowledgebaseConfig, error) {
	if kbConfig == nil {
		log.Debugf("知识库配置为空")
		return request.AppKnowledgebaseConfig{
			Knowledgebases: make([]request.AppKnowledgeBase, 0),
		}, nil
	}
	if len(kbConfig.KnowledgeBaseIds) == 0 {
		log.Debugf("知识库配置为空")
		return request.AppKnowledgebaseConfig{
			Knowledgebases: make([]request.AppKnowledgeBase, 0),
		}, nil
	}

	// 获取知识库详情列表
	kbInfoList, err := knowledgeBase.SelectKnowledgeDetailByIdList(ctx, &knowledgeBase_service.KnowledgeDetailSelectListReq{
		KnowledgeIds: kbConfig.KnowledgeBaseIds,
	})
	if err != nil {
		return request.AppKnowledgebaseConfig{
			Knowledgebases: make([]request.AppKnowledgeBase, 0),
		}, err
	}

	knowledgeBases := buildKnowledgeBases(kbInfoList, kbConfig.KnowledgeBaseIds, kbConfig.AppKnowledgeBaseList)

	return request.AppKnowledgebaseConfig{
		Knowledgebases: knowledgeBases,
		Config: request.AppKnowledgebaseParams{
			MaxHistory:        kbConfig.MaxHistory,
			Threshold:         kbConfig.Threshold,
			TopK:              kbConfig.TopK,
			MatchType:         kbConfig.MatchType,
			PriorityMatch:     kbConfig.PriorityMatch,
			SemanticsPriority: kbConfig.SemanticsPriority,
			KeywordPriority:   kbConfig.KeywordPriority,
		},
	}, nil

}

func buildKnowledgeBases(kbInfoList *knowledgeBase_service.KnowledgeDetailSelectListResp, kbIdList []string, kbConfigList []*assistant_service.AppKnowledgeBase) []request.AppKnowledgeBase {
	if len(kbInfoList.List) == 0 {
		return make([]request.AppKnowledgeBase, 0)
	}
	var knowledgeMap = make(map[string]*knowledgeBase_service.KnowledgeInfo)
	for _, kbInfo := range kbInfoList.List {
		knowledgeMap[kbInfo.KnowledgeId] = kbInfo
	}
	var knowledgeBases = make([]request.AppKnowledgeBase, 0)
	if len(kbConfigList) > 0 {
		for _, kbConfig := range kbConfigList {
			info := knowledgeMap[kbConfig.KnowledgeBaseId]
			if info == nil {
				continue
			}
			params := buildAssistantMetaDataFilterParams(kbConfig)
			knowledgeBases = append(knowledgeBases, request.AppKnowledgeBase{
				ID:                   kbConfig.KnowledgeBaseId,
				Name:                 info.Name,
				MetaDataFilterParams: params,
			})
		}
	} else {
		for _, kbId := range kbIdList {
			info := knowledgeMap[kbId]
			if info == nil {
				continue
			}
			knowledgeBases = append(knowledgeBases, request.AppKnowledgeBase{
				ID:   kbId,
				Name: info.Name,
			})
		}
	}

	return knowledgeBases
}

func buildAssistantMetaDataFilterParams(kbConfig *assistant_service.AppKnowledgeBase) *request.MetaDataFilterParams {
	params := kbConfig.MetaDataFilterParams
	if params == nil {
		return nil
	}
	return &request.MetaDataFilterParams{
		FilterEnable:     params.FilterEnable,
		FilterLogicType:  params.FilterLogicType,
		MetaFilterParams: buildAssistantMetaFilterParams(params.MetaFilterParams),
	}
}

func buildAssistantMetaFilterParams(metaFilterList []*assistant_service.MetaFilterParams) []*request.MetaFilterParams {
	if metaFilterList == nil {
		return nil
	}
	var metaList []*request.MetaFilterParams
	for _, m := range metaFilterList {
		metaList = append(metaList, &request.MetaFilterParams{
			Condition: m.Condition,
			Key:       m.Key,
			Type:      m.Type,
			Value:     m.Value,
		})
	}
	return metaList
}

func transRequestFiles(files []*assistant_service.RequestFile) []response.RequestFile {
	if files == nil {
		return nil
	}
	var result []response.RequestFile
	for _, file := range files {
		result = append(result, response.RequestFile{
			FileName: file.FileName,
			FileSize: file.FileSize,
			FileUrl:  file.FileUrl,
		})
	}
	return result
}
