HeFeng 天气 MCP 服务器
-------------------------

一个通过 HeFeng 天气 API 提供中国各地天气预报数据的 Model Context Protocol 服务器。

功能
--------

*   获取实时天气数据
*   获取每小时天气预报（24小时/72小时/168小时）
*   获取每日天气预报（3天/7天/10天/15天/30天）
*   支持通过经纬度坐标查询位置
*   完全中文天气描述

API
---

此 MCP 服务器提供以下工具：

### get-weather

获取特定位置的天气预报数据。

使用 MCP 主机（例如 Claude Desktop）
---------------------------------------

将以下内容添加到你的 claude_desktop_config.json 文件中

NPX
---

```
{
  "mcpServers": {
    "hefeng-weather": {
      "command": "npx",
      "args": ["hefeng-mcp-weather@latest", "--apiKey=${API_KEY}"]
    }
  }
}
```

许可证
-------

该 MCP 服务器根据 MIT 许可证进行许可。这意味着您可以自由使用、修改和分发软件，但须遵守 MIT 许可证的条款和条件。有关更多详细信息，请参阅项目存储库中的 LICENSE 文件。