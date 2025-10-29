package response

import (
	"encoding/json"
	"net/http"
)

type KnowledgeListResp struct {
	KnowledgeList []*KnowledgeInfo `json:"knowledgeList"`
}

type CreateKnowledgeResp struct {
	KnowledgeId string `json:"knowledgeId"`
}

type KnowledgeHitResp struct {
	Prompt     string             `json:"prompt"`     //提示词列表
	SearchList []*ChunkSearchList `json:"searchList"` //种种结果
	Score      []float64          `json:"score"`      //打分信息
}

type RagKnowledgeResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CommonRagKnowledgeError(err error) ([]byte, int) {
	resp := RagKnowledgeResp{Code: 1, Message: err.Error()}
	marshal, err := json.Marshal(resp)
	if err != nil {
		return []byte(err.Error()), http.StatusBadRequest
	}
	return marshal, http.StatusBadRequest
}

type EmbeddingModelInfo struct {
	ModelId string `json:"modelId"`
}

type KnowledgeInfo struct {
	KnowledgeId        string              `json:"knowledgeId"`        //知识库id
	Name               string              `json:"name"`               //知识库名称
	OrgName            string              `json:"orgName"`            //知识库所属名称
	Description        string              `json:"description"`        //知识库描述
	DocCount           int                 `json:"docCount"`           //文档数量
	EmbeddingModelInfo *EmbeddingModelInfo `json:"embeddingModelInfo"` //embedding模型信息
	KnowledgeTagList   []*KnowledgeTag     `json:"knowledgeTagList"`   //知识库标签列表
	CreateUserId       string              `json:"createUserId"`
	CreateAt           string              `json:"createAt"`       //创建时间
	PermissionType     int32               `json:"permissionType"` //权限类型:0: 查看权限; 10: 编辑权限; 20: 授权权限,数值不连续的原因防止后续有中间权限，目前逻辑 授权权限>编辑权限>查看权限
	Share              bool                `json:"share"`          //是分享，还是私有
}

type KnowledgeMetaData struct {
	Key  string `json:"key"`  // key
	Type string `json:"type"` // type(time, string, number)
}

type ChunkSearchList struct {
	Title            string          `json:"title"`
	Snippet          string          `json:"snippet"`
	KnowledgeName    string          `json:"knowledgeName"`
	ChildContentList []*ChildContent `json:"childContentList"`
	ChildScore       []float64       `json:"childScore"`
}

type ChildContent struct {
	ChildSnippet string  `json:"childSnippet"`
	Score        float64 `json:"score"`
}

type GetKnowledgeMetaSelectResp struct {
	MetaList []*KnowledgeMetaItem `json:"knowledgeMetaList"`
}

type KnowledgeMetaItem struct {
	MetaId        string `json:"metaId"`
	MetaKey       string `json:"metaKey"`
	MetaValueType string `json:"metaValueType"`
	MetaValue     string `json:"metaValue"` // 确定值
}

type KnowledgeMetaValueListResp struct {
	KnowledgeMetaValues []*KnowledgeMetaValues `json:"knowledgeMetaValues"`
}

type KnowledgeMetaValues struct {
	MetaId        string   `json:"metaId"`
	MetaKey       string   `json:"metaKey"`
	MetaValue     []string `json:"metaValue"` // 确定值
	MetaValueType string   `json:"metaValueType"`
}
