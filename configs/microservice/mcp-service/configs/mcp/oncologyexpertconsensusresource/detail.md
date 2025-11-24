# OncologyExpertConsensus-resource-mcp

> 📚 肿瘤专家共识文献结构化资源模块（MCP）

---

## 📘 简介

`OncologyExpertConsensus-resource-mcp` 是一个面向肿瘤疾病领域的专家共识知识服务模块，支持对临床指南与共识文献进行结构化管理、筛选、摘要提取与推荐要点解析。该模块可以被医疗大模型、知识图谱、问答系统或CDSS系统集成调用。

---

## 🏷️ 模块信息

- **MCP Square ID**: `oncologyexpertconsensus_resource`
- **模块名称**: 肿瘤专家共识资源库
- **分类**: `medical`
- **来源**: `modelcontextprotocol`
- **功能特性**:
  - 支持癌种、年份、发布机构等维度筛选共识文献
  - 获取文献结构化摘要、核心推荐意见
  - 提取推荐条目与证据等级（如 I-A, II-B 等）
  - 适配临床辅助决策与科研支持场景

---

## 📦 安装与使用

### ✅ 安装模块

```bash
pip install mcp-server-oncology-consensus
