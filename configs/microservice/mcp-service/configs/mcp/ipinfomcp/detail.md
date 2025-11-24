# IP 地理位置 MCP 服务器

这是一个简单的[模型上下文协议](https://modelcontextprotocol.io/)服务器，使用[ipinfo.io](https://ipinfo.io/) API 来获取关于一个 IP 地址的详细信息。 这可以用来确定用户的大致位置以及他们所使用的网络。

## 安装

您需要创建一个 token 来使用 IPInfo API。如果您还没有，请在 https://ipinfo.io/signup 注册一个免费账户。

要与 Claude 桌面版一起使用，请将以下内容添加到您的 `claude_desktop_config.json` 文件中的 `mcpServers` 部分：

```
    "ipinfo": {
      "command": "uvx",
      "args": [
        "--from",
        "git+https://github.com/briandconnelly/mcp-server-ipinfo.git",
        "mcp-server-ipinfo"
      ],
      "env": {
        "IPINFO_API_TOKEN": "<YOUR TOKEN HERE>"
      }
    }
```

## 组件

### 工具

- ```
  get_ip_details
  ```

  : 此工具用于获取有关 IP 地址的详细信息。

  - **输入:** `ip`: 要获取信息的 IP 地址。
  - **输出:** `IPDetails`: 包含有关 IP 的详细信息（包括位置、组织和国家详情）的 Pydantic 模型。

### 资源

*没有包含自定义资源*

### 提示

*没有包含自定义提示*

## 许可证

MIT 许可证 - 详情请参阅 LICENSE 文件。

## 免责声明

此项目与 [IPInfo](https://ipinfo.io/) 无关。