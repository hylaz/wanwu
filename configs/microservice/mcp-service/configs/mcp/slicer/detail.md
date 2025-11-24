# Slicer-MCP-SSE 工具集

基于 FastMCP 的标准化医学图像预处理与基础分析智能体工具，专为 DICOM/NIfTI 等医学影像多场景预处理和快速分析设计。

---

## 简介

本项目封装了典型医学图像分析流程，包括：图像配准、空间重采样、三维裁剪、快速阈值分割和常用滤波等。每个处理流程以标准 MCP 工具形式暴露，可被智能体、自动化流水线或第三方平台远程自动调用。

- 支持全自动批量医学影像处理和一键集成到AI工作流
- 支持 DICOM、NIfTI 及 ZIP 等主流医学数据格式
- 代码清晰、易扩展、API参数规范、便于医学和AI工程集成

---

## 安装依赖

```bash
pip install fastmcp SimpleITK ants pydicom numpy

## API 工具接口

每个工具均通过 FastMCP 智能体自动发现与调用：

### 1. image_registration

- **功能**：基于 ANTsPy 的医学图像空间配准
- **参数**：
  - `moving_path`：待配准图像路径
  - `fixed_path`：参考图像路径
  - `output_path`：配准结果保存路径
  - `mode`：配准类型，可选 `'Rigid'`、`'Affine'`、`'SyN'`

### 2. resample_image

- **功能**：医学图像空间重采样
- **参数**：
  - `image_path`：原始图像路径
  - `spacing`：新像素间距，格式如 `'1,1,1'`
  - `output_path`：输出图像路径

### 3. crop_volume

- **功能**：三维体积ROI裁剪
- **参数**：
  - `image_path`：输入图像路径
  - `roi`：ROI，格式 `'xmin,xmax,ymin,ymax,zmin,zmax'`
  - `output_path`：输出图像路径

### 4. threshold_segmentation

- **功能**：基于阈值的二值分割
- **参数**：
  - `image_path`：输入图像路径
  - `lower`：下阈值
  - `upper`：上阈值
  - `output_path`：输出掩膜路径

### 5. basic_filter

- **功能**：基础图像滤波
- **参数**：
  - `image_path`：输入图像路径
  - `output_path`：输出图像路径
  - `filter_type`：滤波类型，可选 `'gaussian'`、`'mean'`、`'median'`、`'bilateral'`、`'laplacianofgaussian'`
  - `param`：滤波参数

---

## 命令行启动

```bash
python mcp_tool.py
