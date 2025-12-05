package util

import (
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/xuri/excelize/v2"
)

func ReadExcelColumn(filePath string, columnNo int) ([]string, error) {
	// 1. 打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Errorf("打开Excel文件失败: %v", err)
	}
	defer func() {
		// 关闭文件
		if err := f.Close(); err != nil {
			log.Errorf("关闭Excel文件时出错: %v", err)
		}
	}()

	// 2. 获取工作表列表
	sheets := f.GetSheetList()
	sheet := sheets[0]

	// 3. 获取工作表中的所有行
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	var result []string
	// 4. 遍历行和单元格
	for _, row := range rows {
		result = append(result, row[columnNo])
	}

	return result, nil
}
