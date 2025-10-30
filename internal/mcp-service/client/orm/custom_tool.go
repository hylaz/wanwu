package orm

import (
	"context"
	"errors"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client/model"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client/orm/sqlopt"
	"gorm.io/gorm"
)

func (c *Client) CreateCustomTool(ctx context.Context, customTool *model.CustomTool) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 检查是否已存在相同的记录
		if err := sqlopt.SQLOptions(
			sqlopt.WithName(customTool.Name),
			sqlopt.WithToolSquareID(customTool.ToolSquareId),
			sqlopt.WithOrgID(customTool.OrgID),
			sqlopt.WithUserID(customTool.UserID),
		).Apply(tx).First(&model.CustomTool{}).Error; err == nil {
			return toErrStatus("mcp_create_duplicate_custom_tool")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return toErrStatus("mcp_create_custom_tool_err", err.Error())
		}
		// 创建
		if err := tx.Create(customTool).Error; err != nil {
			return toErrStatus("mcp_create_custom_tool_err", err.Error())
		}
		return nil
	})
}

func (c *Client) GetCustomTool(ctx context.Context, customTool *model.CustomTool) (*model.CustomTool, *err_code.Status) {
	info := &model.CustomTool{}
	if err := sqlopt.SQLOptions(
		sqlopt.WithID(customTool.ID),
		sqlopt.WithOrgID(customTool.OrgID),
		sqlopt.WithUserID(customTool.UserID),
	).Apply(c.db).WithContext(ctx).First(info).Error; err != nil {
		return nil, toErrStatus("mcp_get_custom_tool_info_err", err.Error())
	}
	return info, nil
}

func (c *Client) ListCustomTools(ctx context.Context, orgID, userID, name string) ([]*model.CustomTool, *err_code.Status) {
	var customToolInfos []*model.CustomTool
	if err := sqlopt.SQLOptions(
		sqlopt.WithOrgID(orgID),
		sqlopt.WithUserID(userID),
		sqlopt.LikeName(name),
		sqlopt.WithToolSquareIDEmpty(),
	).Apply(c.db).WithContext(ctx).Find(&customToolInfos).Error; err != nil {
		return nil, toErrStatus("mcp_get_custom_tool_list_err", err.Error())
	}
	return customToolInfos, nil
}

func (c *Client) ListCustomToolsByCustomToolIDs(ctx context.Context, ids []uint32) ([]*model.CustomTool, *err_code.Status) {
	var customToolInfos []*model.CustomTool
	if err := sqlopt.WithIDs(ids).Apply(c.db).WithContext(ctx).Find(&customToolInfos).Error; err != nil {
		return nil, toErrStatus("mcp_get_custom_tool_list_err", err.Error())
	}
	return customToolInfos, nil
}

func (c *Client) UpdateCustomTool(ctx context.Context, customTool *model.CustomTool) *err_code.Status {
	if err := sqlopt.SQLOptions(
		sqlopt.WithID(customTool.ID),
	).Apply(c.db).WithContext(ctx).Model(customTool).Updates(map[string]interface{}{
		"name":               customTool.Name,
		"description":        customTool.Description,
		"schema":             customTool.Schema,
		"privacy_policy":     customTool.PrivacyPolicy,
		"api_key":            customTool.APIKey,
		"auth_type":          customTool.AuthType,
		"custom_header_name": customTool.CustomHeaderName,
		"type":               customTool.Type,
	}).Error; err != nil {
		return toErrStatus("mcp_update_custom_tool_err", err.Error())
	}
	return nil
}

func (c *Client) DeleteCustomTool(ctx context.Context, ID uint32) *err_code.Status {
	if err := sqlopt.WithID(ID).Apply(c.db).WithContext(ctx).Delete(&model.CustomTool{}).Error; err != nil {
		return toErrStatus("mcp_delete_custom_tool_err", err.Error())
	}
	return nil
}

func (c *Client) ListBuiltinTools(ctx context.Context, orgID, userID string) ([]*model.CustomTool, *err_code.Status) {
	var customToolInfos []*model.CustomTool
	if err := sqlopt.SQLOptions(
		sqlopt.WithOrgID(orgID),
		sqlopt.WithUserID(userID),
		sqlopt.WithToolSquareIDNotEmpty(),
	).Apply(c.db).WithContext(ctx).Find(&customToolInfos).Error; err != nil {
		return nil, toErrStatus("mcp_get_custom_tool_list_err", err.Error())
	}
	return customToolInfos, nil
}

func (c *Client) GetBuiltinTool(ctx context.Context, customTool *model.CustomTool) (*model.CustomTool, *err_code.Status) {
	info := &model.CustomTool{}
	if err := sqlopt.SQLOptions(
		sqlopt.WithToolSquareID(customTool.ToolSquareId),
		sqlopt.WithOrgID(customTool.OrgID),
		sqlopt.WithUserID(customTool.UserID),
	).Apply(c.db).WithContext(ctx).First(info).Error; err != nil {
		return nil, toErrStatus("mcp_get_custom_tool_info_err", err.Error())
	}
	return info, nil
}
