package knowledge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	knowledgebase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/generator"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/util"
	rag_service "github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"
	pkg_util "github.com/UnicomAI/wanwu/pkg/util"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	MetaValueTypeNumber   = "number"
	MetaValueTypeTime     = "time"
	MetaConditionEmpty    = "empty"
	MetaConditionNotEmpty = "not empty"
	MetaOperationAdd      = "add"
	MetaOperationUpdate   = "update"
	MetaOperationDelete   = "delete"
)

func (s *Service) SelectKnowledgeList(ctx context.Context, req *knowledgebase_service.KnowledgeSelectReq) (*knowledgebase_service.KnowledgeSelectListResp, error) {
	list, permissionMap, err := orm.SelectKnowledgeList(ctx, req.UserId, req.OrgId, req.Name, req.TagIdList)
	if err != nil {
		log.Errorf(fmt.Sprintf("获取知识库列表失败(%v)  参数(%v)", err, req))
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	var tagMap = make(map[string][]*orm.TagRelationDetail)
	var knowledgeIdList []string
	if len(list) > 0 {
		for _, k := range list {
			knowledgeIdList = append(knowledgeIdList, k.KnowledgeId)
		}
		relation := orm.SelectKnowledgeTagListWithRelation(ctx, req.UserId, req.OrgId, "", knowledgeIdList)
		tagMap = buildKnowledgeTagMap(relation)
	}
	return buildKnowledgeListResp(list, tagMap, permissionMap), nil
}

func (s *Service) SelectKnowledgeListByIdList(ctx context.Context, req *knowledgebase_service.BatchKnowledgeSelectReq) (*knowledgebase_service.KnowledgeSelectListResp, error) {
	list, permissionMap, err := orm.SelectKnowledgeByIdList(ctx, req.KnowledgeIdList, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("获取知识库列表失败(%v)  参数(%v)", err, req))
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	return buildKnowledgeListResp(list, nil, permissionMap), nil
}

func (s *Service) SelectKnowledgeDetailById(ctx context.Context, req *knowledgebase_service.KnowledgeDetailSelectReq) (*knowledgebase_service.KnowledgeInfo, error) {
	knowledgeInfo, err := orm.SelectKnowledgeById(ctx, req.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("获取知识库详情(%v)  参数(%v)", err, req))
		return nil, err
	}
	return buildKnowledgeInfo(knowledgeInfo), nil
}

func (s *Service) SelectKnowledgeDetailByName(ctx context.Context, req *knowledgebase_service.KnowledgeDetailSelectReq) (*knowledgebase_service.KnowledgeInfo, error) {
	knowledgeInfo, err := orm.SelectKnowledgeByName(ctx, req.KnowledgeName, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("根据名称获取知识库详情失败(%v)  参数(%v)", err, req))
		return nil, err
	}
	return buildKnowledgeInfo(knowledgeInfo), nil
}

func (s *Service) SelectKnowledgeDetailByIdList(ctx context.Context, req *knowledgebase_service.KnowledgeDetailSelectListReq) (*knowledgebase_service.KnowledgeDetailSelectListResp, error) {
	knowledgeInfoList, _, err := orm.SelectKnowledgeByIdList(ctx, req.KnowledgeIds, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("根据id列表获取知识库详情列表失败(%v)  参数(%v)", err, req))
		return nil, err
	}
	return buildKnowledgeInfoList(knowledgeInfoList), nil
}

func (s *Service) CreateKnowledge(ctx context.Context, req *knowledgebase_service.CreateKnowledgeReq) (*knowledgebase_service.CreateKnowledgeResp, error) {
	//1.重名校验
	err := orm.CheckSameKnowledgeName(ctx, req.UserId, req.OrgId, req.Name, "")
	if err != nil {
		return nil, err
	}
	//2.创建创建知识库
	knowledgeModel, err := buildKnowledgeBaseModel(req)
	if err != nil {
		log.Errorf("buildKnowledgeBaseModel error %s", err)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseCreateFailed)
	}
	err = orm.CreateKnowledge(ctx, knowledgeModel, req.EmbeddingModelInfo.ModelId)
	if err != nil {
		log.Errorf("CreateKnowledge error %v params %v", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseCreateFailed)
	}
	//3.返回结果
	return &knowledgebase_service.CreateKnowledgeResp{
		KnowledgeId: knowledgeModel.KnowledgeId,
	}, nil
}

func (s *Service) UpdateKnowledge(ctx context.Context, req *knowledgebase_service.UpdateKnowledgeReq) (*emptypb.Empty, error) {
	//1.查询知识库详情,这里前置做了前置权限校验，所以这里不需要再次校验
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeId, "", "")
	if err != nil {
		log.Errorf(fmt.Sprintf("没有操作该知识库的权限 参数(%v)", req))
		return nil, err
	}
	//2.重名校验
	err = orm.CheckSameKnowledgeName(ctx, req.UserId, req.OrgId, req.Name, knowledge.KnowledgeId)
	if err != nil {
		return nil, err
	}
	//3.更新知识库
	err = orm.UpdateKnowledge(ctx, req.Name, req.Description, knowledge)
	if err != nil {
		log.Errorf("知识库更新失败(%v)  参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseUpdateFailed)
	}
	return &emptypb.Empty{}, nil
}

// DeleteKnowledge 删除知识库
func (s *Service) DeleteKnowledge(ctx context.Context, req *knowledgebase_service.DeleteKnowledgeReq) (*emptypb.Empty, error) {
	//1.查询知识库详情
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("没有操作该知识库的权限 参数(%v)", req))
		return nil, err
	}
	//2.校验导入状态
	err = orm.SelectKnowledgeRunningImportTask(ctx, knowledge.KnowledgeId)
	if err != nil {
		return nil, err
	}
	//3.先删除知识库，异步删除资源数据
	err = orm.DeleteKnowledge(ctx, knowledge)
	if err != nil {
		log.Errorf("删除知识库失败 error %v params %v", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseDeleteFailed)
	}
	return &emptypb.Empty{}, nil
}

// KnowledgeHit 知识库命中测试
func (s *Service) KnowledgeHit(ctx context.Context, req *knowledgebase_service.KnowledgeHitReq) (*knowledgebase_service.KnowledgeHitResp, error) {
	// 1.获取知识库信息列表
	if len(req.KnowledgeList) == 0 || req.Question == "" || req.KnowledgeMatchParams == nil {
		return nil, util.ErrCode(errs.Code_KnowledgeInvalidArguments)
	}
	var knowledgeIdList []string
	for _, k := range req.KnowledgeList {
		knowledgeIdList = append(knowledgeIdList, k.KnowledgeId)
	}
	list, _, err := orm.SelectKnowledgeByIdList(ctx, knowledgeIdList, req.UserId, req.OrgId)
	if err != nil {
		return nil, err
	}
	knowledgeIDToName := make(map[string]string)
	for _, k := range list {
		if _, exists := knowledgeIDToName[k.KnowledgeId]; !exists {
			knowledgeIDToName[k.KnowledgeId] = k.Name
		}
	}
	// 2.RAG请求
	ragHitParams, err := buildRagHitParams(req, list, knowledgeIDToName)
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeBaseHitFailed)
	}
	hitResp, err := rag_service.RagKnowledgeHit(ctx, ragHitParams)
	if err != nil {
		log.Errorf("RagKnowledgeHit error %s", err)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseHitFailed)
	}
	return buildKnowledgeBaseHitResp(hitResp), nil
}

func (s *Service) GetKnowledgeMetaSelect(ctx context.Context, req *knowledgebase_service.SelectKnowledgeMetaReq) (*knowledgebase_service.SelectKnowledgeMetaResp, error) {
	metaList, err := orm.SelectMetaByKnowledgeId(ctx, "", "", req.KnowledgeId)
	if err != nil {
		log.Errorf("获取知识库元数据列表失败(%v)  参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeMetaFetchFailed)
	}
	return buildKnowledgeMetaSelectResp(metaList), nil
}

func (s *Service) GetKnowledgeMetaValueList(ctx context.Context, req *knowledgebase_service.KnowledgeMetaValueListReq) (*knowledgebase_service.KnowledgeMetaValueListResp, error) {
	metaList, err := orm.SelectMetaByDocIds(ctx, req.UserId, req.OrgId, req.DocIdList)
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeMetaFetchFailed)
	}
	return buildKnowledgeMetaValueListResp(metaList), nil
}

func (s *Service) UpdateKnowledgeMetaValue(ctx context.Context, req *knowledgebase_service.UpdateKnowledgeMetaValueReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, req.DocIdList, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法修改元数据 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocUpdateMetaStatusFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.查询元数据
	docMetaList, err := orm.SelectMetaByDocIds(ctx, req.UserId, req.OrgId, req.DocIdList)
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeMetaFetchFailed)
	}
	//5.构造文档元数据map
	docMetaMap := buildDocMetaMap(docMetaList)
	//6.构造元数据列表
	addList, updateList, deleteList := buildMetaList(req, docMetaMap, doc.KnowledgeId)
	//7.更新数据库并发送rag请求
	err = orm.BatchUpdateDocMetaValue(ctx, addList, updateList, deleteList, knowledge, docList, req.UserId, req.DocIdList)
	if err != nil {
		log.Errorf("更新文档元数据失败(%v)  参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeMetaUpdateFailed)
	}
	return nil, nil
}

func buildDocMetaMap(docMetaList []*model.KnowledgeDocMeta) map[string]map[string][]*model.KnowledgeDocMeta {
	docMetaMap := make(map[string]map[string][]*model.KnowledgeDocMeta)
	for _, v := range docMetaList {
		if _, exists := docMetaMap[v.DocId]; !exists {
			docMetaMap[v.DocId] = make(map[string][]*model.KnowledgeDocMeta)
		}
		if v.Value != "" {
			docMetaMap[v.DocId][v.Key] = append(docMetaMap[v.DocId][v.Key], v)
		}
	}
	return docMetaMap
}

func buildMetaList(req *knowledgebase_service.UpdateKnowledgeMetaValueReq, docMetaMap map[string]map[string][]*model.KnowledgeDocMeta, knowledgeId string) (addList, updateList []*model.KnowledgeDocMeta, deleteList []string) {
	// 处理请求数据
	reqMetaList := handleReqMetaList(req.MetaList)
	for _, meta := range reqMetaList {
		switch meta.Option {
		case MetaOperationAdd:
			handleAddMeta(req, meta, docMetaMap, knowledgeId, &addList, &updateList, &deleteList)
		case MetaOperationUpdate:
			handleUpdateMeta(req, meta, docMetaMap, knowledgeId, &addList, &updateList, &deleteList)
		case MetaOperationDelete:
			handleDeleteMeta(req, meta, docMetaMap, &deleteList)
		}
	}
	return
}

func handleReqMetaList(metaList []*knowledgebase_service.MetaValueOperation) (reqMetaList []*knowledgebase_service.MetaValueOperation) {
	if len(metaList) > 100 {
		log.Infof("metaList size exceeds 100")
		metaList = metaList[:100]
	}
	keyMap := make(map[string]*knowledgebase_service.MetaValueOperation)
	for _, meta := range metaList {
		if _, exists := keyMap[meta.MetaInfo.Key]; !exists {
			keyMap[meta.MetaInfo.Key] = meta
		} else {
			// 同一key优先级：删除 > 更新 > 新增
			if meta.Option == MetaOperationDelete {
				keyMap[meta.MetaInfo.Key] = meta
			} else if meta.Option == MetaOperationUpdate {
				if keyMap[meta.MetaInfo.Key].Option == MetaOperationAdd {
					keyMap[meta.MetaInfo.Key] = meta
				}
			}
		}
	}
	for _, meta := range keyMap {
		reqMetaList = append(reqMetaList, meta)
	}
	return
}

func handleAddMeta(req *knowledgebase_service.UpdateKnowledgeMetaValueReq, meta *knowledgebase_service.MetaValueOperation, docMetaMap map[string]map[string][]*model.KnowledgeDocMeta, knowledgeId string, addList, updateList *[]*model.KnowledgeDocMeta, deleteList *[]string) {
	for _, docId := range req.DocIdList {
		existMetaList := docMetaMap[docId][meta.MetaInfo.Key]
		if len(existMetaList) > 0 {
			existMetaList[0].Value = meta.MetaInfo.Value
			*updateList = append(*updateList, existMetaList[0])
			for i := 1; i < len(existMetaList); i++ {
				*deleteList = append(*deleteList, existMetaList[i].MetaId)
			}
		} else {
			*addList = append(*addList, &model.KnowledgeDocMeta{
				MetaId:      generator.GetGenerator().NewID(),
				DocId:       docId,
				KnowledgeId: knowledgeId,
				UserId:      req.UserId,
				OrgId:       req.OrgId,
				Key:         meta.MetaInfo.Key,
				Value:       meta.MetaInfo.Value,
				ValueType:   meta.MetaInfo.Type,
			})
		}
	}
}

func handleUpdateMeta(req *knowledgebase_service.UpdateKnowledgeMetaValueReq, meta *knowledgebase_service.MetaValueOperation, docMetaMap map[string]map[string][]*model.KnowledgeDocMeta, knowledgeId string, addList, updateList *[]*model.KnowledgeDocMeta, deleteList *[]string) {
	for _, docId := range req.DocIdList {
		existMetaList := docMetaMap[docId][meta.MetaInfo.Key]
		if len(existMetaList) > 0 {
			existMetaList[0].Value = meta.MetaInfo.Value
			*updateList = append(*updateList, existMetaList[0])
			for i := 1; i < len(existMetaList); i++ {
				*deleteList = append(*deleteList, existMetaList[i].MetaId)
			}
		} else if req.ApplyToSelected {
			*addList = append(*addList, &model.KnowledgeDocMeta{
				MetaId:      generator.GetGenerator().NewID(),
				DocId:       docId,
				KnowledgeId: knowledgeId,
				UserId:      req.UserId,
				OrgId:       req.OrgId,
				Key:         meta.MetaInfo.Key,
				Value:       meta.MetaInfo.Value,
				ValueType:   meta.MetaInfo.Type,
			})
		}
	}
}

func handleDeleteMeta(req *knowledgebase_service.UpdateKnowledgeMetaValueReq, meta *knowledgebase_service.MetaValueOperation, docMetaMap map[string]map[string][]*model.KnowledgeDocMeta, deleteList *[]string) {
	for _, docId := range req.DocIdList {
		existMetaList := docMetaMap[docId][meta.MetaInfo.Key]
		for _, v := range existMetaList {
			*deleteList = append(*deleteList, v.MetaId)
		}
	}
}

func buildRagHitParams(req *knowledgebase_service.KnowledgeHitReq, list []*model.KnowledgeBase, knowledgeIDToName map[string]string) (*rag_service.KnowledgeHitParams, error) {
	matchParams := req.KnowledgeMatchParams
	priorityMatch := matchParams.PriorityMatch
	filterEnable, metaParams, err := buildRagHitMetaParams(req, knowledgeIDToName)
	if err != nil {
		return nil, err
	}
	idList, nameList := buildKnowledgeList(list)
	ret := &rag_service.KnowledgeHitParams{
		UserId:               req.UserId,
		Question:             req.Question,
		KnowledgeIdList:      idList,
		KnowledgeBase:        nameList,
		TopK:                 matchParams.TopK,
		Threshold:            float64(matchParams.Score),
		RerankModelId:        buildRerankId(priorityMatch, matchParams.RerankModelId),
		RetrieveMethod:       buildRetrieveMethod(matchParams.MatchType),
		RerankMod:            buildRerankMod(priorityMatch),
		Weight:               buildWeight(priorityMatch, matchParams.SemanticsPriority, matchParams.KeywordPriority),
		TermWeight:           buildTermWeight(matchParams.TermWeight, matchParams.TermWeightEnable),
		MetaFilter:           filterEnable,
		MetaFilterConditions: metaParams,
	}
	return ret, nil
}

func buildRagHitMetaParams(req *knowledgebase_service.KnowledgeHitReq, knowledgeIDToName map[string]string) (bool, []*rag_service.MetadataFilterItem, error) {
	filterEnable := false // 标记是否有启用的元数据过滤
	var metaFilterConditions []*rag_service.MetadataFilterItem
	for _, k := range req.KnowledgeList {
		// 检查元数据过滤参数是否有效
		filterParams := k.MetaDataFilterParams
		if !isValidFilterParams(k.MetaDataFilterParams) {
			continue
		}
		// 校验合法值
		if k.MetaDataFilterParams.FilterLogicType == "" {
			return false, nil, errors.New("FilterLogicType is empty")
		}
		// 标记元数据过滤生效
		filterEnable = true
		// 构建元数据过滤条件
		metaItems, err := buildRagHitMetaItems(k.KnowledgeId, filterParams.MetaFilterParams)
		if err != nil {
			return false, nil, err
		}
		// 添加过滤项到结果
		metaFilterConditions = append(metaFilterConditions, &rag_service.MetadataFilterItem{
			FilterKnowledgeName: knowledgeIDToName[k.KnowledgeId],
			LogicalOperator:     filterParams.FilterLogicType,
			Conditions:          metaItems,
		})
	}
	return filterEnable, metaFilterConditions, nil
}

// 构建元数据项列表
func buildRagHitMetaItems(knowledgeID string, params []*knowledgebase_service.MetaFilterParams) ([]*rag_service.MetaItem, error) {
	var metaItems []*rag_service.MetaItem
	for _, param := range params {
		// 基础参数校验
		if err := validateMetaFilterParam(knowledgeID, param); err != nil {
			return nil, err
		}
		// 转换参数值
		ragValue, err := convertValue(param.Value, param.Type)
		if err != nil {
			log.Errorf("kbId: %s, convert value failed: %v", knowledgeID, err)
			return nil, fmt.Errorf("convert value for key %s: %s", param.Key, err.Error())
		}
		metaItems = append(metaItems, &rag_service.MetaItem{
			MetaName:           param.Key,
			MetaType:           param.Type,
			ComparisonOperator: param.Condition,
			Value:              ragValue,
		})
	}
	return metaItems, nil
}

// 校验元数据过滤参数
func validateMetaFilterParam(knowledgeID string, param *knowledgebase_service.MetaFilterParams) error {
	// 检查关键参数是否为空
	if param.Key == "" || param.Type == "" || param.Condition == "" {
		errMsg := "key/type/condition cannot be empty"
		log.Errorf("kbId: %s, %s", knowledgeID, errMsg)
		return errors.New(errMsg)
	}

	// 检查空条件与值的匹配性
	if param.Condition == MetaConditionEmpty || param.Condition == MetaConditionNotEmpty {
		if param.Value != "" {
			errMsg := "condition is empty/non-empty, value should be empty"
			log.Errorf("kbId: %s, %s", knowledgeID, errMsg)
			return errors.New(errMsg)
		}
	} else {
		if param.Value == "" {
			errMsg := "value is empty"
			log.Errorf("kbId: %s, %s", knowledgeID, errMsg)
			return errors.New(errMsg)
		}
	}

	return nil
}

func isValidFilterParams(params *knowledgebase_service.MetaDataFilterParams) bool {
	return params != nil &&
		params.FilterEnable &&
		params.MetaFilterParams != nil &&
		len(params.MetaFilterParams) > 0
}

func convertValue(value, valueType string) (interface{}, error) {
	if len(value) == 0 {
		return nil, nil
	}
	// 根据类型转换value
	if valueType == MetaValueTypeNumber || valueType == MetaValueTypeTime {
		ragValue, err := pkg_util.I64(value)
		if err != nil {
			log.Errorf("convertMetaValue fail %v", err)
			return nil, err
		}
		return ragValue, nil
	}
	return value, nil
}

func buildKnowledgeMetaSelectResp(metaList []*model.KnowledgeDocMeta) *knowledgebase_service.SelectKnowledgeMetaResp {
	if len(metaList) == 0 {
		return &knowledgebase_service.SelectKnowledgeMetaResp{}
	}
	var retMetaList []*knowledgebase_service.KnowledgeMetaData
	newMetaList := checkRepeatedMetaKey(metaList)
	for _, meta := range newMetaList {
		if meta.Key != "" {
			retMetaList = append(retMetaList, &knowledgebase_service.KnowledgeMetaData{
				MetaId: meta.MetaId,
				Key:    meta.Key,
				Type:   meta.ValueType,
			})
		}
	}
	return &knowledgebase_service.SelectKnowledgeMetaResp{
		MetaList: retMetaList,
	}
}

// buildKnowledgeListResp 构造知识库列表返回结果
func buildKnowledgeListResp(knowledgeList []*model.KnowledgeBase, knowledgeTagMap map[string][]*orm.TagRelationDetail, permissionMap map[string]int) *knowledgebase_service.KnowledgeSelectListResp {
	if len(knowledgeList) == 0 {
		return &knowledgebase_service.KnowledgeSelectListResp{}
	}
	var retList []*knowledgebase_service.KnowledgeInfo
	for _, knowledge := range knowledgeList {
		knowledgeInfo := buildKnowledgeInfo(knowledge)
		knowledgeInfo.KnowledgeTagInfoList = buildKnowledgeTagList(knowledge.KnowledgeId, knowledgeTagMap)
		knowledgeInfo.PermissionType = buildKnowledgePermission(knowledge.KnowledgeId, permissionMap)
		retList = append(retList, knowledgeInfo)
	}
	return &knowledgebase_service.KnowledgeSelectListResp{
		KnowledgeList: retList,
	}
}

func buildKnowledgeTagMap(tagRelation *orm.TagRelation) map[string][]*orm.TagRelationDetail {
	if tagRelation.RelationErr != nil || tagRelation.TagErr != nil {
		return make(map[string][]*orm.TagRelationDetail)
	}
	var knowledgeTagMap = make(map[string][]*orm.TagRelationDetail)
	for _, relation := range tagRelation.RelationList {
		details := knowledgeTagMap[relation.KnowledgeId]
		if details == nil {
			details = make([]*orm.TagRelationDetail, 0)
		}
		for _, tag := range tagRelation.TagList {
			if tag.TagId == relation.TagId {
				details = append(details, &orm.TagRelationDetail{
					TagId:   tag.TagId,
					TagName: tag.Name,
				})
			}
		}
		knowledgeTagMap[relation.KnowledgeId] = details
	}
	return knowledgeTagMap
}

func buildKnowledgeTagList(knowledgeId string, knowledgeTagMap map[string][]*orm.TagRelationDetail) []*knowledgebase_service.KnowledgeTagInfo {
	if len(knowledgeTagMap) == 0 {
		return []*knowledgebase_service.KnowledgeTagInfo{}
	}
	tagList := knowledgeTagMap[knowledgeId]
	if len(tagList) == 0 {
		return []*knowledgebase_service.KnowledgeTagInfo{}
	}
	var retList []*knowledgebase_service.KnowledgeTagInfo
	for _, tag := range tagList {
		retList = append(retList, &knowledgebase_service.KnowledgeTagInfo{
			TagId:   tag.TagId,
			TagName: tag.TagName,
		})
	}
	return retList
}

func buildKnowledgePermission(knowledgeId string, permissionMap map[string]int) int32 {
	return int32(permissionMap[knowledgeId])
}

func checkRepeatedMetaKey(metaList []*model.KnowledgeDocMeta) []*model.KnowledgeDocMeta {
	if len(metaList) == 0 {
		return []*model.KnowledgeDocMeta{}
	}
	return lo.UniqBy(metaList, func(item *model.KnowledgeDocMeta) string {
		return item.Key
	})
}

// buildKnowledgeInfo 构造知识库信息
func buildKnowledgeInfo(knowledge *model.KnowledgeBase) *knowledgebase_service.KnowledgeInfo {
	embeddingModelInfo := &knowledgebase_service.EmbeddingModelInfo{}
	_ = json.Unmarshal([]byte(knowledge.EmbeddingModel), embeddingModelInfo)
	return &knowledgebase_service.KnowledgeInfo{
		KnowledgeId:        knowledge.KnowledgeId,
		Name:               knowledge.Name,
		Description:        knowledge.Description,
		DocCount:           int32(knowledge.DocCount),
		ShareCount:         int32(knowledge.ShareCount),
		EmbeddingModelInfo: embeddingModelInfo,
		CreatedAt:          pkg_util.Time2Str(knowledge.CreatedAt),
		CreateOrgId:        knowledge.OrgId,
		CreateUserId:       knowledge.UserId,
	}
}

// buildKnowledgeInfoList 构造知识库信息列表
func buildKnowledgeInfoList(knowledgeList []*model.KnowledgeBase) *knowledgebase_service.KnowledgeDetailSelectListResp {
	var retList []*knowledgebase_service.KnowledgeInfo
	for _, v := range knowledgeList {
		info := buildKnowledgeInfo(v)
		retList = append(retList, info)
	}
	return &knowledgebase_service.KnowledgeDetailSelectListResp{
		List:  retList,
		Total: int32(len(retList)),
	}
}

// buildKnowledgeBaseModel 构造知识库模型
func buildKnowledgeBaseModel(req *knowledgebase_service.CreateKnowledgeReq) (*model.KnowledgeBase, error) {
	embeddingModelInfo, err := json.Marshal(req.EmbeddingModelInfo)
	if err != nil {
		return nil, err
	}
	return &model.KnowledgeBase{
		KnowledgeId:    generator.GetGenerator().NewID(),
		Name:           req.Name,
		Description:    req.Description,
		OrgId:          req.OrgId,
		UserId:         req.UserId,
		EmbeddingModel: string(embeddingModelInfo),
		CreatedAt:      time.Now().UnixMilli(),
		UpdatedAt:      time.Now().UnixMilli(),
	}, nil
}

// buildKnowledgeList 构造知识库名称
func buildKnowledgeList(knowledgeList []*model.KnowledgeBase) (knowledgeIdList []string, knowledgeNameList []string) {
	if len(knowledgeList) == 0 {
		return make([]string, 0), make([]string, 0)
	}
	for _, knowledge := range knowledgeList {
		knowledgeNameList = append(knowledgeNameList, knowledge.Name)
		knowledgeIdList = append(knowledgeIdList, knowledge.KnowledgeId)
	}
	return
}

// buildKnowledgeBaseHitResp 构造知识库命中返回
func buildKnowledgeBaseHitResp(ragKnowledgeHitResp *rag_service.RagKnowledgeHitResp) *knowledgebase_service.KnowledgeHitResp {
	knowledgeHitData := ragKnowledgeHitResp.Data
	var searchList = make([]*knowledgebase_service.KnowledgeSearchInfo, 0)
	list := knowledgeHitData.SearchList
	if len(list) > 0 {
		for _, search := range list {
			childContentList := make([]*knowledgebase_service.ChildContent, 0)
			for _, child := range search.ChildContentList {
				childContentList = append(childContentList, &knowledgebase_service.ChildContent{
					ChildSnippet: child.ChildSnippet,
					Score:        float32(child.Score),
				})
			}
			childScore := make([]float32, 0)
			for _, score := range search.ChildScore {
				childScore = append(childScore, float32(score))
			}
			searchList = append(searchList, &knowledgebase_service.KnowledgeSearchInfo{
				Title:            search.Title,
				Snippet:          search.Snippet,
				KnowledgeName:    search.KbName,
				ChildContentList: childContentList,
				ChildScore:       childScore,
			})
		}
	}
	return &knowledgebase_service.KnowledgeHitResp{
		Prompt:     knowledgeHitData.Prompt,
		Score:      knowledgeHitData.Score,
		SearchList: searchList,
	}
}

// buildRerankId 构造重排序模型id
func buildRerankId(priorityType int32, rerankId string) string {
	if priorityType == 1 {
		return ""
	}
	return rerankId
}

// buildRetrieveMethod 构造检索方式
func buildRetrieveMethod(matchType string) string {
	switch matchType {
	case "vector":
		return "semantic_search"
	case "text":
		return "full_text_search"
	case "mix":
		return "hybrid_search"
	}
	return ""
}

// buildRerankMod 构造重排序模式
func buildRerankMod(priorityType int32) string {
	if priorityType == 1 {
		return "weighted_score"
	}
	return "rerank_model"
}

// buildWeight 构造权重信息
func buildWeight(priorityType int32, semanticsPriority float32, keywordPriority float32) *rag_service.WeightParams {
	if priorityType != 1 {
		return nil
	}
	return &rag_service.WeightParams{
		VectorWeight: semanticsPriority,
		TextWeight:   keywordPriority,
	}
}

// buildTermWeight 构造关键词系数信息
func buildTermWeight(termWeight float32, termWeightEnable bool) float32 {
	if termWeightEnable {
		return termWeight
	}
	return 0.0
}

func buildKnowledgeMetaValueListResp(metaList []*model.KnowledgeDocMeta) *knowledgebase_service.KnowledgeMetaValueListResp {
	retMap := make(map[string]*knowledgebase_service.KnowledgeMetaValues)
	var retList []*knowledgebase_service.KnowledgeMetaValues
	for _, meta := range metaList {
		if meta.Value == "" || meta.Key == "" || meta.ValueType == "" {
			continue
		}
		if _, exists := retMap[meta.Key]; !exists {
			retMap[meta.Key] = &knowledgebase_service.KnowledgeMetaValues{
				MetaId:    meta.MetaId,
				Key:       meta.Key,
				Type:      meta.ValueType,
				ValueList: []string{meta.Value},
			}
		} else {
			retMap[meta.Key].ValueList = append(retMap[meta.Key].ValueList, meta.Value)
		}
	}
	for _, retMeta := range retMap {
		retMeta.ValueList = lo.Uniq(retMeta.ValueList)
		retList = append(retList, retMeta)
	}
	return &knowledgebase_service.KnowledgeMetaValueListResp{
		MetaList: retList,
	}
}
