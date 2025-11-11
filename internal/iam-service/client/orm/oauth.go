package orm

import (
	"context"
	"errors"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/model"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/pkg/util"
	"gorm.io/gorm"
)

func (c *Client) CreateOauthApp(ctx context.Context, req *model.OauthApp) *errs.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		if err := sqlopt.SQLOptions(
			sqlopt.WithUserID(req.UserID),
			sqlopt.WithName(req.Name),
		).Apply(tx).First(&model.OauthApp{}).Error; err != gorm.ErrRecordNotFound {
			if err == nil {
				return toErrStatus("oauth_app_create", "app name already exists")
			}
			return toErrStatus("oauth_app_create", err.Error())
		}
		req.ClientID = util.GenUUID()
		req.ClientSecret = util.GenUUID()
		req.Status = true
		if err := tx.Create(req).Error; err != nil {
			return toErrStatus("oauth_app_create", err.Error())
		}
		return nil
	})
}

func (c *Client) DeleteOauthApp(ctx context.Context, clientID string) *errs.Status {
	if err := sqlopt.WithClientID(clientID).Apply(c.db.WithContext(ctx)).Delete(&model.OauthApp{}).Error; err != nil {
		return toErrStatus("oauth_app_delete", clientID, err.Error())
	}
	return nil
}

func (c *Client) UpdateOauthApp(ctx context.Context, req *model.OauthApp) *errs.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		var existingApp model.OauthApp
		err := sqlopt.SQLOptions(
			sqlopt.WithName(req.Name),
		).Apply(tx).First(&existingApp).Error
		if err == nil && existingApp.ClientID != req.ClientID {
			return toErrStatus("oauth_app_update", req.ClientID, "app name already exists")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return toErrStatus("oauth_app_update", req.ClientID, err.Error())
		}

		updates := map[string]interface{}{
			"name":         req.Name,
			"redirect_uri": req.RedirectURI,
			"description":  req.Description,
		}

		if updateErr := sqlopt.WithClientID(req.ClientID).Apply(tx).Model(&model.OauthApp{}).Updates(updates).Error; updateErr != nil {
			return toErrStatus("oauth_app_update", req.ClientID, updateErr.Error())
		}
		return nil
	})
}

func (c *Client) GetOauthAppList(ctx context.Context, userID uint32) ([]*OauthAppInfo, *errs.Status) {
	var apps []*OauthAppInfo
	if err := sqlopt.WithUserID(userID).Apply(c.db.WithContext(ctx)).Model(&model.OauthApp{}).Find(&apps).Error; err != nil {
		return nil, toErrStatus("oauth_app_list", "get list failed", err.Error())
	}
	return apps, nil
}

func (c *Client) UpdateOauthAppStatus(ctx context.Context, clientID string, status bool) *errs.Status {
	updateErr := sqlopt.WithClientID(clientID).Apply(c.db.WithContext(ctx)).Model(&model.OauthApp{}).Update("status", status).Error
	if updateErr != nil {
		return toErrStatus("oauth_app_status_update", clientID, updateErr.Error())
	}
	return nil
}

func (c *Client) GetOauthApp(ctx context.Context, clientID string) (*OauthAppInfo, *errs.Status) {
	var app *OauthAppInfo
	if err := sqlopt.WithClientID(clientID).Apply(c.db.WithContext(ctx)).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, toErrStatus("oauth_app_get", clientID, "oauth app not found")
		}
		return nil, toErrStatus("oauth_app_get", clientID, err.Error())
	}
	return app, nil
}
