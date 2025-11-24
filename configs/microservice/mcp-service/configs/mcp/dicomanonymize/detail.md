# DicomAnonymize-MCP-SSE

## 简介

**DicomAnonymize-MCP-SSE** 是基于 FastMCP 的医学影像 DICOM 自动脱敏服务，支持**元数据标准化脱敏**（DICOM PS3.15）与**像素级可识别特征自动脱敏**。  
适用于医学AI、科研数据共享、自动化数据流等场景，API可被各类AI智能体或自动化系统直接调用。

---

## 功能特性

- **DICOM元数据一键脱敏**：自动清除姓名、ID、医院、医生、设备、描述、日期等敏感标签，支持UID重编码和日期偏移。
- **像素级可识别内容脱敏**：OCR自动检测像素内烧录的文字并无损遮盖，支持2D/3D CT、MR等主流模态。
- **MCP工具接口，API极简易用**：服务自动注册为MCP工具，兼容智能体和RPA/自动化流水线。
- **SSE流式推理，高并发安全**：适用于大规模DICOM隐私合规处理。

---

## 安装依赖

```bash
pip install fastmcp pydicom paddleocr scikit-image opencv-python numpy
