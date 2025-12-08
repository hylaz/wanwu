package sqlopt

import (
	"gorm.io/gorm"
)

type SqlOptions []SQLOption

func SQLOptions(opts ...SQLOption) SqlOptions {
	return opts
}

func (s SqlOptions) Apply(db *gorm.DB, model interface{}) *gorm.DB {
	if model != nil {
		db = db.Model(model)
	}
	for _, opt := range s {
		db = opt.Apply(db)
	}
	return db
}

type SQLOption interface {
	Apply(db *gorm.DB) *gorm.DB
}

type funcSQLOption func(db *gorm.DB) *gorm.DB

func (f funcSQLOption) Apply(db *gorm.DB) *gorm.DB {
	return f(db)
}

func WithID(id uint32) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})
}

func WithKnowledgeID(id string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("knowledge_id = ?", id)
	})
}

func WithQAPairID(id string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("qa_pair_id = ?", id)
	})
}

func WithQuestionMd5(questionMd5 string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(questionMd5) > 0 {
			return db.Where("question_md5 = ?", questionMd5)
		}
		return db
	})
}

func WithOverKnowledgePermission(id int) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("permission_type >= ?", id)
	})
}

func WithPermissionId(id string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("permission_id = ?", id)
	})
}

func WithoutKnowledgeID(knowledgeId string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(knowledgeId) == 0 {
			return db
		}
		return db.Where("knowledge_id != ?", knowledgeId)
	})
}

func WithKnowledgeIDList(idList []string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(idList) > 0 {
			return db.Where("knowledge_id IN (?)", idList)
		}
		return db
	})
}

func WithCategory(category int) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("category = ?", category)
	})
}

func WithTagID(id string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("tag_id = ?", id)
	})
}

func WithSplitterID(id string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("splitter_id = ?", id)
	})
}

func WithImportID(id string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("import_id = ?", id)
	})
}

func WithExportID(id string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("export_id = ?", id)
	})
}

func WithImportIDs(idList []string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("import_id in ?", idList)
	})
}

func WithDocIDs(ids []string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("doc_id in ?", ids)
	})
}

func WithDocIDsNonEmpty(ids []string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(ids) == 0 {
			return db
		}
		return db.Where("doc_id in ?", ids)
	})
}

func WithQAPairIDsNonEmpty(ids []string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(ids) == 0 {
			return db
		}
		return db.Where("qa_pair_id in ?", ids)
	})
}

func WithFileTypeFilter(fileType string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(fileType) > 0 {
			return db.Where("file_type != ?", fileType)
		}
		return db
	})
}

func WithQAPairIDs(ids []string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("qa_pair_id in ?", ids)
	})
}

func WithDocID(id string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("doc_id = ?", id)
	})
}

func WithKey(key string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("`key` = ?", key)
	})
}

func WithType(metaValueType string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("`type` = ?", metaValueType)
	})
}

func WithNonType(metaValueType string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("value_type != ?", metaValueType)
	})
}

func WithIDs(ids []uint32) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	})
}

func WithMetaId(metaId string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("meta_id = ?", metaId)
	})
}

func WithMetaIds(metaIds []string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("meta_id IN ?", metaIds)
	})
}

func WithOrgID(orgID string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("org_id = ?", orgID)
	})
}

func WithUserID(userID string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	})
}

// WithPermit 权限查询条件
func WithPermit(orgID, userID string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(orgID) > 0 {
			db = db.Where("org_id = ?", orgID)
		}
		if len(userID) > 0 {
			db = db.Where("user_id = ?", userID)
		}
		return db
	})
}

func WithStatusList(status []int) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(status) == 0 {
			return db
		} else if len(status) == 1 {
			return db.Where("status = ?", status[0])
		}
		return db.Where("status IN ?", status)
	})
}

func WithStatus(status int) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if status == -1 {
			return db
		}
		return db.Where("status = ?", status)
	})
}

func WithoutStatus(status int) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("status != ?", status)
	})
}

func WithName(name string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(name) > 0 {
			return db.Where("name = ?", name)
		}
		return db
	})
}

func WithoutID(id uint32) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if id != 0 {
			return db.Where("id != ?", id)
		}
		return db
	})
}

func WithValue(value string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(value) > 0 {
			return db.Where("value = ?", value)
		}
		return db
	})
}

func WithNonEmptyValue() SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("value != ''")
	})
}

func WithNameOrValue(name, value string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(name) > 0 || len(value) > 0 {
			return db.Where("name = ? OR value = ?", name, value)
		}
		return db
	})
}

func WithNameOrAliasLike(name string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(name) > 0 {
			// 使用 OR 条件组合模糊查询
			return db.Where("name LIKE ? OR alias LIKE ?", "%"+name+"%", "%"+name+"%")
		}
		return db
	})
}

func WithFilePathMd5(filePathMd5 string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if len(filePathMd5) > 0 {
			return db.Where("file_path_md5 = ?", filePathMd5)
		}
		return db
	})
}

func LikeName(name string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if name != "" {
			return db.Where("name LIKE ?", "%"+name+"%")
		}
		return db
	})
}

func LikeMetaValue(value string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if value != "" {
			return db.Where("value_main LIKE ?", "%"+value+"%")
		}
		return db
	})
}

func LikeQuestion(question string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if question != "" {
			return db.Where("question LIKE ?", "%"+question+"%")
		}
		return db
	})
}

func LikeTag(tag string) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		if tag != "" {
			return db.Where("tag LIKE ?", "%"+tag+"%")
		}
		return db
	})
}

func WithDelete(deleted int) SQLOption {
	return funcSQLOption(func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted = ?", deleted)
	})
}
