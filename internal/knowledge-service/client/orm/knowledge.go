package orm

import (
	"context"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/generator"
	"github.com/samber/lo"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm/sqlopt"
	async_task "github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/async-task"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/db"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/util"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"
	"gorm.io/gorm"
)

// SelectKnowledgeList 查询知识库列表
func SelectKnowledgeList(ctx context.Context, userId, orgId, name string, tagIdList []string) ([]*model.KnowledgeBase, map[string]int, error) {
	var knowledgeIdList []string
	var err error
	if len(tagIdList) > 0 {
		knowledgeIdList, err = SelectKnowledgeIdByTagId(ctx, tagIdList)
		if err != nil {
			return nil, nil, err
		}
	}
	//查询有权限的知识库列表，获取有权限的知识库id，目前是getALL，没有通过连表实现
	permissionKnowledgeList, err := SelectKnowledgeIdByPermission(ctx, userId, orgId, model.PermissionTypeView)
	if err != nil {
		return nil, nil, err
	}
	if len(permissionKnowledgeList) == 0 {
		return make([]*model.KnowledgeBase, 0), nil, nil
	}
	knowledgeIdList = intersectionKnowledgeIdList(knowledgeIdList, buildPermissionKnowledgeIdList(permissionKnowledgeList))
	var knowledgeList []*model.KnowledgeBase
	err = sqlopt.SQLOptions(sqlopt.WithKnowledgeIDList(knowledgeIdList), sqlopt.LikeName(name), sqlopt.WithDelete(0)).
		Apply(db.GetHandle(ctx), &model.KnowledgeBase{}).
		Order("create_at desc").
		Find(&knowledgeList).
		Error
	if err != nil {
		return nil, nil, err
	}
	return knowledgeList, buildPermissionKnowledgeIdMap(permissionKnowledgeList), nil
}

// SelectKnowledgeById 查询知识库信息,todo
func SelectKnowledgeById(ctx context.Context, knowledgeId, userId, orgId string) (*model.KnowledgeBase, error) {
	var knowledge model.KnowledgeBase
	err := sqlopt.SQLOptions(sqlopt.WithPermit(orgId, userId), sqlopt.WithKnowledgeID(knowledgeId), sqlopt.WithDelete(0)).
		Apply(db.GetHandle(ctx), &model.KnowledgeBase{}).
		First(&knowledge).Error
	if err != nil {
		log.Errorf("SelectKnowledgeById userId %s err: %v", userId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseAccessDenied)
	}
	return &knowledge, nil
}

// SelectKnowledgeByIdList 查询知识库信息
func SelectKnowledgeByIdList(ctx context.Context, knowledgeIdList []string, userId, orgId string) ([]*model.KnowledgeBase, map[string]int, error) {
	//查询有权限的知识库列表，获取有权限的知识库id，目前是getALL，没有通过连表实现
	permissionKnowledgeList, err := SelectKnowledgeIdByPermission(ctx, userId, orgId, model.PermissionTypeView)
	if err != nil {
		return nil, nil, err
	}
	if len(permissionKnowledgeList) == 0 {
		return make([]*model.KnowledgeBase, 0), nil, nil
	}
	knowledgeIdList = intersectionKnowledgeIdList(knowledgeIdList, buildPermissionKnowledgeIdList(permissionKnowledgeList))
	var knowledgeList []*model.KnowledgeBase
	err = sqlopt.SQLOptions(sqlopt.WithKnowledgeIDList(knowledgeIdList), sqlopt.WithDelete(0)).
		Apply(db.GetHandle(ctx), &model.KnowledgeBase{}).
		Find(&knowledgeList).Error
	if err != nil {
		log.Errorf("SelectKnowledgeByIdList userId %s err: %v", userId, err)
		return nil, nil, util.ErrCode(errs.Code_KnowledgeBaseAccessDenied)
	}
	return knowledgeList, buildPermissionKnowledgeIdMap(permissionKnowledgeList), nil
}

// SelectKnowledgeByName 查询知识库信息
func SelectKnowledgeByName(ctx context.Context, knowledgeName, userId, orgId string) (*model.KnowledgeBase, error) {
	var knowledge model.KnowledgeBase
	err := sqlopt.SQLOptions(sqlopt.WithPermit(orgId, userId), sqlopt.WithName(knowledgeName), sqlopt.WithDelete(0)).
		Apply(db.GetHandle(ctx), &model.KnowledgeBase{}).
		First(&knowledge).Error
	if err != nil {
		log.Errorf("SelectKnowledgeByName userId %s err: %v", userId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseAccessDenied)
	}
	return &knowledge, nil
}

// SelectKnowledgeByIdNoDeleteCheck 查询知识库信息
func SelectKnowledgeByIdNoDeleteCheck(ctx context.Context, knowledgeId, userId, orgId string) (*model.KnowledgeBase, error) {
	var knowledge model.KnowledgeBase
	err := sqlopt.SQLOptions(sqlopt.WithPermit(orgId, userId), sqlopt.WithKnowledgeID(knowledgeId)).
		Apply(db.GetHandle(ctx), &model.KnowledgeBase{}).
		First(&knowledge).Error
	if err != nil {
		log.Errorf("SelectKnowledgeById userId %s err: %v", userId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseAccessDenied)
	}
	return &knowledge, nil
}

// CheckSameKnowledgeName 知识库名称是否存在同名
func CheckSameKnowledgeName(ctx context.Context, userId, orgId, name, knowledgeId string) error {
	var count int64
	err := sqlopt.SQLOptions(sqlopt.WithPermit(orgId, userId), sqlopt.WithName(name), sqlopt.WithoutKnowledgeID(knowledgeId), sqlopt.WithDelete(0)).
		Apply(db.GetHandle(ctx), &model.KnowledgeBase{}).
		Count(&count).Error
	if err != nil {
		log.Errorf("KnowledgeNameExist userId %s name %s err: %v", userId, name, err)
		return util.ErrCode(errs.Code_KnowledgeBaseDuplicateName)
	}
	if count > 0 {
		return util.ErrCode(errs.Code_KnowledgeBaseDuplicateName)
	}
	return nil
}

// CreateKnowledge 创建知识库
func CreateKnowledge(ctx context.Context, knowledge *model.KnowledgeBase, embeddingModelId string) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.插入数据
		err := createKnowledge(tx, knowledge)
		if err != nil {
			return err
		}
		//2.插入权限信息
		err = CreateKnowledgeIdPermission(tx, buildKnowledgePermission(knowledge))
		if err != nil {
			return err
		}
		//3.通知rag创建知识库
		return service.RagKnowledgeCreate(ctx, &service.RagCreateParams{
			UserId:           knowledge.UserId,
			Name:             knowledge.Name,
			KnowledgeBaseId:  knowledge.KnowledgeId,
			EmbeddingModelId: embeddingModelId,
		})
	})
}

// UpdateKnowledge 更新知识库
func UpdateKnowledge(ctx context.Context, name, description string, knowledgeBase *model.KnowledgeBase) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.更新数据
		err := updateKnowledge(tx, knowledgeBase.Id, name, description)
		if err != nil {
			return err
		}
		//2.通知rag更新知识库
		return service.RagKnowledgeUpdate(ctx, &service.RagUpdateParams{
			UserId:          knowledgeBase.UserId,
			KnowledgeBaseId: knowledgeBase.KnowledgeId,
			OldKbName:       knowledgeBase.Name,
			NewKbName:       name,
		})
	})
}

// UpdateKnowledgeShareCount 更新知识库分享数量
func UpdateKnowledgeShareCount(tx *gorm.DB, knowledgeId string, count int64) error {
	var updateParams = map[string]interface{}{
		"share_count": count,
	}
	return tx.Model(&model.KnowledgeBase{}).Where("knowledge_id=?", knowledgeId).Updates(updateParams).Error
}

// DeleteKnowledge 删除知识库
func DeleteKnowledge(ctx context.Context, knowledgeBase *model.KnowledgeBase) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.逻辑删除数据
		err := logicDeleteKnowledge(tx, knowledgeBase)
		if err != nil {
			return err
		}
		//2.通知rag更新知识库
		return async_task.SubmitTask(ctx, async_task.KnowledgeDeleteTaskType, &async_task.KnowledgeDeleteParams{
			KnowledgeId: knowledgeBase.KnowledgeId,
		})
	})
}

// ExecuteDeleteKnowledge 删除知识库
func ExecuteDeleteKnowledge(tx *gorm.DB, id uint32) error {
	return tx.Unscoped().Model(&model.KnowledgeBase{}).Where("id = ?", id).Delete(&model.KnowledgeBase{}).Error
}

// UpdateKnowledgeFileInfo 更新知识库文档信息
func UpdateKnowledgeFileInfo(tx *gorm.DB, knowledgeId string, resultList []*model.DocInfo) error {
	var docSize int64
	for _, result := range resultList {
		docSize += result.DocSize
	}
	return tx.Model(&model.KnowledgeBase{}).Where("knowledge_id = ?", knowledgeId).
		Update("doc_size", gorm.Expr("doc_size + ?", docSize)).
		Update("doc_count", gorm.Expr("doc_count + ?", len(resultList))).Error
}

// DeleteKnowledgeFileInfo 删除知识库文档信息
func DeleteKnowledgeFileInfo(tx *gorm.DB, knowledgeId string, resultList []*model.DocInfo) error {
	var docSize int64
	for _, result := range resultList {
		docSize += result.DocSize
	}
	return tx.Model(&model.KnowledgeBase{}).Where("knowledge_id = ?", knowledgeId).
		Update("doc_size", gorm.Expr("doc_size - ?", docSize)).
		Update("doc_count", gorm.Expr("doc_count - ?", len(resultList))).Error
}

func createKnowledge(tx *gorm.DB, knowledge *model.KnowledgeBase) error {
	return tx.Create(knowledge).Error
}

func updateKnowledge(tx *gorm.DB, id uint32, name, description string) error {
	var updateParams = map[string]interface{}{
		"name":        name,
		"description": description,
	}
	return tx.Model(&model.KnowledgeBase{}).Where("id=?", id).Updates(updateParams).Error
}

// 逻辑删除
func logicDeleteKnowledge(tx *gorm.DB, knowledge *model.KnowledgeBase) error {
	var updateParams = map[string]interface{}{
		"deleted": 1,
	}
	return tx.Model(&model.KnowledgeBase{}).Where("id=?", knowledge.Id).Updates(updateParams).Error
}

// buildKnowledgePermission 构建知识库权限信息
func buildKnowledgePermission(knowledge *model.KnowledgeBase) *model.KnowledgePermission {
	return &model.KnowledgePermission{
		PermissionId:   generator.GetGenerator().NewID(),
		KnowledgeId:    knowledge.KnowledgeId,
		GrantUserId:    knowledge.UserId,
		GrantOrgId:     knowledge.OrgId,
		PermissionType: model.PermissionTypeSystem,
		CreatedAt:      knowledge.CreatedAt,
		UpdatedAt:      knowledge.UpdatedAt,
		UserId:         knowledge.UserId,
		OrgId:          knowledge.OrgId,
	}
}

func buildPermissionKnowledgeIdList(permissionList []*model.KnowledgePermission) []string {
	return lo.Map(permissionList, func(item *model.KnowledgePermission, index int) string {
		return item.KnowledgeId
	})
}

func buildPermissionKnowledgeIdMap(permissionList []*model.KnowledgePermission) map[string]int {
	var permissionMap = make(map[string]int)
	for _, permission := range permissionList {
		permissionMap[permission.KnowledgeId] = permission.PermissionType
	}
	return permissionMap
}

// intersectionKnowledgeIdList 计算两个知识库id 列表的交集
func intersectionKnowledgeIdList(knowledgeIdList, permissionKnowledgeIdList []string) []string {
	//特殊逻辑，如果用户没有指定tag，则返回用户有权限的知识库id列表
	if len(knowledgeIdList) == 0 {
		return permissionKnowledgeIdList
	}
	var knowledgeIdMap = make(map[string]bool)
	for _, permissionKnowledgeId := range permissionKnowledgeIdList {
		knowledgeIdMap[permissionKnowledgeId] = true
	}
	var retKnowledgeIdList []string
	for _, knowledgeId := range knowledgeIdList {
		if knowledgeIdMap[knowledgeId] {
			retKnowledgeIdList = append(retKnowledgeIdList, knowledgeId)
		}
	}
	return retKnowledgeIdList
}
