# 获取 MCP Server

一个模型上下文协议（Model Context Protocol）服务器，提供网页内容获取功能。该服务器使得大型语言模型（LLMs）能够从网页中检索和处理内容，将 HTML 转换为 markdown，以便于使用。

fetch 工具会截断响应，但通过使用 `start_index` 参数，您可以指定内容提取的起始位置。这使得模型可以分块读取网页，直到找到所需的信息。

### 可用工具

- `fetch` - 从互联网获取 URL 并提取其内容为 markdown。
    - `url` (字符串，必需): 要获取的 URL
    - `max_length` (整数，可选): 返回的最大字符数（默认: 5000）
    - `start_index` (整数，可选): 从此字符索引开始内容（默认: 0）
    - `raw` (布尔值，可选): 获取未经过 markdown 转换的原始内容（默认: false）

### 提示

- **fetch**
  - 获取一个 URL 并提取其内容为 markdown
  - 参数：
    - `url` (字符串，必需): 要获取的 URL

## 安装

可选：安装 node.js，这将使 fetch 服务器使用更强大的 HTML 简化工具。

### 使用 uv（推荐）

使用 `uv` 时，无需特定安装。我们将使用 `uvx` 直接运行 *mcp-server-fetch*。

### 使用 PIP

或者，您可以通过 pip 安装 `mcp-server-fetch`：

```
pip install mcp-server-fetch
```

安装后，您可以使用以下命令作为脚本运行它：

```
python -m mcp_server_fetch
```

## 配置

### 为 Claude.app 配置

添加到您的 Claude 设置中：

<details>
<summary>使用 uvx</summary>

```json
"mcpServers": {
  "fetch": {
    "command": "uvx",
    "args": ["mcp-server-fetch"]
  }
}
```
</details>

<details>
<summary>使用 docker</summary>

```json
"mcpServers": {
  "fetch": {
    "command": "docker",
    "args": ["run", "-i", "--rm", "mcp/fetch"]
  }
}
```
</details>

<details>
<summary>使用 pip 安装</summary>

```json
"mcpServers": {
  "fetch": {
    "command": "python",
    "args": ["-m", "mcp_server_fetch"]
  }
}
```
</details>

### 自定义 - robots.txt

默认情况下，如果请求来自模型（通过工具），服务器将遵循网站的 robots.txt 文件，但如果请求是用户发起的（通过提示），则不会遵循。通过在配置的 `args` 列表中添加参数 `--ignore-robots-txt` 可以禁用此功能。

### 自定义 - User-agent

默认情况下，根据请求是来自模型（通过工具）还是用户发起（通过提示），服务器将使用以下其中一个 user-agent：
```
ModelContextProtocol/1.0 (Autonomous; +https://github.com/modelcontextprotocol/servers)
```
或
```
ModelContextProtocol/1.0 (User-Specified; +https://github.com/modelcontextprotocol/servers)
```

您可以通过在配置的 `args` 列表中添加参数 `--user-agent=YourUserAgent` 来进行自定义。

### 自定义 - 代理

服务器可以通过使用 `--proxy-url` 参数配置为使用代理。

## 调试

您可以使用 MCP 检查器来调试服务器。对于 uvx 安装：

```
npx @modelcontextprotocol/inspector uvx mcp-server-fetch
```

或者如果您在特定目录中安装了该包或正在开发它：

```
cd path/to/servers/src/fetch
npx @modelcontextprotocol/inspector uv run mcp-server-fetch
```

## 贡献

我们鼓励大家贡献，以帮助扩展和改进 mcp-server-fetch。无论您是想添加新工具、增强现有功能，还是改善文档，您的意见都是宝贵的。



欢迎提交拉取请求！欢迎随时贡献新想法、修复错误或增强功能，使 mcp-server-fetch 更加强大和实用。

## 许可证

mcp-server-fetch 采用 MIT 许可证。这意味着您可以自由使用、修改和分发该软件，前提是遵守 MIT 许可证的条款和条件。有关更多详细信息，请参见项目存储库中的 LICENSE 文件。
