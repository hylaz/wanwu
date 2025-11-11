package response

type ReportPageResult struct {
	List      []*ReportInfo `json:"list"`      // 社区报告内容列表
	Total     int32         `json:"total"`     // 社区报告数量：如果为0显示-
	PageNo    int           `json:"pageNo"`    // 当前页码
	PageSize  int           `json:"pageSize"`  // 每页数量
	CreatedAt string        `json:"createdAt"` // 生成时间：unix时间戳，若为空串显示-
	Status    int32         `json:"status"`    // 状态：0.未生成(-) 1.生成中 2.已生成 3.生成失败
}

type ReportInfo struct {
	ContentId string `json:"contentId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}
