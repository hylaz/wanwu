MCP - 12306-MCP车票查询工具

基于 Model Context Protocol (MCP) 的12306购票搜索服务器。提供了简单的API接口，允许大模型利用接口搜索12306购票信息。

🚩Features
----------

| 功能描述                | 状态     |
| ----------------------- | -------- |
| 查询12306购票信息       | ✅ 已完成 |
| 过滤列车信息            | ✅ 已完成 |
| 过站查询                | ✅ 已完成 |
| 中转查询                | 🚧 计划内 |
| 其余接口，欢迎提feature | 🚧 计划内 |

![Image 7](https://s2.loli.net/2025/04/15/UjbrG5esaSEmJxN.jpg)

![Image 8](https://s2.loli.net/2025/04/15/rm1j8zX7sqiyafP.jpg)

⚙️Installation
--------------

```
git clone https://github.com/Joooook/12306-mcp.git
npm i
```

▶️Quick Start
-------------

### CLI

```
npm run build
node ./build/index.js
```

### MCP sever configuration

```
{
    "mcpServers": {
        "12306-mcp": {
            "command": "npx",
            "args": [
                "-y",
                "12306-mcp"
            ]
        }
    }
}
```

👉️Reference
------------

*   [modelcontextprotocol/modelcontextprotocol](https://github.com/modelcontextprotocol/modelcontextprotocol)
*   [modelcontextprotocol/typescript-sdk](https://github.com/modelcontextprotocol/typescript-sdk)

