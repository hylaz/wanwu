package response

type DocMenu struct {
	Name     string     `json:"name"`     // 目录名称
	Index    string     `json:"index"`    // 目录索引
	Path     string     `json:"path"`     // 目录路径（转码后）
	PathRaw  string     `json:"pathRaw"`  // 目录路径
	Children []*DocMenu `json:"children"` // 目录

	content string
}

func (dm *DocMenu) SetContent(content string) {
	dm.content = content
}

type DocSearchResp struct {
	Title       string             `json:"title"` // 文档名
	ContentList []DocSearchContent `json:"list"`  // 内容列表
}

type DocSearchContent struct {
	Title   string `json:"title"`   // 文档中的子标题
	Content string `json:"content"` // 内容
	Url     string `json:"url"`     // 文档链接
}
