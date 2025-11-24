package orm

import (
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/app-service/client/model"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient(db *gorm.DB) (*Client, error) {
	// auto migrate
	if err := db.AutoMigrate(
		model.AppConversation{},
		model.ApiKey{},
		model.AppHistory{},
		model.App{},
		model.AppFavorite{},
		model.SensitiveWordTable{},
		model.SensitiveWordVocabulary{},
		model.AppUrl{},
	); err != nil {
		return nil, err
	}
	return &Client{
		db: db,
	}, nil
}

type ApiKey struct {
	ApiId     string `json:"apiId"`
	CreatedAt int64  `json:"createdAt"`
	ApiKey    string `json:"apiKey" `
}

func toErrStatus(key string, args ...string) *err_code.Status {
	return &err_code.Status{
		TextKey: key,
		Args:    args,
	}
}

type ExplorationAppInfo struct {
	AppId       string
	AppType     string
	CreatedAt   int64
	UpdatedAt   int64
	IsFavorite  bool
	PublishType string
	UserID      string
}

type SensitiveWordTableWithWord struct {
	model.SensitiveWordTable
	SensitiveWords []string
}
