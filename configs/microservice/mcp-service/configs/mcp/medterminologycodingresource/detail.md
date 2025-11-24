# MedTerminologyCoding-resource-mcp

> 🧾 医学术语编码结构化资源服务模块（MCP）

---

## 📘 简介

`MedTerminologyCoding-resource-mcp` 是一个标准医学术语编码资源模块，支持 ICD-10、ICD-11、SNOMED CT、LOINC 等主流国际医学编码体系的**术语匹配、映射、解析与结构化返回**功能。它可被语言模型、EMR 系统、知识图谱平台或智能问诊引擎直接调用，实现临床文本的标准化对接与术语抽取。

---

## 🏷️ 模块信息

- **MCP Square ID**: `medterminologycoding_resource`
- **模块名称**: 医学术语编码资源库
- **分类**: `medical`
- **来源**: `modelcontextprotocol`
- **主要特性**:
  - 多编码体系支持：ICD-10、ICD-11、SNOMED CT、LOINC 等
  - 支持中英文输入
  - 自然语言术语 → 标准编码映射
  - 编码 → 名称 + 分类结构反查
  - 模糊匹配、自适应纠错支持

---

## 📦 安装与运行

### ✅ 安装模块

```bash
pip install mcp-server-medcoding
