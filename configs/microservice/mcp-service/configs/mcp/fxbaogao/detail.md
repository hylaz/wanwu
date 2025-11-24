# 📊 发现报告（MCP Tool）

本工具集成于 MCP 协议，提供对发现报告网站研究报告的搜索与内容提取能力，适用于金融、产业研究、投资分析等场景。

---

## ✨ 工具一：`search_reports`

通过关键词、作者、机构名称、时间范围等条件检索研究报告列表。

### 🔧 参数说明

| 参数名        | 类型              | 必填 | 说明 |
|---------------|-------------------|------|------|
| `keywords`    | `str`             | 否   | 搜索关键词，支持中英文。 |
| `authors`     | `List[str]`       | 否   | 作者姓名列表，如 `["张三", "李四"]`。 |
| `org_names`   | `List[str]`       | 否   | 机构名称列表，如 `["中信证券", "华泰证券"]`。 |
| `start_time`  | `int`             | 否   | 起始时间，毫秒级时间戳，如 `1640995200000` 表示 2022-01-01 00:00:00。 |
| `end_time`    | `int` / `str`     | 否   | 结束时间，支持毫秒时间戳或相对时间字符串：<br>• `"last3day"`<br>• `"last7day"`<br>• `"last1mon"`<br>• `"last3mon"`<br>• `"last1year"` |
| `page_size`   | `int`             | 否   | 返回结果数量，默认 10，最大 100。 |

### 📥 使用示例

```python
# 按关键词搜索
search_reports(keywords="人工智能")

# 按机构搜索
search_reports(org_names=["中信证券"])

# 搜索最近一周某位作者的报告
search_reports(authors=["王磊"], end_time="last7day")

# 精确时间段搜索
search_reports(
    keywords="新能源",
    start_time=1748707200000,
    end_time=1749398399999
)
```

## ✨ 工具二：get_report_content

根据报告 ID（`doc_id`）获取研报的详细内容与总结信息。

---

### 🔧 参数说明

| 参数名   | 类型   | 必填 | 说明 |
|----------|--------|------|------|
| `doc_id` | `int`  | ✅ 是 | 研报文档 ID，来自 `search_reports` 返回结果中的 `docId` 字段。 |

---

### 📥 使用示例

```python
# 获取该研报的内容
content = await get_report_content(doc_id)
```

## 服务器配置：
```json
{
    "mcpServers": {
        "fxbaogao-mcp": {
            "command": "uv",
            "args": [
                "run",
                "report.py"
            ]
        }
    }
}
```

## 注意事项
本工具仅供学习和研究使用，请勿用于商业目的
请遵守必应的使用条款和相关法律法规