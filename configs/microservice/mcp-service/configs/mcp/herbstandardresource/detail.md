# HerbStandard-resource-mcp

> 🌿 中药材、药对、功效标准化智能资源接口服务（MCP模块）

---

## 📘 简介

`HerbStandard-resource-mcp` 是一个基于模型上下文协议（MCP）的服务模块，专注于中药材、中成药、药对搭配、配伍禁忌、药食同源等领域的**标准化结构化知识查询服务**。该模块可被大型语言模型（LLM）或中医智能体调用，为智能问答、中药推荐、方剂审核等任务提供高质量的基础知识支持。

---

## 🏷️ 模块信息

- **MCP Square ID**: `herbstandard_resource`
- **模块名称**: 中药标准资源库
- **分类**: `medical`
- **来源**: `modelcontextprotocol`
- **核心特性**:
  - 中药材别名统一解析
  - 支持查询性味归经、功效主治、来源等信息
  - 可校验药对/方剂是否存在配伍禁忌
  - 支持按功效反向推荐药材
  - 提供常用药对搭配及现代研究依据

---

## 📦 安装与运行

### ✅ 安装方式

```bash
pip install mcp-server-herb-standard
