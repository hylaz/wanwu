# MedImage-resource-mcp

> 🏥 医学影像智能结构化与处理服务（MCP资源模块）

---

## 📘 简介

`MedImage-resource-mcp` 是一个专为医学影像数据处理设计的模型上下文协议（MCP）服务模块，支持 DICOM 等影像数据的读取、结构化解析、格式转换与基础分析。该模块可广泛用于影像 AI 模型训练、医学图像质控、数据平台建设等任务场景。

---

## 🏷️ 模块信息

- **MCP Square ID**: `medimage_resource`
- **名称**: 医学影像资源库
- **分类**: `medical`
- **来源**: `modelcontextprotocol`
- **关键特性**:
  - 支持加载 DICOM 图像与像素数组
  - 提取患者/设备/序列相关元数据
  - 支持 DICOM → NIfTI 转换
  - 获取图像分辨率、体素尺寸等结构信息
  - 适配 AI 训练/分析/质控流程

---

## 📦 安装与运行

### ✅ 安装模块

```bash
pip install mcp-server-medimage
