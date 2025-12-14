package request

import mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"

type FileUrlConvertBase64Req struct {
	FileUrl      string `form:"fileUrl" json:"fileUrl" validate:"required"` // 文件URL
	AddPrefix    bool   `form:"addPrefix" json:"addPrefix"`                 // 是否添加 data:xxx;base64, 前缀
	CustomPrefix string `form:"customPrefix" json:"customPrefix"`           // 自定义前缀（如 "data:video/mp4;base64,"）
}

func (f *FileUrlConvertBase64Req) Check() error {
	return nil
}

type AudioBase64ConvertTextReq struct {
	File   string                 `form:"file" json:"file" validate:"required"` // base64格式
	Config mp_common.AsrConfigOut `form:"config" json:"config" validate:"required"`
	ApiKey string                 `form:"apiKey" json:"apiKey" validate:"required"`
}

func (u *AudioBase64ConvertTextReq) Check() error {
	return nil
}
