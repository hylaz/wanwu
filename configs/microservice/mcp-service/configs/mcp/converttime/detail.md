# Time MCP Server

一个模型上下文协议（Model Context Protocol）服务器，提供时间和时区转换功能。该服务器使得大型语言模型（LLMs）能够获取当前时间信息，并使用 IANA 时区名称进行时区转换，同时自动检测系统时区。

### 可用工具

- `get_current_time` - 获取特定时区或系统时区的当前时间。
  - 必需参数：
    - `timezone` (string): IANA 时区名称（例如，'America/New_York'，'Europe/London'）

- `convert_time` - 在时区之间转换时间。
  - 必需参数：
    - `source_timezone` (string): 源 IANA 时区名称
    - `time` (string): 24小时格式的时间（HH:MM）
    - `target_timezone` (string): 目标 IANA 时区名称

## 安装

### 使用 uv（推荐）

使用 `uv` 时无需特定安装。我们将使用 `uvx` 直接运行 *mcp-server-time*。

### 使用 PIP

或者，您可以通过 pip 安装 `mcp-server-time`：

```bash
pip install mcp-server-time
```

安装后，您可以使用以下命令作为脚本运行它：

```bash
python -m mcp_server_time
```

## 配置

### 为 Claude.app 配置

添加到您的 Claude 设置中：

<details>
<summary>使用 uvx</summary>

```json
"mcpServers": {
  "time": {
    "command": "uvx",
    "args": ["mcp-server-time"]
  }
}
```
</details>

<details>
<summary>使用 docker</summary>

```json
"mcpServers": {
  "time": {
    "command": "docker",
    "args": ["run", "-i", "--rm", "mcp/time"]
  }
}
```
</details>

<details>
<summary>使用 pip 安装</summary>

```json
"mcpServers": {
  "time": {
    "command": "python",
    "args": ["-m", "mcp_server_time"]
  }
}
```
</details>

### 为 Zed 配置

添加到您的 Zed settings.json：

<details>
<summary>使用 uvx</summary>

```json
"context_servers": [
  "mcp-server-time": {
    "command": "uvx",
    "args": ["mcp-server-time"]
  }
],
```
</details>

<details>
<summary>使用 pip 安装</summary>

```json
"context_servers": {
  "mcp-server-time": {
    "command": "python",
    "args": ["-m", "mcp_server_time"]
  }
},
```
</details>

### 自定义 - 系统时区

默认情况下，服务器会自动检测您的系统时区。您可以通过在配置中的 `args` 列表中添加参数 `--local-timezone` 来覆盖此设置。

示例：
```json
{
  "command": "python",
  "args": ["-m", "mcp_server_time", "--local-timezone=America/New_York"]
}
```

## 示例交互

1. 获取当前时间：
```json
{
  "name": "get_current_time",
  "arguments": {
    "timezone": "Europe/Warsaw"
  }
}
```
响应：
```json
{
  "timezone": "Europe/Warsaw",
  "datetime": "2024-01-01T13:00:00+01:00",
  "is_dst": false
}
```

2. 在时区之间转换时间：
```json
{
  "name": "convert_time",
  "arguments": {
    "source_timezone": "America/New_York",
    "time": "16:30",
    "target_timezone": "Asia/Tokyo"
  }
}
```
响应：
```json
{
  "source": {
    "timezone": "America/New_York",
    "datetime": "2024-01-01T12:30:00-05:00",
    "is_dst": false
  },
  "target": {
    "timezone": "Asia/Tokyo",
    "datetime": "2024-01-01T12:30:00+09:00",
    "is_dst": false
  },
  "time_difference": "+13.0h",
}
```

## 调试

您可以使用 MCP 检查器来调试服务器。对于 uvx 安装：

```bash
npx @modelcontextprotocol/inspector uvx mcp-server-time
```

或者如果您在特定目录中安装了该包或正在开发它：

```bash
cd path/to/servers/src/time
npx @modelcontextprotocol/inspector uv run mcp-server-time
```

## Claude 的问题示例

1. "现在几点了？"（将使用系统时区）
2. "东京现在几点？"
3. "当纽约是下午4点时，伦敦几点？"
4. "将东京时间上午9:30 转换为纽约时间"

## 构建

Docker 构建：

```bash
cd src/time
docker build -t mcp/time .
```

## 贡献

我们鼓励贡献，以帮助扩展和改进 mcp-server-time。无论您想添加新的时间相关工具、增强现有功能，还是改善文档，您的意见都是宝贵的。



欢迎提交拉取请求！欢迎贡献新想法、修复错误或增强功能，使 mcp-server-time 更加强大和实用。

## 许可证

mcp-server-time 采用 MIT 许可证。这意味着您可以自由使用、修改和分发该软件，但需遵循 MIT 许可证的条款和条件。有关更多详细信息，请参见项目存储库中的 LICENSE 文件。
