# 🧾 MedOCR-MCP-SSE · 医疗影像OCR服务

基于 PaddleOCR 与 FastMCP 框架的**医疗文档影像OCR服务**，支持文本识别、表格解析、结构化文档理解和批量图像处理。通过 SSE 接口提供轻量级实时服务，适用于医学文档智能分析、电子病历结构化等场景。

---

## 🚀 快速预览

- ✅ 支持 PDF / 多图识别  
- ✅ 文本、表格、完整文档结构识别  
- ✅ FastMCP 工具接口，支持多端协作  
- ✅ SSE 实时推送，低延迟任务响应  

---

## 🧠 背景：基于 Med-MaaS 中医药大模型项目

本项目为中医药大模型生态系统的一部分，旨在构建中医药 + AI 的多模态感知与认知服务基础设施。

---

## 🧩 服务功能概览

| 功能模块             | 描述                           |
|----------------------|--------------------------------|
| `ocr_text`           | 文字识别（纯文本输出）         |
| `ocr_table`          | 表格识别（HTML结构）           |
| `ocr_document`       | 结构化文档识别（表格 + 文本）    |
| `ocr_image_folder`   | 文件夹批量图像识别              |
| `check_file_exists`  | 检查文件路径是否有效           |
| `list_supported_formats` | 返回支持的输入/输出格式  |
| `get_ocr_config`     | 获取当前OCR系统配置信息        |

---

## 📦 安装依赖

建议在虚拟环境中执行：

```bash
pip install paddleocr paddlepaddle pypdfium2 opencv-python fastmcp
```

---

## ▶️ 启动服务

```bash
python medocr_mcp_sse.py
```

启动后服务默认监听地址：

```
http://0.0.0.0:9077/sse
```

日志文件写入 `medocr_mcp_sse.log`

---

## 📁 项目结构

```bash
.
├── medocr_mcp_sse.py           # 主服务入口
├── medocr_mcp_sse.log          # 日志记录文件
├── requirements.txt            # 可选：依赖列表
└── README.md                   # 项目说明文档
```

---

## 📌 依赖说明

| 依赖库       | 说明                           |
|--------------|--------------------------------|
| `PaddleOCR`  | 基础OCR模型库（结构化识别）    |
| `pypdfium2`  | PDF渲染与图像转换               |
| `opencv-python` | 图像读取与预处理           |
| `fastmcp`    | MCP服务通信框架（支持SSE）     |
