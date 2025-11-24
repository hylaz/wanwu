标题：MCP - 必应搜索中文

URL来源：http://modelscope.cn/mcp/servers/@yan5236/bing-cn-mcp-server

Markdown内容：
必应中文MCP
-----------

一款基于MCP（模型上下文协议）的中文微软必应搜索工具。它可以直接搜索必应，并通过像Claude或其他支持MCP的AI工具获取网页内容。

特点
--------

*   **支持中文搜索结果**
*   **无需API密钥**，直接抓取必应搜索结果
*   **提供网页内容获取功能**
*   **轻量级，易于安装和使用**
*   **针对中文用户进行了优化**
*   **支持被Claude等AI工具调用**

安装
------------

### 全局安装

```
npm install -g bing-cn-mcp
```

### 或者使用npx直接运行

```
npx bing-cn-mcp
```

使用方法
-----

### 启动服务器

```
bing-cn-mcp
```

或者使用npx：

```
npx bing-cn-mcp
```

### 在支持MCP的环境中使用

在支持MCP的环境（如Cursor）中，配置MCP服务器以使用它：

1.   找到你的MCP配置文件（例如，`.cursor/mcp.json`）
2.   添加服务器配置：

```
{
  "mcpServers": {
    "bingcn": {
      "command": "npx",
      "args": [
        "bing-cn-mcp"
      ]
    }
  }
}
```

对于Windows用户：

```
{
  "mcpServers": {
    "bingcnmcp": {
        "command": "cmd",
        "args": [
          "/c",
          "npx",
          "bing-cn-mcp"
      ]
    }
  }
}
```

3.   现在你可以在Claude中使用`mcp__bing_search`和`mcp__fetch_webpage`工具了

支持的工具
---------------

### bing_search

搜索必应并检索结果列表。

**参数：**

*   `query`：搜索关键词
*   `num_results`：要返回的结果数量（默认是5）

### fetch_webpage

根据`bing_search`返回的结果ID获取特定网页的内容。

**参数：**

*   `result_id`：从`bing_search`返回的结果ID

自定义配置
--------------------

你可以通过创建`.env`文件来自定义设置，例如：

```
# 用户代理设置
USER_AGENT=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36
```

注意事项
-----

*   一些网站可能采用了反抓取措施，这可能会阻止`fetch_webpage`获取内容
*   本工具仅**用于学习和研究目的**，不得用于商业用途
*   请遵守**必应的服务条款**以及所有适用的法律法规

slcatwujian

许可证
-------

MIT