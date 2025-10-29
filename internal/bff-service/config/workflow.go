package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigDir = "configs/microservice/bff-service/configs/workflow-temp"
)

type WorkflowTempConfig struct {
	TemplateId string `json:"templateId" mapstructure:"templateId"`
	Category   string `json:"category" mapstructure:"category"`
	Avatar     string `json:"avatar"`
	Name       string `json:"name"`
	SchemaPath string `json:"schemaPath" mapstructure:"schemaPath"`
	Schema     string `json:"-" mapstructure:"-"`
	Desc       string `json:"desc" mapstructure:"desc"`
	Author     string `json:"author" mapstructure:"author"`
	Summary    string `json:"summary" mapstructure:"summary"`
	Feature    string `json:"feature" mapstructure:"feature"`
	Scenario   string `json:"scenario" mapstructure:"scenario"`
	Note       string `json:"note" mapstructure:"note"`
}

type WorkflowTemplateSchema struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Schema string `json:"schema"`
}

func (wtf *WorkflowTempConfig) load() error {
	schemaPath := filepath.Join(ConfigDir, wtf.SchemaPath)
	b, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("load workflowtemp %v schema path %v err: %v", wtf.TemplateId, schemaPath, err)
	}
	// 解析外层结构
	var templateSchema WorkflowTemplateSchema
	if err := json.Unmarshal(b, &templateSchema); err != nil {
		return fmt.Errorf("load workflowtemp %v schema unmarshal err: %v", wtf.TemplateId, err)
	}
	wtf.Schema = templateSchema.Schema
	// avatarPath := filepath.Join(ConfigDir, wtf.Avatar)
	// if _, err = os.ReadFile(avatarPath); err != nil {
	//     return fmt.Errorf("load workflow %v avatar path %v err: %v", wtf.TemplateId, avatarPath, err)
	// }
	return nil
}
