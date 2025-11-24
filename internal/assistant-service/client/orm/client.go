package orm

import (
	"context"
	"errors"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/gromitlee/access"
	"github.com/gromitlee/depend/v2"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient(db *gorm.DB) (*Client, error) {
	// rbac
	if err := access.InitAccessRBAC0Controller(db); err != nil {
		return nil, err
	}
	// depend
	if err := depend.Init(db); err != nil {
		return nil, err
	}
	// auto migrate
	if err := db.AutoMigrate(
		model.Assistant{},
		model.Conversation{},
		model.AssistantWorkflow{},
		model.AssistantMCP{},
		model.AssistantTool{},
		model.CustomPrompt{},
	); err != nil {
		return nil, err
	}
	return &Client{
		db: db,
	}, nil
}

func (c *Client) transaction(ctx context.Context, fc func(tx *gorm.DB) *err_code.Status) *err_code.Status {
	var status *err_code.Status
	_ = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if status = fc(tx); status != nil {
			return errors.New(status.String())
		}
		return nil
	})
	return status
}

func toErrStatus(code string, args ...string) *err_code.Status {
	return &err_code.Status{
		TextKey: code,
		Args:    args,
	}
}
