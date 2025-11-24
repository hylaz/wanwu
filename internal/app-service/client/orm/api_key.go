package orm

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/app-service/client/model"
	"github.com/UnicomAI/wanwu/internal/app-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/pkg/util"
)

func (c *Client) GetApiKeyList(ctx context.Context, userId, orgId, appId, appType string) ([]*model.ApiKey, *errs.Status) {
	var apiKeys []*model.ApiKey
	if err := sqlopt.SQLOptions(
		sqlopt.WithAppType(appType),
		sqlopt.WithAppID(appId),
		sqlopt.WithOrgID(orgId),
		sqlopt.WithUserID(userId),
	).Apply(c.db.WithContext(ctx)).
		Find(&apiKeys).Error; err != nil {
		return nil, toErrStatus("app_api_keys_get", err.Error())
	}
	return apiKeys, nil
}

func (c *Client) DelApiKey(ctx context.Context, apiId uint32) *errs.Status {
	if err := sqlopt.WithID(apiId).Apply(c.db.WithContext(ctx)).Delete(&model.ApiKey{}).Error; err != nil {
		return toErrStatus("app_api_key_delete", util.Int2Str(apiId), err.Error())
	}
	return nil
}

func (c *Client) GenApiKey(ctx context.Context, userId, orgId, appId, appType, apiKey string) (*model.ApiKey, *errs.Status) {
	newKey := &model.ApiKey{
		OrgID:   orgId,
		UserID:  userId,
		AppID:   appId,
		AppType: appType,
		ApiKey:  apiKey,
	}
	if err := c.db.WithContext(ctx).Create(newKey).Error; err != nil {
		return nil, toErrStatus("app_api_keys_gen", err.Error())
	}
	return newKey, nil
}

func (c *Client) GetApiKeyByKey(ctx context.Context, apiKey string) (*model.ApiKey, *errs.Status) {
	ret := &model.ApiKey{}
	if err := sqlopt.WithApiKey(apiKey).Apply(c.db).WithContext(ctx).First(ret).Error; err != nil {
		return nil, toErrStatus("app_api_key_get_by_key", err.Error())
	}
	return ret, nil
}
