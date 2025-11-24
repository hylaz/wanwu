# Perplexity Ask MCP Server

一个MCP服务器实现，集成了Sonar API，为Claude提供无与伦比的实时、全网研究。

## 工具

- **perplexity_ask**
  - 与Sonar API进行对话，以进行实时网页搜索。
  - **输入:**
    - `messages` (数组): 一组对话消息。
      - 每条消息必须包含：
        - `role` (字符串): 消息的角色（例如，`system`、`user`、`assistant`）。
        - `content` (字符串): 消息的内容。

## 配置

### 第一步：

克隆此仓库：

```bash
git clone git@github.com:ppl-ai/modelcontextprotocol.git
```

导航到`perplexity-ask`目录并安装必要的依赖：

```bash
cd modelcontextprotocol/perplexity-ask && npm install
```

### 第二步：获取Sonar API密钥

1. 注册一个Sonar API账户。
2. 按照账户设置说明生成您的API密钥。
3. 将API密钥设置为环境变量`PERPLEXITY_API_KEY`。

### 第三步：配置Claude桌面

1. 在这里下载Claude桌面。

2. 将以下内容添加到您的`claude_desktop_config.json`中：

```json
{
  "mcpServers": {
    "perplexity-ask": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "PERPLEXITY_API_KEY",
        "mcp/perplexity-ask"
      ],
      "env": {
        "PERPLEXITY_API_KEY": "YOUR_API_KEY_HERE"
      }
    }
  }
}
```

### NPX

```json
{
  "mcpServers": {
    "perplexity-ask": {
      "command": "npx",
      "args": [
        "-y",
        "server-perplexity-ask"
      ],
      "env": {
        "PERPLEXITY_API_KEY": "YOUR_API_KEY_HERE"
      }
    }
  }
}
```

您可以使用以下命令访问该文件：

```bash
vim ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

### 第四步：构建Docker镜像

Docker构建：

```bash
docker build -t mcp/perplexity-ask:latest -f Dockerfile .
```

### 第五步：测试

让我们确保Claude桌面能够识别我们在`perplexity-ask`服务器中暴露的两个工具。您可以通过查找锤子图标来做到这一点：


点击锤子图标后，您应该能看到与文件系统MCP服务器一起提供的工具：

如果您看到这两个工具，这意味着集成是活跃的。恭喜！这意味着Claude现在可以向Perplexity提问。您可以像使用Perplexity网页应用一样简单地使用它。

### 第六步：高级参数

目前，使用的搜索参数是默认参数。您可以直接在`index.ts`脚本中修改任何搜索参数。有关详细信息，请参阅官方API文档。

### 故障排除

Claude文档提供了一个优秀的故障排除指南，您可以参考。不过，您仍然可以通过api@perplexity.ai与我们联系以获取额外支持或报告错误。

## 许可证

此MCP服务器遵循MIT许可证。这意味着您可以自由使用、修改和分发该软件，但需遵循MIT许可证的条款和条件。有关更多详细信息，请参见项目仓库中的LICENSE文件。
