# Brave Search MCP Server

一个集成了 Brave Search API 的 MCP 服务器实现，提供网页和本地搜索功能。

## 特性

- **网页搜索**：一般查询、新闻、文章，带有分页和新鲜度控制
- **本地搜索**：查找商家、餐厅和服务，提供详细信息
- **灵活过滤**：控制结果类型、安全级别和内容新鲜度
- **智能回退**：当未找到本地结果时，本地搜索会自动回退到网页搜索

## 工具

- **brave_web_search**
  - 执行带有分页和过滤的网页搜索
  - 输入：
    - `query` (string)：搜索词
    - `count` (number, optional)：每页结果数（最多 20）
    - `offset` (number, optional)：分页偏移（最多 9）

- **brave_local_search**
  - 搜索本地商家和服务
  - 输入：
    - `query` (string)：本地搜索词
    - `count` (number, optional)：结果数量（最多 20）
  - 如果未找到本地结果，则自动回退到网页搜索


## 配置

### 获取 API 密钥
1. 注册一个 Brave Search API 账户
2. 选择一个计划（免费套餐可用，每月 2,000 次查询）
3. 从 开发者仪表板 生成你的 API 密钥

### 与 Claude Desktop 一起使用
将以下内容添加到你的 `claude_desktop_config.json`：

### Docker

```json
{
  "mcpServers": {
    "brave-search": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "BRAVE_API_KEY",
        "mcp/brave-search"
      ],
      "env": {
        "BRAVE_API_KEY": "YOUR_API_KEY_HERE"
      }
    }
  }
}
```

### NPX

```json
{
  "mcpServers": {
    "brave-search": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-brave-search"
      ],
      "env": {
        "BRAVE_API_KEY": "YOUR_API_KEY_HERE"
      }
    }
  }
}
```


## 构建

Docker 构建：

```bash
docker build -t mcp/brave-search:latest -f src/brave-search/Dockerfile .
```

## 许可证

该 MCP 服务器根据 MIT 许可证授权。这意味着你可以自由使用、修改和分发该软件，但需遵循 MIT 许可证的条款和条件。有关更多详细信息，请参见项目库中的 LICENSE 文件。
