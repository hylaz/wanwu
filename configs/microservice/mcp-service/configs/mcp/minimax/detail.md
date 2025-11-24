
<p align="center">
  官方的 MiniMax 模型上下文协议 (MCP) 服务器，支持与强大的文本转语音和视频/图像生成 API 进行交互。该服务器允许像 Claude Desktop、Cursor、Windsurf、OpenAI Agents 等 MCP 客户端生成语音、克隆声音、生成视频、生成图像等。
</p>

## 使用 MCP 客户端的快速入门
1. 从 MiniMax 获取您的 API 密钥。
2. 安装 `uv`（Python 包管理器），可以使用 `curl -LsSf https://astral.sh/uv/install.sh | sh` 安装，或查看 `uv` 仓库 获取其他安装方法。

### Claude Desktop
前往 `Claude > Settings > Developer > Edit Config > claude_desktop_config.json`，添加以下内容：

```
{
  "mcpServers": {
    "MiniMax": {
      "command": "uvx",
      "args": [
        "minimax-mcp"
      ],
      "env": {
        "MINIMAX_API_KEY": "<insert-your-api-key-here>",
        "MINIMAX_MCP_BASE_PATH": "<local-output-dir-path>",
        "MINIMAX_API_HOST": "https://api.minimaxi.chat"
      }
    }
  }
}
```

如果您使用的是 Windows，您需要在 Claude Desktop 中启用“开发者模式”以使用 MCP 服务器。点击左上角汉堡菜单中的“帮助”，然后选择“启用开发者模式”。

### Cursor
前往 `Cursor -> Preferences -> Cursor Settings -> MCP -> Add new global MCP Server`，添加上述配置。

就这样。您的 MCP 客户端现在可以通过这些工具与 MiniMax 进行交互：

## 示例用法

⚠️ 警告：使用这些工具可能会产生费用。

### 1. 播放一段晚间新闻
<img src="https://public-cdn-video-data-algeng.oss-cn-wulanchabu.aliyuncs.com/Snipaste_2025-04-09_20-07-53.png?x-oss-process=image/resize,p_50/format,webp" style="display: inline-block; vertical-align: middle;"/>

### 2. 克隆一个声音
<img src="https://public-cdn-video-data-algeng.oss-cn-wulanchabu.aliyuncs.com/Snipaste_2025-04-09_19-45-13.png?x-oss-process=image/resize,p_50/format,webp" style="display: inline-block; vertical-align: middle;"/>

### 3. 生成一个视频
<img src="https://public-cdn-video-data-algeng.oss-cn-wulanchabu.aliyuncs.com/Snipaste_2025-04-09_19-58-52.png?x-oss-process=image/resize,p_50/format,webp" style="display: inline-block; vertical-align: middle;"/>
<img src="https://public-cdn-video-data-algeng.oss-cn-wulanchabu.aliyuncs.com/Snipaste_2025-04-09_19-59-43.png?x-oss-process=image/resize,p_50/format,webp" style="display: inline-block; vertical-align: middle; "/>

### 4. 生成图像
<img src="https://public-cdn-video-data-algeng.oss-cn-wulanchabu.aliyuncs.com/gen_image.png?x-oss-process=image/resize,p_50/format,webp" style="display: inline-block; vertical-align: middle;"/>
<img src="https://public-cdn-video-data-algeng.oss-cn-wulanchabu.aliyuncs.com/gen_image1.png?x-oss-process=image/resize,p_50/format,webp" style="display: inline-block; vertical-align: middle; "/>
