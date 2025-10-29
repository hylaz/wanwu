package knowledge_doc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	knowledgebase_doc_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-doc-service"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/generator"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/util"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"
	util2 "github.com/UnicomAI/wanwu/pkg/util"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"strings"
	"time"
)

const (
	fiveMinutes               int64 = 5 * 60 * 1000
	noSplitter                      = "未设置"
	segmentImportingMessage         = "分段内容正在上传解析中"
	segmentCompleteFormat           = "分段内容解析完成，成功%d，失败%d"
	segmentPartCompleteFormat       = "分段内容解析完成，成功%d"
	segmentCompleteFail             = "分段内容解析失败"
	DocImportIng                    = 1
	DocImportFinish                 = 2
	DocImportError                  = 3
	MetaOptionDelete                = "delete"
	MetaOptionAdd                   = "add"
	MetaOptionUpdate                = "update"
	MetaStatusFailed                = "failed"
	MetaStatusPartial               = "partial"
)

func (s *Service) GetDocList(ctx context.Context, req *knowledgebase_doc_service.GetDocListReq) (*knowledgebase_doc_service.GetDocListResp, error) {
	//入口层已经校验过用户权限，此处无需校验
	list, total, err := orm.GetDocList(ctx, "", "", req.KnowledgeId,
		req.DocName, req.DocTag, util.BuildDocReqStatusList(int(req.Status)), req.PageSize, req.PageNum)
	if err != nil {
		log.Errorf("获取知识库列表失败(%v)  参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	//查询配置信息
	var importTaskList []*model.KnowledgeImportTask
	if len(list) > 0 {
		importTaskList, err = orm.SelectKnowledgeImportTaskByIdList(ctx, buildImportTaskIdList(list))
		if err != nil {
			log.Errorf("获取知识库列表失败(%v)  参数(%v)", err, req)
		}
	}

	return buildDocListResp(list, importTaskList, total, req.PageSize, req.PageNum), nil
}

func (s *Service) GetDocDetail(ctx context.Context, req *knowledgebase_doc_service.GetDocDetailReq) (*knowledgebase_doc_service.DocInfo, error) {
	doc, err := orm.GetDocDetail(ctx, req.UserId, req.OrgId, req.DocId)
	if err != nil {
		log.Errorf("获取知识库列表失败(%v)  参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	return buildDocInfo(doc, make(map[string]*model.SegmentConfig)), nil
}

func (s *Service) ImportDoc(ctx context.Context, req *knowledgebase_doc_service.ImportDocReq) (*emptypb.Empty, error) {
	task, err := buildImportTask(req)
	if err != nil {
		return nil, err
	}
	//创建导入任务
	err = orm.CreateKnowledgeImportTask(ctx, task)
	if err != nil {
		log.Errorf("import doc fail %v", err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocImportFail)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateDocStatus(ctx context.Context, req *knowledgebase_doc_service.UpdateDocStatusReq) (*emptypb.Empty, error) {
	err := orm.UpdateDocStatusDocId(ctx, req.DocId, int(req.Status), buildMetaParamsList(removeDuplicateMeta(req.MetaDataList)))
	if err != nil {
		log.Errorf("docId: %v update doc fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocUpdateStatusFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateDocMetaData(ctx context.Context, req *knowledgebase_doc_service.UpdateDocMetaDataReq) (*emptypb.Empty, error) {
	if len(req.MetaDataList) == 0 {
		return &emptypb.Empty{}, nil
	}
	// 更新文档元数据
	if len(req.DocId) > 0 {
		return updateDocMetaData(ctx, req)
	}
	// 更新知识库元数据
	if len(req.KnowledgeId) > 0 {
		return updateKnowledgeMetaData(ctx, req)
	}
	return nil, util.ErrCode(errs.Code_KnowledgeDocUpdateMetaStatusFailed)
}

func buildMetaDocMap(metaList []*model.KnowledgeDocMeta) map[string][]*model.KnowledgeDocMeta {
	dataMap := make(map[string][]*model.KnowledgeDocMeta)
	if len(metaList) == 0 {
		return dataMap
	}
	for _, meta := range metaList {
		metas := dataMap[meta.Key]
		if len(metas) == 0 {
			metas = make([]*model.KnowledgeDocMeta, 0)
		}
		metas = append(metas, meta)
		dataMap[meta.Key] = metas
	}
	return dataMap
}

// updateKnowledgeMetaData 更新知识库元数据
func updateKnowledgeMetaData(ctx context.Context, req *knowledgebase_doc_service.UpdateDocMetaDataReq) (*emptypb.Empty, error) {
	// 1.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	// 2.查询知识库元数据
	metaList, err := orm.SelectMetaByKnowledgeId(ctx, req.UserId, req.OrgId, req.KnowledgeId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 错误(%v) 参数(%v)", err, req)
		return nil, err
	}
	// 3.构造各种操作列表
	deleteList, updateList, addList := buildOptionList(metaList, req)
	// 4.校验updateList和addList
	err = checkUpdateAndAddMetaList(addList, updateList, metaList)
	if err != nil {
		log.Errorf("更新元数据失败 错误(%v) 参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeMetaDuplicateKey)
	}
	updateStatus := MetaStatusFailed
	// 5.执行批量删除
	if len(deleteList) > 0 {
		err = orm.BatchDeleteMeta(ctx, deleteList, req.KnowledgeId, &service.RagBatchDeleteMetaParams{
			UserId:        knowledge.UserId,
			KnowledgeBase: knowledge.Name,
			KnowledgeId:   req.KnowledgeId,
			Keys:          deleteList,
		})
		if err != nil {
			log.Errorf("删除元数据失败 错误(%v) 删除参数(%v)", err, req)
			return nil, util.ErrCode(errs.Code_KnowledgeMetaDeleteFailed)
		}
		updateStatus = MetaStatusPartial
	}
	// 6.执行批量更新
	if len(updateList) > 0 {
		err = orm.BatchUpdateMetaKey(ctx, updateList, req.KnowledgeId, &service.RagBatchUpdateMetaKeyParams{
			UserId:        req.UserId,
			KnowledgeBase: knowledge.Name,
			KnowledgeId:   req.KnowledgeId,
			Mappings:      updateList,
		})
		if err != nil {
			log.Errorf("更新元数据失败 错误(%v) 更新参数(%v)", err, req)
			if updateStatus == MetaStatusPartial {
				return nil, util.ErrCode(errs.Code_KnowledgeMetaUpdatePartialSuccess)
			}
			return nil, util.ErrCode(errs.Code_KnowledgeMetaUpdateFailed)
		}
		if updateStatus == MetaStatusFailed {
			updateStatus = MetaStatusPartial
		}
	}
	// 7.执行批量新增
	if len(addList) > 0 {
		err = orm.BatchAddMeta(ctx, addList)
		if err != nil {
			log.Errorf("新增元数据失败 错误(%v) 更新参数(%v)", err, req)
			if updateStatus == MetaStatusPartial {
				return nil, util.ErrCode(errs.Code_KnowledgeMetaUpdatePartialSuccess)
			}
			return nil, util.ErrCode(errs.Code_KnowledgeMetaCreateFailed)
		}
	}

	return &emptypb.Empty{}, nil
}

func buildImportTaskIdList(docList []*model.KnowledgeDoc) []string {
	return lo.Map(docList, func(item *model.KnowledgeDoc, index int) string {
		return item.ImportTaskId
	})
}

// updateDocMetaData 更新文档元数据
func updateDocMetaData(ctx context.Context, req *knowledgebase_doc_service.UpdateDocMetaDataReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法增加元数据 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocUpdateMetaStatusFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.查询元数据
	metaDocList, err := orm.SelectMetaByKnowledgeId(ctx, req.UserId, req.OrgId, knowledge.KnowledgeId)
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeDocUpdateMetaStatusFailed)
	}
	docMetaMap := buildMetaDocMap(metaDocList)
	//5.构造元数据操作列表
	metaDataList := removeDuplicateMeta(req.MetaDataList)
	var fileName = service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	addList, updateList, deleteList := buildDocMetaModelList(metaDataList, doc.OrgId, doc.UserId, req.KnowledgeId, req.DocId)
	if err1 := checkMetaKeyType(addList, updateList, docMetaMap); err1 != nil {
		return nil, err1
	}
	//6.构造RAG请求参数
	params, err := buildMetaRagParams(metaDataList)
	if err != nil {
		log.Errorf("docId %v update buildMetaRagParams fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocUpdateMetaFailed)
	}
	//7.更新数据库并发送RAG请求
	err = orm.UpdateDocStatusDocMeta(ctx, req.DocId, addList, updateList, deleteList,
		&service.RagDocMetaParams{
			FileName:      fileName,
			KnowledgeBase: knowledge.Name,
			UserId:        req.UserId,
			MetaList:      params,
		})
	if err != nil {
		log.Errorf("docId %v update doc tag fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocUpdateMetaFailed)
	}
	return &emptypb.Empty{}, nil
}

func buildKnowledgeMetaMap(metaList []*model.KnowledgeDocMeta) map[string]string {
	metaMap := make(map[string]string)
	for _, meta := range metaList {
		metaMap[meta.MetaId] = meta.Key
	}
	return metaMap
}

func buildUpdateMetaMap(metaList []*knowledgebase_doc_service.MetaData, metaMap map[string]string) []*service.RagMetaMapKeys {
	metaMapKeys := make([]*service.RagMetaMapKeys, 0)
	for _, reqMeta := range metaList {
		if reqMeta.Option == MetaOptionUpdate {
			if dbKey, exists := metaMap[reqMeta.MetaId]; !exists {
				log.Errorf("metaId %s doesn't exist", reqMeta.MetaId)
				continue
			} else if dbKey == "" {
				log.Errorf("metaId %s dbKey is empty", reqMeta.MetaId)
				continue
			} else if dbKey != reqMeta.Key {
				metaMapKeys = append(metaMapKeys, &service.RagMetaMapKeys{
					NewKey: reqMeta.Key,
					OldKey: dbKey,
				})
			}
		}
	}
	return metaMapKeys
}

func buildAddMetaList(req *knowledgebase_doc_service.UpdateDocMetaDataReq) []*model.KnowledgeDocMeta {
	addList := make([]*model.KnowledgeDocMeta, 0)
	for _, reqMeta := range req.MetaDataList {
		if reqMeta.Option == MetaOptionAdd {
			addList = append(addList, &model.KnowledgeDocMeta{
				KnowledgeId: req.KnowledgeId,
				MetaId:      generator.GetGenerator().NewID(),
				Key:         reqMeta.Key,
				ValueType:   reqMeta.ValueType,
				Rule:        "",
				OrgId:       req.OrgId,
				UserId:      req.UserId,
				CreatedAt:   time.Now().UnixMilli(),
				UpdatedAt:   time.Now().UnixMilli(),
			})
		}
	}
	return addList
}

func checkMetaKeyType(addList []*model.KnowledgeDocMeta, updateList []*model.KnowledgeDocMeta, docMetaMap map[string][]*model.KnowledgeDocMeta) error {
	if len(addList) > 0 {
		for _, meta := range addList {
			data := docMetaMap[meta.Key]
			if len(data) > 0 {
				for _, datum := range data {
					if datum.ValueType != meta.ValueType {
						log.Errorf("meta key %s datum metaId %s type %s meta type %s error", meta.Key, datum.MetaId, datum.ValueType, meta.ValueType)
						return util.ErrCode(errs.Code_KnowledgeDocUpdateMetaSameKeyFailed)
					}
				}
			}
		}
	}
	if len(updateList) > 0 {
		for _, meta := range updateList {
			data := docMetaMap[meta.Key]
			if len(data) > 0 {
				for _, datum := range data {
					if datum.MetaId != meta.MetaId && datum.ValueType != meta.ValueType {
						log.Errorf("meta key %s datum type %s meta type %s error", meta.Key, datum.ValueType, meta.ValueType)
						return util.ErrCode(errs.Code_KnowledgeDocUpdateMetaSameKeyFailed)
					}
				}
			}
		}
	}
	return nil
}

func (s *Service) InitDocStatus(ctx context.Context, req *knowledgebase_doc_service.InitDocStatusReq) (*emptypb.Empty, error) {
	err := orm.InitDocStatus(ctx, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("init doc fail %v", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeGeneral)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteDoc(ctx context.Context, req *knowledgebase_doc_service.DeleteDocReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, req.Ids, "", "")
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//2.校验导入状态
	docIdList, resultDocList, err := checkDocStatus(docList)
	if err != nil {
		log.Errorf("删除知识库文件失败 error %v params %v", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocDeleteDuringParse)
	}
	if len(docIdList) == 0 {
		return &emptypb.Empty{}, nil
	}
	//3.删除文档
	err = orm.DeleteDocByIdList(ctx, docIdList, resultDocList)
	if err != nil {
		log.Errorf("删除知识库文件失败 error %v params %v", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocDeleteFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetDocCategoryUploadTip(ctx context.Context, req *knowledgebase_doc_service.DocImportTipReq) (*knowledgebase_doc_service.DocImportTipResp, error) {
	//1.查询知识库详情,前置参数校验
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeId, "", "")
	if err != nil {
		return nil, err
	}
	//2.查询第一个异步任务信息
	taskList, err := orm.SelectKnowledgeLatestImportTask(ctx, req.KnowledgeId)
	if err != nil {
		return nil, err
	}
	if len(taskList) == 0 {
		return &knowledgebase_doc_service.DocImportTipResp{
			KnowledgeId:   req.KnowledgeId,
			KnowledgeName: knowledge.Name,
			UploadStatus:  DocImportFinish,
		}, nil
	}
	if len(taskList) > 0 {
		task := taskList[0]
		if task.Status == model.KnowledgeImportError {
			return &knowledgebase_doc_service.DocImportTipResp{
				KnowledgeId:   req.KnowledgeId,
				KnowledgeName: knowledge.Name,
				Message:       "\n" + task.ErrorMsg,
				UploadStatus:  DocImportError,
			}, nil
		} else if task.Status == model.KnowledgeImportFinish {
			return &knowledgebase_doc_service.DocImportTipResp{
				KnowledgeId:   req.KnowledgeId,
				KnowledgeName: knowledge.Name,
				UploadStatus:  DocImportFinish,
			}, nil
		}
	}
	return &knowledgebase_doc_service.DocImportTipResp{
		KnowledgeId:   req.KnowledgeId,
		KnowledgeName: knowledge.Name,
		Message:       "",
		UploadStatus:  DocImportIng,
	}, nil
}

func (s *Service) GetDocSegmentList(ctx context.Context, req *knowledgebase_doc_service.DocSegmentListReq) (*knowledgebase_doc_service.DocSegmentListResp, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, "", "")
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	docInfo := docList[0]
	//2.查询知识库详情
	knowledge, err := orm.SelectKnowledgeById(ctx, docInfo.KnowledgeId, "", "")
	if err != nil {
		log.Errorf("查询知识库详情失败 参数(%v)", req)
		return nil, err
	}
	//3.查询知识库导入详情
	importTask, err := orm.SelectKnowledgeImportTaskById(ctx, docInfo.ImportTaskId)
	if err != nil {
		log.Errorf("查询知识库导入详情失败 参数(%v)", req)
		return nil, err
	}
	//4.查询最新导入详情
	segmentImportTask, err := orm.SelectSegmentLatestImportTaskByDocID(ctx, docInfo.DocId)
	//此处失败不影响详情展示
	if err != nil {
		log.Errorf("查询知识库导入详情失败 参数(%v)", req)
	}
	//4.查询分片信息
	segmentListResp, err := service.RagGetDocSegmentList(ctx, &service.RagGetDocSegmentParams{
		UserId:            knowledge.UserId,
		KnowledgeBaseName: knowledge.Name,
		FileName:          service.RebuildFileName(docInfo.DocId, docInfo.FileType, docInfo.Name),
		PageSize:          req.PageSize,
		SearchAfter:       req.PageSize * (req.PageNo - 1),
	})
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeDocSplitFailed)
	}
	//5.查询文档元数据,忽略错误
	metaDataList, _ := orm.SelectDocMetaList(ctx, req.UserId, req.OrgId, req.DocId)
	return buildSegmentListResp(importTask, docInfo, segmentListResp, req.PageNo, req.PageSize, metaDataList, segmentImportTask)
}

func (s *Service) GetDocChildSegmentList(ctx context.Context, req *knowledgebase_doc_service.GetDocChildSegmentListReq) (*knowledgebase_doc_service.GetDocChildSegmentListResp, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, "", "")
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	docInfo := docList[0]
	//2.查询知识库详情
	knowledge, err := orm.SelectKnowledgeById(ctx, docInfo.KnowledgeId, "", "")
	if err != nil {
		log.Errorf("查询知识库详情失败 参数(%v)", req)
		return nil, err
	}
	//3.查询分片信息
	segmentListResp, err := service.RagGetDocChildSegmentList(ctx, &service.RagGetDocChildSegmentParams{
		UserId:            knowledge.UserId,
		KnowledgeBaseName: knowledge.Name,
		KnowledgeId:       knowledge.KnowledgeId,
		FileName:          service.RebuildFileName(docInfo.DocId, docInfo.FileType, docInfo.Name),
		ChunkId:           req.ContentId,
	})
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeDocSplitFailed)
	}
	return buildChildSegmentListResp(segmentListResp)
}

func (s *Service) AnalysisDocUrl(ctx context.Context, req *knowledgebase_doc_service.AnalysisUrlDocReq) (*knowledgebase_doc_service.AnalysisUrlDocResp, error) {
	analysisResult, err := service.BatchRagDocUrlAnalysis(ctx, req.UrlList)
	if err != nil {
		return nil, err
	}
	var retUrlList []*knowledgebase_doc_service.UrlInfo
	for _, result := range analysisResult {
		retUrlList = append(retUrlList, &knowledgebase_doc_service.UrlInfo{
			Url:      result.Url,
			FileName: util.UrlNameFilter(result.FileName),
			FileSize: result.FileSize,
		})
	}
	return &knowledgebase_doc_service.AnalysisUrlDocResp{UrlList: retUrlList}, nil
}

// BatchUpdateDocMetaData 批量更新文档元数据
func (s *Service) BatchUpdateDocMetaData(ctx context.Context, req *knowledgebase_doc_service.BatchUpdateDocMetaDataReq) (*emptypb.Empty, error) {
	//1.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//2.查询知识库下所有文档
	docList, err := orm.GetDocListByKnowledgeId(ctx, req.UserId, req.OrgId, req.KnowledgeId)
	if err != nil {
		log.Errorf("没有查询到文档列表 错误(%v) 参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}
	if len(docList) == 0 {
		return &emptypb.Empty{}, nil
	}
	//3.查询已经设置key的文档
	metaList, err := orm.SelectMetaByKnowledgeId(ctx, req.UserId, req.OrgId, req.KnowledgeId)
	if err != nil {
		log.Errorf("没有查询到文档列表 错误(%v) 参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}
	//4.更新数据map
	var keyValueMap = make(map[string]string)
	for _, item := range req.MetaDataList {
		keyValueMap[item.Key] = item.Value
	}
	addList, updateList := buildUpdateMetaDataParams(knowledge.KnowledgeId, req.OrgId, req.UserId, docList, metaList, keyValueMap)
	//5.文档Id与名称map
	docNameMap := make(map[string]string)
	for _, doc := range docList {
		docNameMap[doc.DocId] = service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	}
	//6.批量更新
	err = orm.UpdateBatchStatusDocMeta(ctx, req.KnowledgeId, docNameMap, addList, updateList, &service.BatchRagDocMetaParams{
		KnowledgeId:   knowledge.KnowledgeId,
		KnowledgeBase: knowledge.Name,
		UserId:        req.UserId,
	})
	if err != nil {
		log.Errorf("update doc tag fail %v", err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocUpdateMetaFailed)
	}
	return &emptypb.Empty{}, nil
}

func buildUpdateMetaDataParams(knowledgeId, orgId, userId string, docList []*model.KnowledgeDoc, metaList []*model.KnowledgeDocMeta, updateMap map[string]string) (addList []*model.KnowledgeDocMeta,
	updateList []*model.KnowledgeDocMeta) {
	existMetaMap := make(map[string]map[string]*model.KnowledgeDocMeta)
	keyTypeMap := make(map[string]string)
	if len(metaList) > 0 {
		for _, meta := range metaList {
			if _, exists := existMetaMap[meta.DocId]; !exists {
				existMetaMap[meta.DocId] = make(map[string]*model.KnowledgeDocMeta)
			}
			if len(updateMap[meta.Key]) > 0 {
				existMetaMap[meta.DocId][meta.Key] = meta
				keyTypeMap[meta.Key] = meta.ValueType
			}
		}
	}
	for _, doc := range docList {
		dataMap := existMetaMap[doc.DocId]
		if len(dataMap) > 0 {
			for metaKey, metaInfo := range dataMap {
				updateList = append(updateList, &model.KnowledgeDocMeta{
					MetaId:    metaInfo.MetaId,
					DocId:     doc.DocId,
					Key:       metaKey,
					Value:     updateMap[metaKey],
					ValueType: metaInfo.ValueType,
				})
			}
		} else {
			for key, value := range updateMap {
				addList = append(addList, &model.KnowledgeDocMeta{
					KnowledgeId: knowledgeId,
					MetaId:      generator.GetGenerator().NewID(),
					DocId:       doc.DocId,
					Key:         key,
					Value:       value,
					ValueType:   keyTypeMap[key],
					Rule:        "",
					OrgId:       orgId,
					UserId:      userId,
					CreatedAt:   time.Now().UnixMilli(),
					UpdatedAt:   time.Now().UnixMilli(),
				})
			}

		}
	}

	return
}
func checkDocStatus(docList []*model.KnowledgeDoc) ([]uint32, []*model.KnowledgeDoc, error) {
	var docIdList []uint32
	var docResultList []*model.KnowledgeDoc
	for _, doc := range docList {
		if doc.Status == model.DocProcessing {
			return nil, nil, errors.New("解析中的文档无法删除")
		}
		docIdList = append(docIdList, doc.Id)
		docResultList = append(docResultList, doc)
	}
	return docIdList, docResultList, nil
}

// buildDocListResp 构造知识库文档列表
func buildDocListResp(list []*model.KnowledgeDoc, importTaskList []*model.KnowledgeImportTask, total int64, pageSize int32, pageNum int32) *knowledgebase_doc_service.GetDocListResp {
	segmentConfigMap := buildSegmentConfigMap(importTaskList)
	var retList = make([]*knowledgebase_doc_service.DocInfo, 0)
	if len(list) > 0 {
		for _, item := range list {
			retList = append(retList, buildDocInfo(item, segmentConfigMap))
		}
	}
	return &knowledgebase_doc_service.GetDocListResp{
		Total:    total,
		Docs:     retList,
		PageSize: pageSize,
		PageNum:  pageNum,
	}
}

func buildDocInfo(item *model.KnowledgeDoc, segmentConfigMap map[string]*model.SegmentConfig) *knowledgebase_doc_service.DocInfo {
	return &knowledgebase_doc_service.DocInfo{
		DocId:         item.DocId,
		DocName:       item.Name,
		DocSize:       item.FileSize,
		DocType:       item.FileType,
		KnowledgeId:   item.KnowledgeId,
		UploadTime:    util2.Time2Str(item.CreatedAt),
		Status:        int32(util.BuildDocRespStatus(item.Status)),
		ErrorMsg:      item.ErrorMsg,
		SegmentMethod: buildSegmentMethod(item, segmentConfigMap),
	}
}

func buildSegmentMethod(knowledgeDoc *model.KnowledgeDoc, configMap map[string]*model.SegmentConfig) string {
	config := configMap[knowledgeDoc.ImportTaskId]
	if config == nil || config.SegmentMethod == "" {
		return model.CommonSegmentMethod
	}
	return config.SegmentMethod
}

// buildSegmentConfigMap 构造分词配置
func buildSegmentConfigMap(importTaskList []*model.KnowledgeImportTask) map[string]*model.SegmentConfig {
	retMap := make(map[string]*model.SegmentConfig)
	if len(importTaskList) == 0 {
		return retMap
	}
	for _, importTask := range importTaskList {
		var config = &model.SegmentConfig{}
		err := json.Unmarshal([]byte(importTask.SegmentConfig), config)
		if err != nil {
			log.Errorf("SegmentConfig process error %s", err.Error())
			continue
		}
		retMap[importTask.ImportId] = config
	}
	return retMap
}

func removeDuplicateMeta(metaDataList []*knowledgebase_doc_service.MetaData) []*knowledgebase_doc_service.MetaData {
	if len(metaDataList) == 0 {
		return metaDataList
	}
	return lo.UniqBy(metaDataList, func(item *knowledgebase_doc_service.MetaData) string {
		return item.Key
	})
}

// buildImportTask 构造导入任务
func buildImportTask(req *knowledgebase_doc_service.ImportDocReq) (*model.KnowledgeImportTask, error) {
	//是否是自动分段类型
	if autoSegmentType(req.DocSegment.SegmentType, req.DocSegment.SegmentMethod) {
		req.DocSegment.Overlap = 0.0
		req.DocSegment.MaxSplitter = 4000
	}
	segmentConfig, err := json.Marshal(req.DocSegment)
	if err != nil {
		return nil, err
	}
	analyzer, err := json.Marshal(&model.DocAnalyzer{
		AnalyzerList: req.DocAnalyzer,
	})
	if err != nil {
		return nil, err
	}
	docList := make([]*model.DocInfo, 0)
	for _, docInfo := range req.DocInfoList {
		docList = append(docList, &model.DocInfo{
			DocId:   docInfo.DocId,
			DocName: docInfo.DocName,
			DocUrl:  docInfo.DocUrl,
			DocType: docInfo.DocType,
			DocSize: docInfo.DocSize,
		})
	}
	docImportInfo, err := json.Marshal(&model.DocImportInfo{
		DocInfoList: docList,
	})
	if err != nil {
		return nil, err
	}

	preprocess, err := json.Marshal(&model.DocPreProcess{
		PreProcessList: req.DocPreprocess,
	})
	if err != nil {
		return nil, err
	}
	var docImportMetaData string
	if len(req.DocMetaDataList) > 0 {
		metaList := make([]*model.KnowledgeDocMeta, 0)
		for _, metaData := range req.DocMetaDataList {
			metaList = append(metaList, &model.KnowledgeDocMeta{
				Key:       metaData.Key,
				Value:     metaData.Value,
				ValueType: metaData.ValueType,
				Rule:      metaData.Rule,
			})
		}
		importMetaDataByte, err := json.Marshal(&model.DocImportMetaData{
			DocMetaDataList: metaList,
		})
		if err != nil {
			return nil, err
		}
		docImportMetaData = string(importMetaDataByte)
	}
	return &model.KnowledgeImportTask{
		ImportId:      generator.GetGenerator().NewID(),
		KnowledgeId:   req.KnowledgeId,
		ImportType:    int(req.DocImportType),
		SegmentConfig: string(segmentConfig),
		DocAnalyzer:   string(analyzer),
		CreatedAt:     time.Now().UnixMilli(),
		UpdatedAt:     time.Now().UnixMilli(),
		DocInfo:       string(docImportInfo),
		OcrModelId:    req.OcrModelId,
		DocPreProcess: string(preprocess),
		MetaData:      docImportMetaData,
		UserId:        req.UserId,
		OrgId:         req.OrgId,
	}, nil
}

func autoSegmentType(segmentType, segmentMethod string) bool {
	if segmentMethod == model.CommonSegmentMethod || segmentMethod == "" {
		return segmentType == "0"
	}
	return false
}

// buildSegmentListResp 构造文档分段列表
func buildSegmentListResp(importTask *model.KnowledgeImportTask, doc *model.KnowledgeDoc,
	segmentListResp *service.ContentListResp, pageNo, pageSize int32, metaDataList []*model.KnowledgeDocMeta,
	segmentImportTask *model.DocSegmentImportTask) (*knowledgebase_doc_service.DocSegmentListResp, error) {
	var config = &model.SegmentConfig{}
	err := json.Unmarshal([]byte(importTask.SegmentConfig), config)
	if err != nil {
		log.Errorf("SegmentConfig process error %s", err.Error())
		return nil, err
	}
	segmentConfigMap := buildSegmentConfigMap([]*model.KnowledgeImportTask{importTask})
	var resp = &knowledgebase_doc_service.DocSegmentListResp{
		FileName:            doc.Name,
		MaxSegmentSize:      int32(config.MaxSplitter),
		SegType:             config.SegmentType,
		CreatedAt:           util2.Time2Str(doc.CreatedAt),
		Splitter:            buildSplitter(config.Splitter),
		PageTotal:           buildPageTotal(int32(segmentListResp.ChunkTotalNum), pageSize),
		SegmentTotalNum:     int32(segmentListResp.ChunkTotalNum),
		ContentList:         buildContentList(segmentListResp.List),
		MetaDataList:        buildMetaList(metaDataList),
		SegmentImportStatus: buildSegmentImportStatus(segmentImportTask),
		SegmentMethod:       buildSegmentMethod(doc, segmentConfigMap),
	}
	return resp, nil
}

func buildChildSegmentListResp(resp *service.ChildContentListResp) (*knowledgebase_doc_service.GetDocChildSegmentListResp, error) {
	var retList = make([]*knowledgebase_doc_service.ChildSegmentInfo, 0)
	for _, item := range resp.ChildContentList {
		retList = append(retList, &knowledgebase_doc_service.ChildSegmentInfo{
			Content:  item.Content,
			ChildId:  item.ContentId,
			ChildNum: int32(item.MetaData.ChildChunkCurrentNum),
			ParentId: resp.ParentChunkId,
		})
	}
	return &knowledgebase_doc_service.GetDocChildSegmentListResp{
		ContentList: retList,
	}, nil
}

func buildSegmentImportStatus(segmentImportTask *model.DocSegmentImportTask) string {
	if segmentImportTask == nil {
		return ""
	}
	if segmentImportTask.Status == model.DocSegmentImportInit {
		return segmentImportingMessage
	} else if segmentImportTask.Status == model.DocSegmentImportImporting {
		timeSpan := time.Now().UnixMilli() - segmentImportTask.UpdatedAt
		if timeSpan < fiveMinutes {
			return segmentImportingMessage
		}
		//大于5分钟标识异步任务中间断了，todo 异步更新任务
	}
	if segmentImportTask.SuccessCount <= 0 {
		return segmentCompleteFail
	}
	if segmentImportTask.TotalCount <= 0 {
		return fmt.Sprintf(segmentPartCompleteFormat, segmentImportTask.SuccessCount)
	}
	return fmt.Sprintf(segmentCompleteFormat, segmentImportTask.SuccessCount, segmentImportTask.TotalCount-segmentImportTask.SuccessCount)
}

func buildMetaList(metaDataList []*model.KnowledgeDocMeta) []*knowledgebase_doc_service.MetaData {
	if len(metaDataList) == 0 {
		return make([]*knowledgebase_doc_service.MetaData, 0)
	}
	return lo.Map(metaDataList, func(item *model.KnowledgeDocMeta, index int) *knowledgebase_doc_service.MetaData {
		var valueType = item.ValueType
		if valueType == "" {
			valueType = model.MetaTypeString
		}
		return &knowledgebase_doc_service.MetaData{
			MetaId:    item.MetaId,
			Key:       item.Key,
			Value:     item.Value,
			ValueType: valueType,
			Rule:      item.Rule,
		}
	})
}

func buildMetaParamsList(metaDataList []*knowledgebase_doc_service.MetaData) []*model.KnowledgeDocMeta {
	if len(metaDataList) == 0 {
		return make([]*model.KnowledgeDocMeta, 0)
	}
	return lo.Map(metaDataList, func(item *knowledgebase_doc_service.MetaData, index int) *model.KnowledgeDocMeta {
		return &model.KnowledgeDocMeta{
			MetaId: item.MetaId,
			Key:    item.Key,
			Value:  item.Value,
		}
	})
}

func buildDeleteMetaKeys(reqMetaList []*knowledgebase_doc_service.MetaData, metaMap map[string]string) []string {
	var deleteKeys []string
	for _, reqMeta := range reqMetaList {
		if reqMeta.Option == MetaOptionDelete {
			if dbKey, exists := metaMap[reqMeta.MetaId]; !exists {
				log.Errorf("metaId %s doesn't exist", reqMeta.MetaId)
				continue
			} else if dbKey == "" {
				log.Errorf("metaId %s dbKey is empty", reqMeta.MetaId)
				continue
			} else {
				deleteKeys = append(deleteKeys, dbKey)
			}
		}
	}
	return deleteKeys
}

func buildDocMetaModelList(metaDataList []*knowledgebase_doc_service.MetaData, orgId, userId, knowledgeId, docId string) (addList []*model.KnowledgeDocMeta,
	updateList []*model.KnowledgeDocMeta, deleteDataIdList []string) {
	if len(metaDataList) == 0 {
		return
	}
	for _, data := range metaDataList {
		if data.Option == MetaOptionDelete {
			deleteDataIdList = append(deleteDataIdList, data.MetaId)
			continue
		}
		if data.Option == MetaOptionUpdate {
			updateList = append(updateList, &model.KnowledgeDocMeta{
				MetaId:    data.MetaId,
				DocId:     docId,
				Key:       data.Key,
				Value:     data.Value,
				ValueType: data.ValueType,
			})
			continue
		}
		if data.Option == MetaOptionAdd {
			addList = append(addList, &model.KnowledgeDocMeta{
				KnowledgeId: knowledgeId,
				MetaId:      generator.GetGenerator().NewID(),
				DocId:       docId,
				Key:         data.Key,
				Value:       data.Value,
				ValueType:   data.ValueType,
				Rule:        "",
				OrgId:       orgId,
				UserId:      userId,
				CreatedAt:   time.Now().UnixMilli(),
				UpdatedAt:   time.Now().UnixMilli(),
			})
		}
	}
	return
}

func buildMetaRagParams(metaDataList []*knowledgebase_doc_service.MetaData) ([]*service.MetaData, error) {
	if len(metaDataList) == 0 {
		return make([]*service.MetaData, 0), nil
	}
	var retList = make([]*service.MetaData, 0)
	for _, data := range metaDataList {
		if data.Option == "delete" {
			continue
		}
		valueData, err := buildValueData(data.ValueType, data.Value)
		if err != nil {
			log.Errorf("buildValueData error %s", err.Error())
			return nil, err
		}
		retList = append(retList, &service.MetaData{
			Key:       data.Key,
			Value:     valueData,
			ValueType: data.ValueType,
		})
	}
	return retList, nil
}

func buildValueData(valueType string, value string) (interface{}, error) {
	switch valueType {
	case model.MetaTypeNumber:
	case model.MetaTypeTime:
		return strconv.ParseInt(value, 10, 64)
	}
	return value, nil
}

func buildSplitter(splitterList []string) string {
	if len(splitterList) == 0 {
		return noSplitter
	}
	return strings.Join(splitterList, " 、 ")
}

func buildPageTotal(totalNum int32, pageSize int32) int32 {
	leftPageSize := totalNum % pageSize
	var leftPage int32 = 0
	if leftPageSize > 0 {
		leftPage = 1
	}
	return totalNum/pageSize + leftPage
}

func buildContentList(contentList []service.FileSplitContent) []*knowledgebase_doc_service.SegmentContent {
	var retList = make([]*knowledgebase_doc_service.SegmentContent, 0)
	for i := 0; i < len(contentList); i++ {
		content := contentList[i]
		retList = append(retList, &knowledgebase_doc_service.SegmentContent{
			Content:    content.Content,
			Available:  content.Status,
			ContentId:  content.ContentId,
			ContentNum: int32(content.MetaData.ChunkCurrentNum),
			Labels:     content.Labels,
			IsParent:   content.IsParent,
			ChildNum:   int32(content.ChildChunkTotalNum),
		})
	}
	return retList
}

func checkUpdateAndAddMetaList(addList []*model.KnowledgeDocMeta, updateList []*service.RagMetaMapKeys, dbMetaList []*model.KnowledgeDocMeta) error {
	// 构造数据库map
	dbKeySet := make(map[string]bool, len(dbMetaList))
	for _, dbMeta := range dbMetaList {
		dbKeySet[dbMeta.Key] = true
	}

	// 校验addList
	addKeySet := make(map[string]bool, len(addList))
	for _, addMeta := range addList {
		if dbKeySet[addMeta.Key] {
			log.Errorf("add meta failed: key %s already exists", addMeta.Key)
			return errors.New("key already exists")
		}
		if addKeySet[addMeta.Key] {
			log.Errorf("add meta failed: key %s repeated", addMeta.Key)
			return errors.New("key repeated")
		}
		addKeySet[addMeta.Key] = true
	}

	// 校验updateList
	for _, updateMeta := range updateList {
		if dbKeySet[updateMeta.NewKey] {
			log.Errorf("update meta failed: key %s already exists", updateMeta.NewKey)
			return errors.New("key already exists")
		}
		if addKeySet[updateMeta.NewKey] {
			log.Errorf("update meta failed: key %s repeated", updateMeta.NewKey)
			return errors.New("key repeated")
		}
	}
	return nil
}

func buildOptionList(metaList []*model.KnowledgeDocMeta, req *knowledgebase_doc_service.UpdateDocMetaDataReq) ([]string, []*service.RagMetaMapKeys, []*model.KnowledgeDocMeta) {
	metaMap := buildKnowledgeMetaMap(metaList)
	deleteList := buildDeleteMetaKeys(req.MetaDataList, metaMap)
	updateList := buildUpdateMetaMap(req.MetaDataList, metaMap)
	addList := buildAddMetaList(req)
	return deleteList, updateList, addList
}
