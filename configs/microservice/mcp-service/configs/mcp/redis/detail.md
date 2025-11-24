Redis
-------------

一个提供访问 Redis 数据库的模型上下文协议服务器。此服务器通过一组标准化工具使 LLM 能够与 Redis 键值存储进行交互。


先决条件
-------------

1.  必须安装并运行 Redis 服务器
    *   [下载 Redis](https://redis.io/download)
    *   对于 Windows 用户：使用 [Windows Subsystem for Linux (WSL)](https://redis.io/docs/getting-started/installation/install-redis-on-windows/)  或 [Memurai](https://www.memurai.com/)（兼容 Redis 的 Windows 服务器）
    *   默认端口：6379


常见问题及解决方案
-------------------------

### 连接错误

**ECONNREFUSED**

*   **原因**：Redis 服务器未运行或无法访问
*   **解决方案**：
    *   验证 Redis 是否正在运行： `redis-cli ping` 应返回 "PONG"
    *   检查 Redis 服务状态：`systemctl status redis` （Linux）或  `brew services list` （macOS）
    *   确保默认端口（6379）未被防火墙阻止
    *   验证 Redis URL 格式： `redis://hostname:port`

### 服务器行为

*  服务器实现了带有最大重试次数 5 次的指数退避策略
*  初始重试延迟：1 秒，最大延迟：30 秒
*  服务器将在达到最大重试次数后退出，以防止无限重连循环


组件
----------

### 工具

*   **set**

    *  设置 Redis 键值对，并可选设置过期时间
    *  输入：
        *   `key`  (字符串)：Redis 键
        *   `value` (字符串)：要存储的值
        *   `expireSeconds` (数字, 可选)：过期时间（秒）

*   **get**

    *   从 Redis 中根据键获取值
    *   输入： `key` (字符串)：要检索的 Redis 键

*   **delete**

    *   从 Redis 中删除一个或多个键
    *   输入： `key`  (字符串 | 字符串数组)：要删除的单个键或键数组

*   **list**

    *   列出与模式匹配的 Redis 键
    *   输入： `pattern` (字符串, 可选)：用于匹配键的模式（默认：*）


与 Claude Desktop 一起使用
-------------------------

要将此服务器与 Claude Desktop 应用程序一起使用，请在您的`claude_desktop_config.json`文件的 "mcpServers" 部分添加以下配置：

### Docker

*  在 macOS 上运行 Docker 时，如果服务器在主机网络上运行（例如 localhost），请使用 host.docker.internal
*  可以作为参数指定 Redis URL，默认为 "redis://localhost:6379"

```
{
  "mcpServers": {
    "redis": {
      "command": "docker",
      "args": [
        "run", 
        "-i", 
        "--rm", 
        "mcp/redis", 
        "redis://host.docker.internal:6379"]
    }
  }
}
```

### NPX

```
{
  "mcpServers": {
    "redis": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-redis",
        "redis://localhost:6379"
      ]
    }
  }
}
```


构建
--------

Docker:

```
docker build -t mcp/redis -f src/redis/Dockerfile . 
```


许可证
-------

该 MCP 服务器基于 MIT 许可证发布。这意味着您可以自由使用、修改和分发软件，但需遵守 MIT 许可证的条款和条件。有关更多详细信息，请参阅项目仓库中的 LICENSE 文件。
