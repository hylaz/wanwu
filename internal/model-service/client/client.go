package client

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/model-service/client/model"
)

type IClient interface {
	ImportModel(ctx context.Context, req *model.ModelImported) *errs.Status
	UpdateModel(ctx context.Context, req *model.ModelImported) *errs.Status
	DeleteModel(ctx context.Context, req *model.ModelImported) *errs.Status
	ChangeModelStatus(ctx context.Context, req *model.ModelImported) *errs.Status
	GetModel(ctx context.Context, req *model.ModelImported) (*model.ModelImported, *errs.Status)
	ListModels(ctx context.Context, req *model.ModelImported) ([]*model.ModelImported, *errs.Status)
	ListTypeModels(ctx context.Context, req *model.ModelImported) ([]*model.ModelImported, *errs.Status)

	GetModelById(ctx context.Context, modelId uint32) (*model.ModelImported, *errs.Status)
	GetModelByIds(ctx context.Context, modelIds []uint32) ([]*model.ModelImported, *errs.Status)
}
