package util

import (
	"bytes"
	"encoding/base64"
	"mime/multipart"
	"net/http"
)

// Base64ToFileHeader Base64 转 *multipart.FileHeader → 清理临时文件
func Base64ToFileHeader(b64, filename string) (*multipart.FileHeader, error) {
	rawFileData, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	// 构造multipart表单
	body := &bytes.Buffer{}
	formWriter := multipart.NewWriter(body)

	fw, err := formWriter.CreateFormFile("file", filename) // 字段名"file"可根据业务调整
	if err != nil {
		return nil, err
	}
	if _, err := fw.Write(rawFileData); err != nil {
		return nil, err
	}
	if err := formWriter.Close(); err != nil {
		return nil, err
	}

	formReader := multipart.NewReader(body, formWriter.Boundary())
	form, err := formReader.ReadForm(int64(body.Len() + 1024))
	if err != nil {
		return nil, err
	}

	files := form.File["file"]
	if len(files) == 0 {
		return nil, http.ErrMissingFile
	}

	return files[0], nil
}
