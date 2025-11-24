# Firecrawl MCP Server

一个与 Firecrawl 集成的模型上下文协议 (MCP) 服务器实现，提供网页抓取功能。

> 特别感谢 @vrknetha、@cawstudios 的初始实现！
>
> 您还可以在 MCP.so 的游乐场上体验我们的 MCP 服务器。感谢 MCP.so 的托管以及 @gstarwd 的集成支持。

## 特性

- 支持抓取、爬取、搜索、提取、深度研究和批量抓取
- 支持 JS 渲染的网页抓取
- URL 发现与爬取
- 内容提取的网页搜索
- 自动重试，采用指数退避策略
  - 内置速率限制的高效批处理
- 云 API 的信用使用监控
- 完善的日志系统
- 支持云和自托管的 Firecrawl 实例
- 移动/桌面视口支持
- 智能内容过滤，支持标签的包含/排除

## 安装

### 使用 npx 运行

```bash
env FIRECRAWL_API_KEY=fc-YOUR_API_KEY npx -y firecrawl-mcp
```

### 手动安装

```bash
npm install -g firecrawl-mcp
```

### 在 Cursor 上运行

配置 Cursor 🖥️  
注意：需要 Cursor 版本 0.45.6 及以上  
有关最新的配置说明，请参考官方 Cursor 文档中的 MCP 服务器配置：  
Cursor MCP 服务器配置指南

在 Cursor **v0.45.6** 中配置 Firecrawl MCP

1. 打开 Cursor 设置
2. 转到功能 > MCP 服务器
3. 点击 "+ 添加新 MCP 服务器"
4. 输入以下内容：
   - 名称： "firecrawl-mcp"（或您喜欢的名称）
   - 类型： "command"
   - 命令： `env FIRECRAWL_API_KEY=your-api-key npx -y firecrawl-mcp`

在 Cursor **v0.48.6** 中配置 Firecrawl MCP

1. 打开 Cursor 设置
2. 转到功能 > MCP 服务器
3. 点击 "+ 添加新全局 MCP 服务器"
4. 输入以下代码：
   ```json
   {
     "mcpServers": {
       "firecrawl-mcp": {
         "command": "npx",
         "args": ["-y", "firecrawl-mcp"],
         "env": {
           "FIRECRAWL_API_KEY": "YOUR-API-KEY"
         }
       }
     }
   }
   ```

> 如果您使用 Windows 并遇到问题，请尝试 `cmd /c "set FIRECRAWL_API_KEY=your-api-key && npx -y firecrawl-mcp"`

将 `your-api-key` 替换为您的 Firecrawl API 密钥。如果您还没有，可以在 https://www.firecrawl.dev/app/api-keys 创建一个帐户并获取它。

添加后，刷新 MCP 服务器列表以查看新工具。Composer Agent 会在适当时自动使用 Firecrawl MCP，但您可以通过描述您的网页抓取需求明确请求它。通过 Command+L（Mac）访问 Composer，选择提交按钮旁边的 "Agent"，并输入您的查询。

### 在 Windsurf 上运行

将以下内容添加到您的 `./codeium/windsurf/model_config.json`：

```json
{
  "mcpServers": {
    "mcp-server-firecrawl": {
      "command": "npx",
      "args": ["-y", "firecrawl-mcp"],
      "env": {
        "FIRECRAWL_API_KEY": "YOUR_API_KEY"
      }
    }
  }
}
```

### 通过 Smithery 安装（遗留）

要通过 Smithery 自动安装 Firecrawl for Claude Desktop：

```bash
npx -y @smithery/cli install @mendableai/mcp-server-firecrawl --client claude
```

## 配置

### 环境变量

#### 云 API 所需

- `FIRECRAWL_API_KEY`：您的 Firecrawl API 密钥
  - 使用云 API 时必需（默认）
  - 使用自托管实例时可选，需设置 `FIRECRAWL_API_URL`
- `FIRECRAWL_API_URL`（可选）：自托管实例的自定义 API 端点
  - 如果未提供，将使用云 API（需要 API 密钥）

#### 可选配置

##### 重试配置

- `FIRECRAWL_RETRY_MAX_ATTEMPTS`：最大重试次数（默认：3）
- `FIRECRAWL_RETRY_INITIAL_DELAY`：第一次重试前的初始延迟（毫秒，默认：1000）
- `FIRECRAWL_RETRY_MAX_DELAY`：重试之间的最大延迟（毫秒，默认：10000）
- `FIRECRAWL_RETRY_BACKOFF_FACTOR`：指数退避乘数（默认：2）

##### 信用使用监控

- `FIRECRAWL_CREDIT_WARNING_THRESHOLD`：信用使用警告阈值（默认：1000）
- `FIRECRAWL_CREDIT_CRITICAL_THRESHOLD`：信用使用临界阈值（默认：100）

### 配置示例

用于云 API 使用的自定义重试和信用监控：

```bash
# 云 API 所需
export FIRECRAWL_API_KEY=your-api-key

# 可选重试配置
export FIRECRAWL_RETRY_MAX_ATTEMPTS=5        # 增加最大重试次数
export FIRECRAWL_RETRY_INITIAL_DELAY=2000    # 初始延迟 2 秒
export FIRECRAWL_RETRY_MAX_DELAY=30000       # 最大延迟 30 秒
export FIRECRAWL_RETRY_BACKOFF_FACTOR=3      # 更激进的退避

# 可选信用监控
export FIRECRAWL_CREDIT_WARNING_THRESHOLD=2000    # 在 2000 信用时发出警告
export FIRECRAWL_CREDIT_CRITICAL_THRESHOLD=500    # 在 500 信用时发出临界警告
```

用于自托管实例：

```bash
# 自托管所需
export FIRECRAWL_API_URL=https://firecrawl.your-domain.com

# 自托管的可选身份验证
export FIRECRAWL_API_KEY=your-api-key  # 如果您的实例需要身份验证

# 自定义重试配置
export FIRECRAWL_RETRY_MAX_ATTEMPTS=10
export FIRECRAWL_RETRY_INITIAL_DELAY=500     # 从更快的重试开始
```

### 与 Claude Desktop 一起使用

将以下内容添加到您的 `claude_desktop_config.json`：

```json
{
  "mcpServers": {
    "mcp-server-firecrawl": {
      "command": "npx",
      "args": ["-y", "firecrawl-mcp"],
      "env": {
        "FIRECRAWL_API_KEY": "YOUR_API_KEY_HERE",

        "FIRECRAWL_RETRY_MAX_ATTEMPTS": "5",
        "FIRECRAWL_RETRY_INITIAL_DELAY": "2000",
        "FIRECRAWL_RETRY_MAX_DELAY": "30000",
        "FIRECRAWL_RETRY_BACKOFF_FACTOR": "3",

        "FIRECRAWL_CREDIT_WARNING_THRESHOLD": "2000",
        "FIRECRAWL_CREDIT_CRITICAL_THRESHOLD": "500"
      }
    }
  }
}
```

### 系统配置

服务器包含多个可通过环境变量设置的可配置参数。如果未配置，以下为默认值：

```typescript
const CONFIG = {
  retry: {
    maxAttempts: 3, // 对于速率限制请求的重试次数
    initialDelay: 1000, // 第一次重试前的初始延迟（毫秒）
    maxDelay: 10000, // 重试之间的最大延迟（毫秒）
    backoffFactor: 2, // 指数退避的乘数
  },
  credit: {
    warningThreshold: 1000, // 当信用使用达到此水平时发出警告
    criticalThreshold: 100, // 当信用使用达到此水平时发出临界警报
  },
};
```

这些配置控制：

1. **重试行为**

   - 自动重试因速率限制而失败的请求
   - 使用指数退避以避免对 API 造成过大压力
   - 示例：使用默认设置时，重试将在以下时间进行：
     - 第一次重试：延迟 1 秒
     - 第二次重试：延迟 2 秒
     - 第三次重试：延迟 4 秒（上限为最大延迟）

2. **信用使用监控**
   - 跟踪云 API 使用的 API 信用消耗
   - 在指定阈值时提供警告
   - 帮助防止意外的服务中断
   - 示例：使用默认设置：
     - 在剩余 1000 信用时发出警告
     - 在剩余 100 信用时发出临界警报

### 速率限制和批处理

服务器利用 Firecrawl 的内置速率限制和批处理能力：

- 自动速率限制处理，采用指数退避
- 高效的并行处理批量操作
- 智能请求排队和节流
- 对瞬态错误的自动重试

## 可用工具

### 1. 抓取工具 (`firecrawl_scrape`)

从单个 URL 抓取内容，支持高级选项。

```json
{
  "name": "firecrawl_scrape",
  "arguments": {
    "url": "https://example.com",
    "formats": ["markdown"],
    "onlyMainContent": true,
    "waitFor": 1000,
    "timeout": 30000,
    "mobile": false,
    "includeTags": ["article", "main"],
    "excludeTags": ["nav", "footer"],
    "skipTlsVerification": false
  }
}
```

### 2. 批量抓取工具 (`firecrawl_batch_scrape`)

高效抓取多个 URL，内置速率限制和并行处理。

```json
{
  "name": "firecrawl_batch_scrape",
  "arguments": {
    "urls": ["https://example1.com", "https://example2.com"],
    "options": {
      "formats": ["markdown"],
      "onlyMainContent": true
    }
  }
}
```

响应包括操作 ID 以供状态检查：

```json
{
  "content": [
    {
      "type": "text",
      "text": "批量操作已排队，ID 为：batch_1。使用 firecrawl_check_batch_status 检查进度。"
    }
  ],
  "isError": false
}
```

### 3. 检查批量状态 (`firecrawl_check_batch_status`)

检查批量操作的状态。

```json
{
  "name": "firecrawl_check_batch_status",
  "arguments": {
    "id": "batch_1"
  }
}
```

### 4. 搜索工具 (`firecrawl_search`)

在网上搜索，并可选择从搜索结果中提取内容。

```json
{
  "name": "firecrawl_search",
  "arguments": {
    "query": "your search query",
    "limit": 5,
    "lang": "en",
    "country": "us",
    "scrapeOptions": {
      "formats": ["markdown"],
      "onlyMainContent": true
    }
  }
}
```

### 5. 爬取工具 (`firecrawl_crawl`)

启动异步爬取，支持高级选项。

```json
{
  "name": "firecrawl_crawl",
  "arguments": {
    "url": "https://example.com",
    "maxDepth": 2,
    "limit": 100,
    "allowExternalLinks": false,
    "deduplicateSimilarURLs": true
  }
}
```

### 6. 提取工具 (`firecrawl_extract`)

使用 LLM 能力从网页中提取结构化信息。支持云 AI 和自托管 LLM 提取。

```json
{
  "name": "firecrawl_extract",
  "arguments": {
    "urls": ["https://example.com/page1", "https://example.com/page2"],
    "prompt": "提取产品信息，包括名称、价格和描述",
    "systemPrompt": "您是一个有用的助手，提取产品信息",
    "schema": {
      "type": "object",
      "properties": {
        "name": { "type": "string" },
        "price": { "type": "number" },
        "description": { "type": "string" }
      },
      "required": ["name", "price"]
    },
    "allowExternalLinks": false,
    "enableWebSearch": false,
    "includeSubdomains": false
  }
}
```

示例响应：

```json
{
  "content": [
    {
      "type": "text",
      "text": {
        "name": "示例产品",
        "price": 99.99,
        "description": "这是一个示例产品描述"
      }
    }
  ],
  "isError": false
}
```

#### 提取工具选项：

- `urls`：要提取信息的 URL 数组
- `prompt`：用于 LLM 提取的自定义提示
- `systemPrompt`：引导 LLM 的系统提示
- `schema`：用于结构化数据提取的 JSON 模式
- `allowExternalLinks`：允许从外部链接提取
- `enableWebSearch`：启用网页搜索以获取额外上下文
- `includeSubdomains`：在提取中包括子域名

使用自托管实例时，提取将使用您配置的 LLM。对于云 API，它使用 Firecrawl 的托管 LLM 服务。

### 7. 深度研究工具 (`firecrawl_deep_research`)

使用智能爬取、搜索和 LLM 分析对查询进行深度网络研究。

```json
{
  "name": "firecrawl_deep_research",
  "arguments": {
    "query": "碳捕集技术是如何工作的？",
    "maxDepth": 3,
    "timeLimit": 120,
    "maxUrls": 50
  }
}
```

参数：

- query (string, required)：要探索的研究问题或主题。
- maxDepth (number, optional)：爬取/搜索的最大递归深度（默认：3）。
- timeLimit (number, optional)：研究会话的时间限制（秒，默认：120）。
- maxUrls (number, optional)：要分析的最大 URL 数量（默认：50）。

返回：

- 基于研究生成的最终分析（data.finalAnalysis）。
- 可能还包括研究过程中使用的结构化活动和来源。

### 8. 生成 llms.txt 工具 (`firecrawl_generate_llmstxt`)

为给定域生成标准化的 llms.txt（可选生成 llms-full.txt）文件。该文件定义了大型语言模型应如何与该站点交互。

```json
{
  "name": "firecrawl_generate_llmstxt",
  "arguments": {
    "url": "https://example.com",
    "maxUrls": 20,
    "showFullText": true
  }
}
```

参数：

- url (string, required)：要分析的网站的基本 URL。
- maxUrls (number, optional)：要包含的最大 URL 数量（默认：10）。
- showFullText (boolean, optional)：是否在响应中包含 llms-full.txt 内容。

返回：

- 生成的 llms.txt 文件内容，及可选的 llms-full.txt（data.llmstxt 和/或 data.llmsfulltxt）

## 日志系统

服务器包含全面的日志记录：

- 操作状态和进度
- 性能指标
- 信用使用监控
- 速率限制跟踪
- 错误条件

示例日志消息：

```
[INFO] Firecrawl MCP 服务器成功初始化
[INFO] 开始抓取 URL: https://example.com
[INFO] 批量操作已排队，ID 为：batch_1
[WARNING] 信用使用已达到警告阈值
[ERROR] 超过速率限制，2 秒后重试...
```

## 错误处理

服务器提供强大的错误处理：

- 对瞬态错误的自动重试
- 速率限制处理，采用退避策略
- 详细的错误消息
- 信用使用警告
- 网络弹性

示例错误响应：

```json
{
  "content": [
    {
      "type": "text",
      "text": "错误：超过速率限制。2 秒后重试..."
    }
  ],
  "isError": true
}
```

## 开发

```bash
# 安装依赖
npm install

# 构建
npm run build

# 运行测试
npm test
```

### 贡献

1. Fork 该仓库
2. 创建您的功能分支
3. 运行测试： `npm test`
4. 提交拉取请求

## 许可证

MIT 许可证 - 详见 LICENSE 文件。
