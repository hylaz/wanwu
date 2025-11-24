# 顺序思维 MCP 服务器

一个提供通过结构化思维过程进行动态和反思性问题解决工具的MCP服务器实现。

## 功能

- 将复杂问题分解为可管理的步骤
- 随着理解的深入修订和完善思路
- 分支进入不同的推理路径
- 动态调整思维总数
- 生成并验证解决方案假设

## 工具

### sequential_thinking

促进详细、逐步的问题解决和分析思考过程。

**输入:**

- `thought` (字符串): 当前思考步骤
- `nextThoughtNeeded` (布尔值): 是否需要另一个思考步骤
- `thoughtNumber` (整数): 当前思考编号
- `totalThoughts` (整数): 估计所需的总思考次数
- `isRevision` (布尔值, 可选): 是否修改之前的思考
- `revisesThought` (整数, 可选): 正在重新考虑哪个思考
- `branchFromThought` (整数, 可选): 分支点的思考编号
- `branchId` (字符串, 可选): 分支标识符
- `needsMoreThoughts` (布尔值, 可选): 如果需要更多思考

## 使用

顺序思维工具设计用于：

- 将复杂问题分解成步骤
- 具有修订空间的规划与设计
- 可能需要路线修正的分析
- 初始时可能不清楚全部范围的问题
- 需要在多个步骤中保持上下文的任务
- 需要过滤掉无关信息的情况

## 配置

### 与 Claude 桌面版一起使用

将以下内容添加到您的 `claude_desktop_config.json` 中：

#### npx

```
{
  "mcpServers": {
    "sequential-thinking": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-sequential-thinking"
      ]
    }
  }
}
```

#### docker

```
{
  "mcpServers": {
    "sequentialthinking": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "mcp/sequentialthinking"
      ]
    }
  }
}
```

## 构建

Docker:

```
docker build -t mcp/sequentialthinking -f src/sequentialthinking/Dockerfile .
```

## 许可证

此MCP服务器根据MIT许可证许可。这意味着您可以在遵守MIT许可证条款和条件的前提下自由使用、修改和分发该软件。更多详情，请参阅项目仓库中的LICENSE文件。