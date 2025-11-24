package service

import (
	"fmt"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

func ImportModel(ctx *gin.Context, userId, orgId string, req *request.ImportOrUpdateModelRequest) error {
	clientReq, err := parseImportAndUpdateClientReq(userId, orgId, req)
	if err != nil {
		return err
	}
	if err = ValidateModel(ctx, clientReq); err != nil {
		return grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("An error occurred during model import validation: Invalid model: %v, err : %v", clientReq.Model, err))
	}
	_, err = model.ImportModel(ctx.Request.Context(), clientReq)
	if err != nil {
		return err
	}
	return nil
}

func UpdateModel(ctx *gin.Context, userId, orgId string, req *request.ImportOrUpdateModelRequest) error {
	if req.ModelId == "" {
		return grpc_util.ErrorStatus(err_code.Code_BFFInvalidArg, "modelId cannot be empty")
	}
	clientReq, err := parseImportAndUpdateClientReq(userId, orgId, req)
	if err != nil {
		return err
	}
	if err = ValidateModel(ctx, clientReq); err != nil {
		return grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("An error occurred during model update validation: Invalid model: %v, err : %v", clientReq.Model, err))
	}
	_, err = model.UpdateModel(ctx, clientReq)
	if err != nil {
		return err
	}
	return nil
}

func DeleteModel(ctx *gin.Context, userId, orgId string, req *request.DeleteModelRequest) error {
	_, err := model.DeleteModel(ctx.Request.Context(), &model_service.DeleteModelReq{
		ModelId: req.ModelId,
		UserId:  userId,
		OrgId:   orgId,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetModel(ctx *gin.Context, userId, orgId string, req *request.GetModelRequest) (*response.ModelInfo, error) {
	resp, err := model.GetModel(ctx.Request.Context(), &model_service.GetModelReq{
		ModelId: req.ModelId,
		UserId:  userId,
		OrgId:   orgId,
	})
	if err != nil {
		return nil, err
	}
	return toModelInfo(ctx, resp)
}

func ListModels(ctx *gin.Context, userId, orgId string, req *request.ListModelsRequest) (*response.ListResult, error) {
	resp, err := model.ListModels(ctx.Request.Context(), &model_service.ListModelsReq{
		Provider:    req.Provider,
		ModelType:   req.ModelType,
		DisplayName: req.DisplayName,
		IsActive:    req.IsActive,
		UserId:      userId,
		OrgId:       orgId,
	})
	if err != nil {
		return nil, err
	}
	list, err := toModelInfos(ctx, resp.Models)
	if err != nil {
		return nil, err
	}
	return &response.ListResult{
		List:  list,
		Total: resp.Total,
	}, nil
}

func ChangeModelStatus(ctx *gin.Context, userId, orgId string, req *request.ModelStatusRequest) error {
	_, err := model.ChangeModelStatus(ctx.Request.Context(), &model_service.ModelStatusReq{
		ModelId:  req.ModelId,
		IsActive: req.IsActive,
		UserId:   userId,
		OrgId:    orgId,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetModelById(ctx *gin.Context, req *request.GetModelByIdRequest) (*response.ModelInfo, error) {
	resp, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{
		ModelId: req.ModelId,
	})
	if err != nil {
		return nil, err
	}
	return toModelInfo(ctx, resp)
}

func ListTypeModels(ctx *gin.Context, userId, orgId string, req *request.ListTypeModelsRequest) (*response.ListResult, error) {
	resp, err := model.ListTypeModels(ctx.Request.Context(), &model_service.ListTypeModelsReq{
		ModelType: req.ModelType,
		UserId:    userId,
		OrgId:     orgId,
	})
	if err != nil {
		return nil, err
	}
	list, err := toModelInfos(ctx, resp.Models)
	if err != nil {
		return nil, err
	}
	return &response.ListResult{
		List:  list,
		Total: resp.Total,
	}, nil
}

func parseImportAndUpdateClientReq(userId, orgId string, req *request.ImportOrUpdateModelRequest) (*model_service.ModelInfo, error) {
	clientReq := &model_service.ModelInfo{
		Provider:      req.Provider,
		ModelId:       req.ModelId,
		ModelType:     req.ModelType,
		Model:         req.Model,
		DisplayName:   req.DisplayName,
		ModelIconPath: req.Avatar.Key,
		PublishDate:   req.PublishDate,
		UserId:        userId,
		OrgId:         orgId,
		IsActive:      true,
		ModelDesc:     req.ModelDesc,
	}
	configStr, err := req.ConfigString()
	if err != nil {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFInvalidArg, err.Error())
	}
	clientReq.ProviderConfig = configStr
	return clientReq, nil
}

func toModelInfos(ctx *gin.Context, models []*model_service.ModelInfo) ([]*response.ModelInfo, error) {
	var ret []*response.ModelInfo
	for _, m := range models {
		info, err := toModelInfo(ctx, m)
		if err != nil {
			return nil, err
		}
		ret = append(ret, info)
	}
	return ret, nil
}

func toModelInfo(ctx *gin.Context, modelInfo *model_service.ModelInfo) (*response.ModelInfo, error) {
	modelConfig, err := mp.ToModelConfig(modelInfo.Provider, modelInfo.ModelType, modelInfo.ProviderConfig)
	if err != nil {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v get model config err: %v", modelInfo.ModelId, err))
	}
	tags, err := mp.ToModelTags(modelInfo.Provider, modelInfo.ModelType, modelInfo.ProviderConfig)
	if err != nil {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v get model tags err: %v", modelInfo.ModelId, err))
	}
	res := &response.ModelInfo{
		ModelId:     modelInfo.ModelId,
		Provider:    modelInfo.Provider,
		Model:       modelInfo.Model,
		ModelType:   modelInfo.ModelType,
		DisplayName: modelInfo.DisplayName,
		Avatar:      CacheAvatar(ctx, modelInfo.ModelIconPath, true),
		PublishDate: modelInfo.PublishDate,
		IsActive:    modelInfo.IsActive,
		UserId:      modelInfo.UserId,
		OrgId:       modelInfo.OrgId,
		CreatedAt:   util.Time2Str(modelInfo.CreatedAt),
		UpdatedAt:   util.Time2Str(modelInfo.UpdatedAt),
		ModelDesc:   modelInfo.ModelDesc,
		Config:      modelConfig,
		Tags:        tags,
	}
	if res.DisplayName == "" {
		res.DisplayName = res.Model
	}
	return res, nil
}
