Markdownify MCP 服务器
----------------------

Markdownify 是一个模型上下文协议（MCP）服务器，它可以将各种文件类型和网页内容转换为 Markdown 格式。它提供了一套工具，可以将 PDF、图片、音频文件、网页等内容转换成易于阅读和分享的 Markdown 文本。


功能
--------

*   将多种文件类型转换为 Markdown：
    *   PDF
    *   图片
    *   音频（带转录）
    *   DOCX
    *   XLSX
    *   PPTX


*   将网页内容转换为 Markdown：
    *   YouTube 视频字幕
    *   Bing 搜索结果
    *   一般网页


*   获取现有的 Markdown 文件


开始使用
---------------

1.   克隆此仓库
2.   安装依赖项：
```
pnpm install
```

注意：这也会安装 `uv` 和相关的 Python 依赖项。

3.    构建项目：
```
pnpm run build
```
4.    启动服务器：
```
pnpm start
```


开发
-----------

*   使用 `pnpm run dev` 以监视模式启动 TypeScript 编译器
*   修改 `src/server.ts` 来自定义服务器行为
*   在  `src/tools.ts`中添加或修改工具


与桌面应用程序一起使用
----------------------

要将此服务器与桌面应用程序集成，请在应用程序的服务器配置中添加以下内容：

```
{
  "mcpServers": {
    "markdownify": {
      "command": "node",
      "args": [
        "{ABSOLUTE PATH TO FILE HERE}/dist/index.js"
      ],
      "env": {
        // By default, the server will use the default install location of `uv`
        "UV_PATH": "/path/to/uv"
      }
    }
  }
}
```

可用工具
---------------
*   `youtube-to-markdown`：将 YouTube 视频转换为 Markdown
*   `pdf-to-markdown`：将 PDF 文件转换为 Markdown
*   `bing-search-to-markdown`：将 Bing 搜索结果转换为 Markdown
*   `webpage-to-markdown`：将网页转换为 Markdown
*   `image-to-markdown`：将图片转换为带有元数据的 Markdown
*   `audio-to-markdown`：将音频文件转换为带有转录的 Markdown
*   `docx-to-markdown`：将 DOCX 文件转换为 Markdown
*   `xlsx-to-markdown`：将 XLSX 文件转换为 Markdown
*   `pptx-to-markdown`：将 PPTX 文件转换为 Markdown
*   `get-markdown-file`：获取现有的 Markdown 文件


贡献
------------

欢迎贡献！请随时提交 Pull Request。


许可证
-------

此项目根据 MIT 许可证发布 - 详情请参阅 LICENSE 文件。