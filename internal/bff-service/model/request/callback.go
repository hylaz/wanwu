package request

type FileUrlConvertBase64Req struct {
	FileUrl string `form:"fileUrl" json:"fileUrl" validate:"required"` // 文件URL
}

func (f *FileUrlConvertBase64Req) Check() error {
	return nil
}
