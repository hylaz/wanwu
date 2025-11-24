标题：MCP - AI Word文档助手

URL来源：http://modelscope.cn/mcp/servers/@GongRzhe/Office-Word-MCP-Server

Markdown内容：
Office-Word-MCP-Server
----------------------

一款基于模型上下文协议（MCP）的服务器，用于创建、读取和操作Microsoft Word文档。该服务器能让AI助手通过标准化接口处理Word文档，提供丰富的文档编辑功能。

![图片1：MCP服务器](https://badge.mcpx.dev/?type=server)

概述
--------

Office-Word-MCP-Server实现了[模型上下文协议](https://modelcontextprotocol.io/)，将Word文档操作以工具和资源的形式开放出来。它充当AI助手与Microsoft Word文档之间的桥梁，支持文档创建、内容添加、格式设置和分析等操作。

### 示例

#### 提示词

![图片2：图片](http://modelscope.cn/mcp/servers/@GongRzhe/Office-Word-MCP-Server)

#### 输出

![图片3：图片](http://modelscope.cn/mcp/servers/@GongRzhe/Office-Word-MCP-Server)

功能特点
--------

### 文档管理

*   创建带有元数据的新Word文档
*   提取文本并分析文档结构
*   查看文档属性和统计信息
*   列出目录中的可用文档
*   创建现有文档的副本

### 内容创建

*   添加不同级别的标题
*   插入带有可选样式的段落
*   创建包含自定义数据的表格
*   添加按比例缩放的图片
*   插入分页符

### 富文本格式设置

*   格式化特定文本部分（粗体、斜体、下划线）
*   更改文本颜色和字体属性
*   为文本元素应用自定义样式
*   在整个文档中搜索和替换文本

### 表格格式设置

*   为表格设置边框和样式
*   创建具有独特格式的标题行
*   应用单元格底纹和自定义边框
*   构建表格以提高可读性

### 高级文档操作

*   删除段落
*   创建自定义文档样式
*   在整个文档中应用统一的格式
*   精确控制特定文本范围的格式

安装步骤
------------

### 前置要求

*   Python 3.8或更高版本
*   pip包管理器

### 基本安装

```
# 克隆仓库
git clone https://github.com/GongRzhe/Office-Word-MCP-Server.git
cd Office-Word-MCP-Server

# 安装依赖
pip install -r requirements.txt
```

### 使用安装脚本

或者，你可以使用提供的安装脚本，它会处理以下事项：

*   检查前置要求
*   设置虚拟环境
*   安装依赖项
*   生成MCP配置

```
python setup_mcp.py
```

在Claude桌面版中使用
-----------------------------

### 配置方法

#### 方法1：本地安装后

1.  安装完成后，将服务器添加到你的Claude桌面版配置文件中：

```
{
  "mcpServers": {
    "word-document-server": {
      "command": "python",
      "args": [
        "/path/to/word_server.py"
      ]
    }
  }
}
```

#### 方法2：不安装（使用uvx）

1.  你也可以配置Claude桌面版，通过uvx包管理器使用服务器，无需本地安装：

```
{
  "mcpServers": {
    "word-document-server": {
      "command": "uvx",
      "args": [
        "--from", "office-word-mcp-server", "word_mcp_server"
      ]
    }
  }
}
```

2.  配置文件位置：

    *   macOS系统：`~/Library/Application Support/Claude/claude_desktop_config.json`
    *   Windows系统：`%APPDATA%\Claude\claude_desktop_config.json`

3.  重启Claude桌面版以加载配置。

### 示例操作

配置完成后，你可以让Claude执行以下操作：

*   “创建一个名为‘report.docx’的新文档，包含标题页”
*   “在我的文档中添加一个标题和三个段落”
*   “插入一个包含销售数据的4x4表格”
*   “将第2段中的‘important’一词格式化为粗体红色”
*   “搜索并替换所有‘old term’为‘new term’”
*   “为章节标题创建自定义样式”
*   “为我文档中的表格应用格式”

API参考
-------------

### 文档创建和属性

```
create_document(filename, title=None, author=None)
get_document_info(filename)
get_document_text(filename)
get_document_outline(filename)
list_available_documents(directory=".")
copy_document(source_filename, destination_filename=None)
```

### 内容添加

```
add_heading(filename, text, level=1)
add_paragraph(filename, text, style=None)
add_table(filename, rows, cols, data=None)
add_picture(filename, image_path, width=None)
add_page_break(filename)
```

### 文本格式设置

```
format_text(filename, paragraph_index, start_pos, end_pos, bold=None, 
            italic=None, underline=None, color=None, font_size=None, font_name=None)
search_and_replace(filename, find_text, replace_text)
delete_paragraph(filename, paragraph_index)
create_custom_style(filename, style_name, bold=None, italic=None, 
                    font_size=None, font_name=None, color=None, base_style=None)
```

### 表格格式设置

```
format_table(filename, table_index, has_header_row=None, 
             border_style=None, shading=None)
```

故障排除
---------------

### 常见问题

1.  **样式缺失**

    *   某些文档可能缺少标题和表格操作所需的样式
    *   服务器会尝试创建缺失的样式或使用直接格式
    *   为获得最佳效果，请使用带有标准Word样式的模板

2.  **权限问题**

    *   确保服务器有权读取/写入文档路径
    *   使用`copy_document`函数创建锁定文档的可编辑副本
    *   如果操作失败，请检查文件所有权和权限

3.  **图片插入问题**

    *   对图片文件使用绝对路径
    *   验证图片格式兼容性（推荐JPEG、PNG）
    *   检查图片文件大小和权限

### 调试

通过设置环境变量启用详细日志：

```
export MCP_DEBUG=1  # Linux/macOS系统
set MCP_DEBUG=1     # Windows系统
```

贡献指南
------------

欢迎贡献！请随时提交拉取请求。

1.   Fork仓库
2.  创建你的功能分支（`git checkout -b feature/amazing-feature`）
3.  提交你的更改（`git commit -m 'Add some amazing feature'`）
4.  推送到分支（`git push origin feature/amazing-feature`）
5.  打开拉取请求

许可证
-------

本项目基于MIT许可证授权 - 详见LICENSE文件了解详情。

鸣谢
---------------

*   [模型上下文协议](https://modelcontextprotocol.io/)提供协议规范
*   [python-docx](https://python-docx.readthedocs.io/)用于Word文档操作
*   [FastMCP](https://github.com/modelcontextprotocol/python-sdk)提供Python MCP实现

* * *

_注意：此服务器会与你系统上的文档文件进行交互。在Claude桌面版或其他MCP客户端中确认请求的操作前，务必验证这些操作是否合适。_