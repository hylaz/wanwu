package util

import (
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func FileEOF(err error) bool {
	return errors.Is(err, io.EOF) || (err != nil && err.Error() == "EOF")
}

func BuildFilePath(fileDir, fileExt string) string {
	return fileDir + uuid.NewV4().String() + fileExt
}

func ReplaceLast(s, old, new string) string {
	i := strings.LastIndex(s, old)
	if i == -1 {
		return s
	}
	return s[:i] + new + s[i+len(old):]
}

func Img2base64(imgPath string) (string, error) {
	// 读取图片文件
	data, err := os.ReadFile(imgPath)
	if err != nil {
		return "", err
	}

	// 获取文件扩展名（不含点）
	ext := strings.TrimPrefix(filepath.Ext(imgPath), ".")

	// 对文件内容进行base64编码
	encodedImage := base64.StdEncoding.EncodeToString(data)

	// 构建完整的base64数据URI
	imgBase64Str := "data:image/" + ext + ";base64," + encodedImage
	return imgBase64Str, nil
}
